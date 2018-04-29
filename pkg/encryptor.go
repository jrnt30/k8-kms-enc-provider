// Copyright © 2018 Justin Nauman <justin@spantree.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package pkg

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"

	"github.com/jrnt30/k8-kms-enc-provider/v1beta1"
)

// Ensure that our implemetnation stays in contract with the Protobuf's specification
var _ v1beta1.KeyManagementServiceServer = AwsKmsProvider{}

// AwsKmsProviderConfiguration allows for the customization of the
// KMS provider with some sensible defaults
type AwsKmsProviderConfiguration struct {
	// KeyId is the identifier for KMS key to use for encryption.
	// Can be either the Key ARN or the Key ID.
	// NOTE: Key Alias support is currently not implemented with the existing
	// validation logic
	KeyId *string

	// AwsRegion is the specifier on which AWS Region the KMS key resides in.
	AwsRegion *string
}

// NewAwsKmsProvider is a helper for generating a new KMS Key proxy that provide sensible defaults
func NewAwsKmsProvider(cfg *AwsKmsProviderConfiguration) (*AwsKmsProvider, error) {
	if cfg.KeyId == nil {
		return nil, errors.New("KeyId is a required attribute that must be provided ")
	}

	var sess *session.Session
	awsConfigs := make([]*aws.Config, 0)
	if cfg.AwsRegion != nil {
		awsConfigs = append(awsConfigs, &aws.Config{
			Region: cfg.AwsRegion,
		})
	}
	sess, err := session.NewSession(awsConfigs...)
	if err != nil {
		return nil, err
	}
	svc := kms.New(sess)

	foundKey := false
	svc.ListKeysPages(&kms.ListKeysInput{}, func(page *kms.ListKeysOutput, lastPage bool) bool {
		for _, kmsKey := range page.Keys {
			fmt.Println(*kmsKey.KeyArn)
			if *kmsKey.KeyArn == *cfg.KeyId || *kmsKey.KeyId == *cfg.KeyId {
				foundKey = true
				break
			}
		}
		return !foundKey && !lastPage
	})

	if !foundKey {
		return nil, fmt.Errorf("Unable to locate the KMS key [%s]. Ensure this is a valid ARN or ID for the region", *cfg.KeyId)
	}

	return &AwsKmsProvider{
		sess: sess,
		kms:  svc,
	}, nil
}

// AwsKmsProvider is an implementation of the K8 KMS provider
// specification https://kubernetes.io/docs/tasks/administer-cluster/kms-provider/
//
// In this implemenation we are using AWS KMS's encryption functionality to convert the plaintext into a
// ciphertext that is capable of being stored securely.
//
type AwsKmsProvider struct {
	sess    *session.Session
	kms     *kms.KMS
	keyName string
	region  string
}

// Version returns API information to consumers (primarily just the K8 masters themselves )
func (a AwsKmsProvider) Version(context.Context, *v1beta1.VersionRequest) (*v1beta1.VersionResponse, error) {
	return &v1beta1.VersionResponse{RuntimeName: "kms-enc-provider", RuntimeVersion: "v1beta1", Version: "v1beta1"}, nil
}

// Decrypt is responsible for converting the *v1beta1.DecryptRequest.Cipher into a plaintext representation
// K8 itself.
func (a AwsKmsProvider) Decrypt(c context.Context, req *v1beta1.DecryptRequest) (*v1beta1.DecryptResponse, error) {
	decResp, err := a.kms.Decrypt(&kms.DecryptInput{
		CiphertextBlob: req.Cipher})

	if err != nil {
		return nil, fmt.Errorf("KMS Dec. Error encountered: %s", err)
	}

	return &v1beta1.DecryptResponse{Plain: decResp.Plaintext}, nil
}

// Encrypt is responsible for taking the plaintext from *v1beta1.EncryptRequest.Plain and transparently
// encrypting the value for K8.
func (a AwsKmsProvider) Encrypt(c context.Context, req *v1beta1.EncryptRequest) (*v1beta1.EncryptResponse, error) {
	encResp, err := a.kms.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String("alias/k8envelop"),
		Plaintext: req.Plain})

	if err != nil {
		return nil, fmt.Errorf("KMS Dec. Error encountered: %s", err)
	}

	return &v1beta1.EncryptResponse{Cipher: encResp.CiphertextBlob}, nil
}

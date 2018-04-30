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

package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/jrnt30/k8-kms-enc-provider/v1beta1"
)

var cipherText string

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Allows for the decryption of cipher-text that was encrypted previously",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		gc, err := grpc.Dial(socketPath,
			grpc.WithInsecure(),
			grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
				return net.DialTimeout("unix", addr, timeout)
			}))

		if err != nil {
			log.Fatal(err)
		}
		client := v1beta1.NewKeyManagementServiceClient(gc)

		decoded, _ := base64.StdEncoding.DecodeString(cipherText)

		resp, err := client.Decrypt(context.Background(), &v1beta1.DecryptRequest{
			Cipher: decoded,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", string(resp.Plain))
	},
}

func init() {
	clientCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringVar(&cipherText, "cipher-text", "", "Ciphertext To Decrypt")
	decryptCmd.MarkFlagRequired("cipher-text")
}

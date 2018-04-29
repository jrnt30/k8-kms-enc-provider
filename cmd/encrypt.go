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

	"github.com/jrnt30/k8-kms-enc-provider/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var plainText string

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		resp, err := client.Encrypt(context.Background(), &v1beta1.EncryptRequest{
			Plain: []byte(plainText),
		})

		encodedCipherString := base64.StdEncoding.EncodeToString(resp.Cipher)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", string(encodedCipherString))
	},
}

func init() {
	clientCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringVar(&plainText, "plain-text", "", "Plain text to encrypt")
	encryptCmd.MarkFlagRequired("plain-text")
}

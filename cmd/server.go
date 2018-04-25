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
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/jrnt30/aws-kms-k8-enc-provider/pkg"
	"github.com/jrnt30/aws-kms-k8-enc-provider/v1beta1"
)

var awsRegion string
var keyID string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		keyProviderServer, err := pkg.NewAwsKmsProvider(&pkg.AwsKmsProviderConfiguration{
			AwsRegion: aws.String(awsRegion),
			KeyId:     aws.String(keyID),
		})

		server := grpc.NewServer()
		lis, err := net.Listen("unix", socketPath)
		defer lis.Close()
		if err != nil {
			log.Fatal("Error creating the socket listener, existing", err)
		}

		sigTerm := make(chan os.Signal, 1)
		signal.Notify(sigTerm, os.Interrupt, os.Kill, syscall.SIGTERM)

		waits := sync.WaitGroup{}
		waits.Add(1)

		go func() {
			v1beta1.RegisterKeyManagementServiceServer(server, keyProviderServer)
			server.Serve(lis)
			waits.Done()
		}()

		go func(term chan os.Signal) {
			t := <-term
			log.Printf("Closing due to %s signal", t)
			waits.Done()
		}(sigTerm)

		waits.Wait()

	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&awsRegion, "region", "", "Region to load the associated KMS Key from")
	serverCmd.Flags().StringVar(&keyID, "key-id", "", "KMS Key Identifier (ID or ARN) to be used for encryption")
	serverCmd.MarkFlagRequired("region")
	serverCmd.MarkFlagRequired("key-id")
}

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

package server

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

	"github.com/jrnt30/k8-kms-enc-provider/cmd"
	"github.com/jrnt30/k8-kms-enc-provider/pkg"
	"github.com/jrnt30/k8-kms-enc-provider/v1beta1"
)

var awsRegion string
var keyID string
var socketPath string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launches the K8 KMS server component that listens on a socket",
	Run: func(cmd *cobra.Command, args []string) {
		keyProviderServer, err := pkg.NewAwsKmsProvider(&pkg.AwsKmsProviderConfiguration{
			AwsRegion: aws.String(awsRegion),
			KeyId:     aws.String(keyID),
		})
		if err != nil {
			log.Fatal("Error creating the backing KMS provider: ", err)
		}

		server := grpc.NewServer()
		addr, err := net.ResolveUnixAddr("unix", socketPath)
		if err != nil {
			log.Fatal("Error resolving the socket, existing", err)
		}

		lis, err := net.ListenUnix("unix", addr)
		if err != nil {
			log.Fatal("Error creating the socket listener, existing", err)
		}
		defer lis.Close()

		sigTerm := make(chan os.Signal, 1)
		signal.Notify(sigTerm, os.Interrupt, os.Kill, syscall.SIGTERM)

		waits := sync.WaitGroup{}
		waits.Add(1)

		go func() {
			log.Print("Registering listener on the GRPC server")
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
	cmd.RootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&socketPath, "socket", "/tmp/kms-grpc", "path to the socket to use")
	serverCmd.Flags().StringVar(&awsRegion, "region", "", "Region to load the associated KMS Key from")
	serverCmd.Flags().StringVar(&keyID, "key-id", "", "KMS Key Identifier (ID or ARN) to be used for encryption")
	serverCmd.MarkFlagRequired("region")
	serverCmd.MarkFlagRequired("key-id")
}

/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	gossip "github.com/sheikhshack/distributed-chaos-50.041/node/gossip/proto"
	"google.golang.org/grpc"
	"log"
	"os"
	"github.com/spf13/cobra"
)

var fileName string
var content string

// writeFileCmd represents the writeFile command
var writeFileCmd = &cobra.Command{
	Use:   "writefile",
	Short: "Writes to the chord node",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("writeFile called")
		attachedNode := os.Getenv("APP_NODE")
		writeExternalFile(attachedNode, fileName, content)
	},
}

var readFileCmd = &cobra.Command{
	Use:   "readfile",
	Short: "Read file form chord node",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("read called")
	},
}

func init() {
	rootCmd.AddCommand(writeFileCmd)
	rootCmd.AddCommand(readFileCmd)
	writeFileCmd.PersistentFlags().StringVarP(&fileName, "fileName", "f","", "filename to insert")
	writeFileCmd.PersistentFlags().StringVarP(&content, "content", "c","", "content to insert")


}
func writeExternalFile(nodeAddr, fileName, containerIP string) {
	log.Printf("Attempting to write file %v to node %v", fileName, nodeAddr)
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, 9000)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()
	client := gossip.NewInternalListenerClient(conn)
	res, err := client.StoreKeyHash(context.Background(), &gossip.DLUploadRequest{
		Filename:    fileName,
		ContainerIP: containerIP,
	})
	if err != nil {
		log.Fatalf("-- MIKE external fail %s", err)
	}

	log.Printf("\nSuccess upload info to the following chord node: %+v\n", res)
}

// Unimplemented
func resolveFile(nodeAddr, fileName string) {
	//log.Printf("Attempting to lookup string %v", key)
	var conn *grpc.ClientConn
	connectionParams := fmt.Sprintf("%s:%v", nodeAddr, 9000)

	conn, err := grpc.Dial(connectionParams, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to: %s", err)

	}
	defer conn.Close()
	client := gossip.NewInternalListenerClient(conn)
	res, err := client.GetFileLocation(context.Background(), &gossip.DLDownloadRequest{
		Filename: fileName,
	})
	if err != nil {
		log.Fatalf("-- MIKE external fail %s", err)
	}

	log.Printf("\nSuccess, received the following containerINFO: %+v\n", res)
}
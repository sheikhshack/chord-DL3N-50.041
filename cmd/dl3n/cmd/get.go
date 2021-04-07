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
	"os"

	"github.com/sheikhshack/distributed-chaos-50.041/dl3n"
	"github.com/spf13/cobra"
)

var getChordAddr string

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [meta_file_path]",
	Short: "download a file",
	Long: `get downloads the file using its metadata.
get will download files in chunks from dl3n nodes seeding this file.`,
	Args: cobra.ExactArgs(1),
	Run:  get,
}

func get(cmd *cobra.Command, args []string) {
	path := args[0]
	d, err := dl3n.NewDL3NFromMeta(path)
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}

	nd := dl3n.NewChordNodeDiscovery(getChordAddr)

	ds := dl3n.NewDL3NNode(d, nd)
	err = ds.Get()
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&getChordAddr, "addr", "a", "127.0.0.1", "address of a chord DHT node to connect to.")
}

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

// writemetaCmd represents the writemeta command
var writemetaCmd = &cobra.Command{
	Use:   "writemeta [file_path]",
	Short: "create a metadata file without sharing the file",
	Long: `writemeta generates a .dl3n metadata file which contains metadata about
the file being shared, such as its name, infohash, and chunks.

unlike seed, writemeta will not chunk the file, only generate the metadata file.
`,
	Args: cobra.ExactArgs(1),
	Run:  writemeta,
}

func writemeta(cmd *cobra.Command, args []string) {
	path := args[0]

	d, err := dl3n.NewDL3NFromFileOneChunk(path)
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}

	err = d.WriteMetaFile(path + ".dl3n")
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(writemetaCmd)
}

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
	"fmt"

	"github.com/spf13/cobra"
)

// writemetaCmd represents the writemeta command
var writemetaCmd = &cobra.Command{
	Use:   "writemeta",
	Short: "create a metadata file without sharing the file",
	Long: `writemeta generates a .dl3n metadata file which contains metadata about
the file being shared, such as its name, infohash, and chunks.

unlike seed, writemeta will not chunk the file, only generate the metadata file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("writemeta called")
	},
}

func init() {
	rootCmd.AddCommand(writemetaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// writemetaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// writemetaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

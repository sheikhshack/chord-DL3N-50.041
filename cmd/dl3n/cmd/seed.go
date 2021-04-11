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
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/sheikhshack/distributed-chaos-50.041/dl3n"
)

var seedChordAddr string

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed [file_path]",
	Short: "start sharing a file",
	Long: `seed starts seeding a file, allowing other dl3n clients 
to download chunks of the file from seeding nodes.

seed will generate a .dl3n metadata file which contains metadata about
the file being shared, such as its name, infohash, and chunks.
`,
	Args: cobra.ExactArgs(1),
	Run:  seed,
}

func seed(cmd *cobra.Command, args []string) {
	path := args[0]

	d, err := dl3n.NewDL3NFromFileOneChunk(path)
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}

	d.WriteMetaFile(path + ".dl3n")
	nd := dl3n.NewChordNodeDiscovery(seedChordAddr)
	ds := dl3n.NewDL3NNode(d, nd)

	ip := getIP()
	seederAddr := ip + ":11111"
	err = ds.NodeDiscovery.SetSeederAddr(d.Hash, seederAddr)
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}

	ds.StartSeed(seederAddr)

	// listen for signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// wait for sigint or sigterm
	<-sigs
	ds.StopSeed()
}

func init() {
	rootCmd.AddCommand(seedCmd)

	seedCmd.Flags().StringVarP(&seedChordAddr, "addr", "a", "127.0.0.1", "address of a chord DHT node to connect to.")

}

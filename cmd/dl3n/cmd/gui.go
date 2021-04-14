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
	"github.com/sheikhshack/distributed-chaos-50.041/dl3n"
	"github.com/sheikhshack/distributed-chaos-50.041/dl3n/gui"
	"github.com/spf13/cobra"
)

var guiChordAddr string

// guiCmd represents the gui command
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "run a server to serve the gui",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runGui,
}

func runGui(cmd *cobra.Command, args []string) {
	ip := getIP()
	seederAddr := ip + ":11111"

	nd := dl3n.NewChordNodeDiscovery(guiChordAddr)
	g := gui.NewGuiServer(nd, seederAddr)
	g.StartServer()
}

func init() {
	rootCmd.AddCommand(guiCmd)

	guiCmd.Flags().StringVarP(&guiChordAddr, "addr", "a", "127.0.0.1", "address of a chord DHT node to connect to.")
}

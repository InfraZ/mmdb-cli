/*
Copyright 2024 The InfraZ Authors.

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

	"github.com/InfraZ/mmdb-cli/pkg/dump"
	"github.com/spf13/cobra"
)

var cmdDumpConfig dump.CmdDumpConfig

const (
	dumpCmdName      = "dump"
	dumpCmdShortDesc = "Dump MMDB data into a json dataset"
	dumpCmdLongDesc  = `This command dumps MMDB data into a json dataset`
)

// dumpCmd represents the generate command
var dumpCmd = &cobra.Command{
	Use:   dumpCmdName,
	Short: dumpCmdShortDesc,
	Long:  dumpCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		err := dump.DumpMMMDB(&cmdDumpConfig)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	// Add flags to the update command
	dumpCmd.Flags().StringVarP(&cmdDumpConfig.InputDatabase, "input", "i", "", "Input path of the MMDB file")
	dumpCmd.Flags().StringVarP(&cmdDumpConfig.OutputFile, "output", "o", "", "Output path of the output JSON dataset file (must have a .json extension)")

	// Mark required flags
	dumpCmd.MarkFlagRequired("input")
	dumpCmd.MarkFlagRequired("output")
}

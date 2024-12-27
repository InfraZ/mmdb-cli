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

	"github.com/InfraZ/mmdb-cli/pkg/generate"
	"github.com/spf13/cobra"
)

var cmdGenerateConfig generate.CmdGenerateConfig

const (
	generateCmdName      = "generate"
	generateCmdShortDesc = "Generate a MMDB database from a JSON dataset"
	generateCmdLongDesc  = `This command generates a MMDB database from a JSON dataset.`
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   generateCmdName,
	Short: generateCmdShortDesc,
	Long:  generateCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		err := generate.GenerateMMDB(&cmdGenerateConfig)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	// Add flags to the update command
	generateCmd.Flags().StringVarP(&cmdGenerateConfig.InputDataset, "input", "i", "", "Input path of the JSON dataset file (must have a .json extension)")
	generateCmd.Flags().StringVarP(&cmdGenerateConfig.OutputDatabase, "output", "o", "", "Output path of the MMDB database file (must have a .mmdb extension)")
	generateCmd.Flags().BoolVarP(&cmdGenerateConfig.Verbose, "verbose", "v", false, "Enable verbose mode")

	generateCmd.Flags().BoolVar(&cmdGenerateConfig.DisableIPv4Aliasing, "disable-ipv4-aliasing", false, "Disable IPv4 aliasing")
	generateCmd.Flags().BoolVar(&cmdGenerateConfig.IncludeReservedNetworks, "include-reserved-networks", false, "Include reserved networks")

	// Mark required flags
	generateCmd.MarkFlagRequired("input")
	generateCmd.MarkFlagRequired("output")
}

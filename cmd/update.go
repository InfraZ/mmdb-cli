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

	"github.com/spf13/cobra"

	"github.com/InfraZ/mmdb-cli/pkg/update"
)

var cmdUpdateConfig update.CmdUpdateConfig

const (
	updateCmdName      = "update"
	updateCmdShortDesc = "Update existing MMDB file"
	updateCmdLongDesc  = `This command updates an existing MMDB file with new data`
)

// updateCmd represents the generate command
var updateCmd = &cobra.Command{
	Use:   updateCmdName,
	Short: updateCmdShortDesc,
	Long:  updateCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		err := update.UpdateMMDB(cmdUpdateConfig)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	// Add flags to the update command
	updateCmd.Flags().StringVarP(&cmdUpdateConfig.InputDatabase, "input", "i", "", "Input path of the MMDB file")
	updateCmd.Flags().StringVarP(&cmdUpdateConfig.InputDataSet, "dataset", "d", "", "Input path of the dataset file")
	updateCmd.Flags().StringVarP(&cmdUpdateConfig.OutputDatabase, "output", "o", "", "Output path of the MMDB file")
	updateCmd.Flags().BoolVarP(&cmdUpdateConfig.Verbose, "verbose", "v", false, "Enable verbose mode")

	updateCmd.Flags().BoolVar(&cmdUpdateConfig.DisableIPv4Aliasing, "disable-ipv4-aliasing", false, "Disable IPv4 aliasing")
	updateCmd.Flags().BoolVar(&cmdUpdateConfig.IncludeReservedNetworks, "include-reserved-networks", false, "Include reserved networks")

	// Mark required flags
	updateCmd.MarkFlagRequired("input")
	updateCmd.MarkFlagRequired("dataset")
	updateCmd.MarkFlagRequired("output")
}

// Copyright 2024 The MMDB CLI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"
	"mmdb-cli/pkg/inspect"
	"mmdb-cli/pkg/output"

	"github.com/spf13/cobra"
)

const (
	inspectCmdName      = "inspect"
	inspectCmdShortDesc = "Inspect an IP address or CIDR in the MMDB file"
	inspectCmdLongDesc  = `This command allows you to inspect an IP address or CIDR in the MMDB file`
)

var cmdInspectConfig inspect.CmdInspectConfig

// inspectCmd represents the generate command
var inspectCmd = &cobra.Command{
	Use:   inspectCmdName,
	Short: inspectCmdShortDesc,
	Long:  inspectCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		// Set the inputs
		cmdInspectConfig.Inputs = cmd.Flags().Args()

		inspectResult, err := inspect.InspectInMMDB(cmdInspectConfig)
		if err != nil {
			log.Fatal(err)
		}

		err = output.Output(inspectResult, outputOptions)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// Add flags to the inspect command
	inspectCmd.Flags().StringVarP(&cmdInspectConfig.InputFile, "input", "i", "", "Input path of the MMDB file")
	inspectCmd.Flags().StringVarP(&outputOptions.Format, "format", "f", "yaml", "Output format (yaml, json, json-pretty)")

	// Mark required flags
	inspectCmd.MarkFlagRequired("input")
}

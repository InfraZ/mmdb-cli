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
	"log"

	"github.com/InfraZ/mmdb-cli/pkg/verify"

	"github.com/spf13/cobra"
)

const (
	verifyCmdName      = "verify"
	verifyCmdShortDesc = "Verify the MMDB file"
	verifyCmdLongDesc  = `This command verifies the MMDB file`
)

var cmdVerifyConfig verify.CmdVerifyConfig

// verifyCmd represents the generate command
var verifyCmd = &cobra.Command{
	Use:   verifyCmdName,
	Short: verifyCmdShortDesc,
	Long:  verifyCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		verifyResult, err := verify.VerifyMMDB(cmdVerifyConfig)
		if err != nil {
			log.Fatal(err)
		}

		if verifyResult {
			fmt.Println("The MMDB file is valid")
		} else {
			fmt.Println("The MMDB file is invalid")
		}
	},
}

func init() {
	// Add flags to the inspect command
	verifyCmd.Flags().StringVarP(&cmdVerifyConfig.InputFile, "input", "i", "", "Input path of the MMDB file")

	// Mark required flags
	verifyCmd.MarkFlagRequired("input")
}

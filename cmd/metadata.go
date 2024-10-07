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
	"mmdb-cli/pkg/metadata"
	"mmdb-cli/pkg/output"

	"github.com/spf13/cobra"
)

const (
	metadataCmdName      = "metadata"
	metadataCmdShortDesc = "Prints metadata of the MMDB file"
	metadataCmdLongDesc  = `This command prints metadata of the MMDB file`
)

var cmdMetadataConfig metadata.CmdMetadataConfig

// metadataCmd represents the generate command
var metadataCmd = &cobra.Command{
	Use:   metadataCmdName,
	Short: metadataCmdShortDesc,
	Long:  metadataCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		metadataJson, err := metadata.MetadataMMDB(cmdMetadataConfig)
		if err != nil {
			log.Fatal(err)
		}

		err = output.Output(metadataJson, outputOptions)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// Add flags to the metadata command
	metadataCmd.Flags().StringVarP(&cmdMetadataConfig.InputFile, "input", "i", "", "Input path of the MMDB file")
	metadataCmd.Flags().StringVarP(&outputOptions.Format, "format", "f", "yaml", "Output format (yaml, json, json-pretty)")

	// Mark required flags
	metadataCmd.MarkFlagRequired("input")
}

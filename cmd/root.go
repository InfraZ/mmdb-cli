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
	"os"

	"github.com/InfraZ/mmdb-cli/pkg/output"

	"github.com/spf13/cobra"
)

var outputOptions output.OutputOptions

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mmdb-cli",
	Short: "InfraZ MMDB CLI is a command line toolkit for working with MMDB",
	Long: `
InfraZ MMDB CLI is a command line toolkit for working with MMDB
Complete documentation is available at https://docs.infraz.io/mmdb-cli`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(metadataCmd)
	rootCmd.AddCommand(inspectCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(dumpCmd)
	rootCmd.AddCommand(generateCmd)
}

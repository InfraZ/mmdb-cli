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
	"log"
	"os"

	"github.com/InfraZ/mmdb-cli/pkg/diff"
	"github.com/spf13/cobra"
)

const (
	diffCmdName      = "diff"
	diffCmdShortDesc = "Show differences between two MMDB files (experimental)"
	diffCmdLongDesc  = `[Experimental] Compare two MMDB files and show added, removed, and modified networks.

This command is experimental: its flags and output format may change in a future release without following semantic versioning.

Exits with code 0 when no differences are found, code 1 when differences exist.

Examples:
  # Full diff
  mmdb-cli diff file-a.mmdb file-b.mmdb

  # Diff scoped to a specific CIDR
  mmdb-cli diff file-a.mmdb file-b.mmdb 1.0.0.0/8

  # Diff scoped to multiple IPs or CIDRs
  mmdb-cli diff file-a.mmdb file-b.mmdb 1.1.1.1 8.8.8.0/24`
)

var diffCmd = &cobra.Command{
	Use:   diffCmdName + " <file-a.mmdb> <file-b.mmdb> [ip-or-cidr ...]",
	Short: diffCmdShortDesc,
	Long:  diffCmdLongDesc,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := diff.CmdDiffConfig{
			InputA:  args[0],
			InputB:  args[1],
			Filters: args[2:],
		}

		result, err := diff.DiffMMDB(cfg)
		if err != nil {
			log.Fatal(err)
		}

		diff.PrintResult(result)

		if result.HasDiff {
			os.Exit(1)
		}
	},
}

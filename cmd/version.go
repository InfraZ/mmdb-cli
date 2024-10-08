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
)

const (
	version    = "v0.1.0-alpha"
	maintainer = "The InfraZ Authors"
)

const (
	versionCmdName      = "version"
	versionCmdShortDesc = "Show version information for mmdb-cli"
	versionCmdLongDesc  = `This command shows version information for mmdb-cli`
)

// versionCmd represents the generate command
var versionCmd = &cobra.Command{
	Use:   versionCmdName,
	Short: versionCmdShortDesc,
	Long:  versionCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: \t%s\nMaintainer: \t%s\n", version, maintainer)
	},
}

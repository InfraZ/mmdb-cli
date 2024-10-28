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

	"github.com/InfraZ/mmdb-cli/internal/metadata"
	"github.com/spf13/cobra"
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

		fmt.Printf("Version: %s\n", metadata.Version)
		fmt.Printf("Licence: %s\n", metadata.License)
		fmt.Printf("Description: %s\n", metadata.ShortDescription)
		fmt.Printf("Documentation: %s\n", metadata.DocumentationURL)
		fmt.Println("Maintainers:")
		for _, maintainer := range metadata.Maintainers {
			fmt.Printf("\t- %s\n", maintainer)
		}
		fmt.Printf("\nPlease support us by starring the project on GitHub: %s\n", metadata.Homepage)

	},
}

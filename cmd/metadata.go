/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"mmdb-cli/pkg/metadata"

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
		metadata.MetadataMMDB(cmdMetadataConfig)
	},
}

func init() {

	metadataCmd.Flags().StringVarP(&cmdMetadataConfig.InputFile, "input", "i", "", "Input path of the MMDB file")

}

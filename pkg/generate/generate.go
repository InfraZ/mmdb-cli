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

package generate

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/InfraZ/mmdb-cli/internal/files"
	"github.com/InfraZ/mmdb-cli/pkg/mmdb"
	"github.com/maxmind/mmdbwriter"
)

type CmdGenerateConfig struct {
	InputDataset   string
	OutputDatabase string
	Verbose        bool

	DisableIPv4Aliasing     bool
	IncludeReservedNetworks bool
}

/*
Structure of the dumped JSON dataset:
{
	"version": "v1",
	"schema": {
		<SCHEMA>
	},
	"metadata": {
		<METADATA>
	},
	"dataset": [
		{
			"network": "<NETWORK>",
			"record": {
				<RECORD>
			}
		}
	]
}
*/

func readDataSet(inputDataSet string) (map[string]interface{}, error) {
	_, err := os.Stat(inputDataSet)
	if os.IsNotExist(err) {
		log.Fatalf("File %s does not exist", inputDataSet)
		return nil, err
	}

	datasetFile, err := os.Open(inputDataSet)
	if err != nil {
		return nil, err
	}
	defer datasetFile.Close()

	byteValue, err := io.ReadAll(datasetFile)
	if err != nil {
		return nil, err
	}

	var dataset map[string]interface{}
	if err := json.Unmarshal(byteValue, &dataset); err != nil {
		return nil, err
	}

	return dataset, nil
}

func mmdbWriterOptions(cfg *CmdGenerateConfig, metadata map[string]interface{}) (*mmdbwriter.Options, error) {

	// BuildEpoch
	if metadata["BuildEpoch"] != nil {
		fmt.Println("[-] BuildEpoch in metadata will be ignored")
	}

	// DatabaseType
	var databaseType string
	if metadata["DatabaseType"] == nil {
		log.Fatalf("\n[!] DatabaseType is a required field in metadata")
	} else {
		databaseType = metadata["DatabaseType"].(string)
	}

	// Description
	description := make(map[string]string)
	if metadata["Description"] == nil {
		log.Fatalf("\n[!] Description is a required field in metadata")
	} else {
		// Convert the description map[string]interface {} to map[string]string
		for descriptionKey, descriptionValue := range metadata["Description"].(map[string]interface{}) {
			description[descriptionKey] = descriptionValue.(string)
		}
	}

	// IPVersion
	var ipVersion int
	if metadata["IPVersion"] == nil {
		// Default to IPv6
		metadata["IPVersion"] = 6

		fmt.Println("[-] IPVersion is not provided in metadata, defaulting to 6 (An IPv6 database supports both IPv4 and IPv6 lookups)")
	} else {
		// Convert the IPVersion to int
		ipVersion = int(metadata["IPVersion"].(float64))

		if ipVersion != 4 && ipVersion != 6 {
			log.Fatalf("\n[!] Invalid value for IPVersion in metadata. The supported values are 4 and 6")
		}
	}

	// Languages
	languages := make([]string, 0)
	if metadata["Languages"] == nil {
		// Default to English
		languages = append(languages, "en")

		fmt.Println("[-] Languages is not provided in metadata, defaulting to English")
	} else {
		// Convert the languages []interface {} to []string
		for _, language := range metadata["Languages"].([]interface{}) {
			languages = append(languages, language.(string))
		}
	}

	var recordSize int
	if metadata["RecordSize"] == nil {
		// Default to 28
		recordSize = 28

		fmt.Println("[-] RecordSize is not provided in metadata, defaulting to 28 (The supported values are 24, 28, and 32)")
	} else {
		// Convert the RecordSize to int
		recordSize = int(metadata["RecordSize"].(float64))

		if recordSize != 24 && recordSize != 28 && recordSize != 32 {
			log.Fatalf("\n[!] Invalid value for RecordSize in metadata. The supported values are 24, 28, and 32")
		}
	}

	// Initialize the MMDB writer options
	mmdbWriterOptions := &mmdbwriter.Options{
		DatabaseType:            databaseType,
		Description:             description,
		DisableIPv4Aliasing:     cfg.DisableIPv4Aliasing,
		IncludeReservedNetworks: cfg.IncludeReservedNetworks,
		IPVersion:               ipVersion,
		Languages:               languages,
		RecordSize:              recordSize,
	}

	return mmdbWriterOptions, nil
}

func initializeMMDBWriter(cfg *CmdGenerateConfig, metadata map[string]interface{}) (*mmdbwriter.Tree, error) {

	mmdbWriterOptions, err := mmdbWriterOptions(cfg, metadata)
	if err != nil {
		log.Fatalf("\n[!] Error initializing MMDB writer options: %v", err)
	}

	writer, err := mmdbwriter.New(*mmdbWriterOptions)
	if err != nil {
		log.Fatalf("\n[!] Error initializing MMDB writer: %v", err)
	}

	return writer, err
}

func GenerateMMDB(cfg *CmdGenerateConfig) error {

	// Validate files
	filesToCheck := []files.FilesListValidation{
		{FilePath: cfg.InputDataset, ExpectedExtension: ".json", ShouldExist: true},
	}

	if err := files.FilesValidation(filesToCheck); err != nil {
		log.Fatal(err)
	}

	// Initialize the record position
	var recordPosition int = 0

	// Read the dataset
	var dataset map[string]interface{}
	dataset, err := readDataSet(cfg.InputDataset)
	if err != nil {
		log.Fatalf("\n[!] Error reading dataset: %v", err)
	}

	// Extract metadata from the dataset
	metadata := dataset["metadata"]

	// Extract schema from the dataset
	var schema map[string]interface{}
	var useDefaultSchema bool = true

	// Extract and validate version if present
	if versionInterface, exists := dataset["version"]; exists {
		if version, ok := versionInterface.(string); ok {
			if version != "v1" {
				log.Fatalf("\n[!] Unsupported dataset version: %s (supported: v1)", version)
			}
			fmt.Printf("[+] Dataset version: %s\n", version)
		}
	}

	if schemaInterface, exists := dataset["schema"]; exists {
		if schemaMap, ok := schemaInterface.(map[string]interface{}); ok {
			schema = schemaMap
			useDefaultSchema = false
			fmt.Printf("[+] Using dynamic schema from dataset: %v\n", schema)
		} else {
			log.Printf("[-] Schema field exists but is not a valid object, falling back to default schema\n")
		}
	} else {
		fmt.Println("[-] No schema found in dataset, using default schema")
	}

	// Initialize the MMDB writer
	writer, err := initializeMMDBWriter(cfg, metadata.(map[string]interface{}))
	if err != nil {
		log.Fatalf("\n[!] Error initializing MMDB writer: %v", err)
	}

	// Iterate over the dataset and write to the MMDB
	for _, dataset := range dataset["dataset"].([]interface{}) {

		recordPosition++

		// Type assertion to map[string]interface{}
		dataMap := dataset.(map[string]interface{})

		// Record network
		_, network, err := net.ParseCIDR(dataMap["network"].(string))
		if err != nil {
			log.Fatalf("\n[!] Invalid network (%s) in the dataset: %v", dataMap["network"].(string), err)
		}

		// Parse dynamic data
		dynamicData, exists := dataMap["record"].(map[string]interface{})
		if !exists {
			log.Fatalf("\n[!] Error parsing data for record %d (network: %s) - %v", recordPosition, network, err)
		}

		// Convert dynamic data to MMDB type map (using dynamic schema if available)
		dynamicMmdbData := mmdb.ConvertToMMDBTypeMap(dynamicData, useDefaultSchema, schema)

		if err := writer.Insert(network, dynamicMmdbData); err != nil {
			log.Fatalf("\n[!] Error inserting record %d (network: %s) - %v", recordPosition, network, err)
		}

		if cfg.Verbose {
			fmt.Printf("[-] Inserting record %d for network %s - data: %v\n", recordPosition, network, dynamicMmdbData)
		} else {
			fmt.Printf("\r[-] Inserted %d records", recordPosition)
		}
	}

	fmt.Printf("\r[+] Total records inserted: %d\n", recordPosition)

	// Write the MMDB database to the output file
	outputFile, err := os.Create(cfg.OutputDatabase)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write the MMDB database to the output file
	fmt.Println("[+] Writing MMDB database to the output file")
	if _, err = writer.WriteTo(outputFile); err != nil {
		return err
	}

	// check the output file size
	outputDatabaseStat, err := outputFile.Stat()
	if err != nil {
		return fmt.Errorf("[!] Failed to get output file stats: %s - %v", cfg.OutputDatabase, err)
	}

	// convert outputDatabaseStat.Size() to MB
	outputDatabaseSizeMB := float64(outputDatabaseStat.Size()) / 1024 / 1024
	fmt.Printf("\r[+] %s file created with size: %.2f MB\n", cfg.OutputDatabase, outputDatabaseSizeMB)

	fmt.Println("[+] MMDB Generated successfully")

	return err
}

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

	"github.com/InfraZ/mmdb-cli/pkg/mmdb"
	"github.com/maxmind/mmdbwriter"
)

type CmdGenerateConfig struct {
	InputDataset   string
	OutputDatabase string

	DisableIPv4Aliasing     bool
	IncludeReservedNetworks bool
}

/*
Structure of the dumped JSON dataset:
{
	"schema": "v1",
	"metadata": {
		<METADATA>
	}
	"data": [
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
		log.Printf("[-] BuildEpoch in metadata will be ignored")
	}

	// DatabaseType
	var databaseType string
	if metadata["DatabaseType"] == nil {
		log.Fatalf("[!] DatabaseType is a required field in metadata")
	} else {
		databaseType = metadata["DatabaseType"].(string)
	}

	// Description
	description := make(map[string]string)
	if metadata["Description"] == nil {
		log.Fatalf("[!] Description is a required field in metadata")
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

		log.Printf("[-] IPVersion is not provided in metadata, defaulting to 6 (An IPv6 database supports both IPv4 and IPv6 lookups)")
	} else {
		// Convert the IPVersion to int
		ipVersion = int(metadata["IPVersion"].(float64))

		if ipVersion != 4 && ipVersion != 6 {
			log.Fatalf("[!] Invalid value for IPVersion in metadata. The supported values are 4 and 6")
		}
	}

	// Languages
	languages := make([]string, 0)
	if metadata["Languages"] == nil {
		// Default to English
		languages = append(languages, "en")

		log.Printf("[-] Languages is not provided in metadata, defaulting to English")
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

		log.Printf("[-] RecordSize is not provided in metadata, defaulting to 28 (The supported values are 24, 28, and 32)")
	} else {
		// Convert the RecordSize to int
		recordSize = int(metadata["RecordSize"].(float64))

		if recordSize != 24 && recordSize != 28 && recordSize != 32 {
			log.Fatalf("[!] Invalid value for RecordSize in metadata. The supported values are 24, 28, and 32")
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
		log.Fatalf("[!] Error initializing MMDB writer options: %v", err)
	}

	writer, err := mmdbwriter.New(*mmdbWriterOptions)
	if err != nil {
		log.Fatalf("[!] Error initializing MMDB writer: %v", err)
	}

	return writer, err
}

func GenerateMMDB(cfg *CmdGenerateConfig) error {

	// Initialize the record position
	var recordPosition int = 0

	// Read the dataset
	var dataset map[string]interface{}
	dataset, err := readDataSet(cfg.InputDataset)
	if err != nil {
		log.Fatalf("[!] Error reading dataset: %v", err)
	}

	// Extract metadata from the dataset
	metadata := dataset["metadata"]

	// Initialize the MMDB writer
	writer, err := initializeMMDBWriter(cfg, metadata.(map[string]interface{}))
	if err != nil {
		log.Fatalf("[!] Error initializing MMDB writer: %v", err)
	}

	// Iterate over the dataset and write to the MMDB
	for _, data := range dataset["data"].([]interface{}) {

		recordPosition++

		// Type assertion to map[string]interface{}
		dataMap := data.(map[string]interface{})

		// Record network
		_, network, err := net.ParseCIDR(dataMap["network"].(string))
		if err != nil {
			log.Fatalf("[!] Invalid network (%s) in the dataset: %v", dataMap["network"].(string), err)
		}

		// Parse dynamic data
		dynamicData, exists := dataMap["record"].(map[string]interface{})
		if !exists {
			log.Fatalf("[!] Error parsing data for record %d (network: %s) - %v", recordPosition, network, err)
		}

		// Convert dynamic data to MMDB type map
		dynamicMmdbData := mmdb.ConvertToMMDBTypeMap(dynamicData)

		if err := writer.Insert(network, dynamicMmdbData); err != nil {
			log.Fatalf("[!] Error inserting record %d (network: %s) - %v", recordPosition, network, err)
		}

		fmt.Printf("\r[+] Inserted %d records", recordPosition)
	}
	fmt.Printf("\r")
	log.Printf("[+] Total records inserted: %d", recordPosition)

	// Write the MMDB database to the output file
	outputFile, err := os.Create(cfg.OutputDatabase)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write the MMDB database to the output file
	log.Println("[+] Writing MMDB database to the output file")
	if _, err = writer.WriteTo(outputFile); err != nil {
		return err
	}
	if err == nil {
		log.Println("[+] MMDB generated successfully")
	}

	return err
}

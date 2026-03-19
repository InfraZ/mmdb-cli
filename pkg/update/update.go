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

package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"

	"github.com/InfraZ/mmdb-cli/internal/files"
	"github.com/InfraZ/mmdb-cli/pkg/mmdb"
)

type CmdUpdateConfig struct {
	InputDatabase  string
	InputDataSet   string
	OutputDatabase string
	Verbose        bool

	DisableIPv4Aliasing     bool
	IncludeReservedNetworks bool
}

func readJsonInput(inputDataSet string) (map[string]interface{}, error) {
	_, err := os.Stat(inputDataSet)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", inputDataSet)
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

func parseInputData(inputDataSet string) ([]map[string]interface{}, map[string]interface{}, string, error) {
	inputData, err := readJsonInput(inputDataSet)
	if err != nil {
		return nil, nil, "", fmt.Errorf("error reading dataset: %w", err)
	}

	datasetInterface, exists := inputData["dataset"]
	if !exists {
		return nil, nil, "", fmt.Errorf("no 'dataset' field found in input data")
	}

	datasetSlice, ok := datasetInterface.([]interface{})
	if !ok {
		return nil, nil, "", fmt.Errorf("dataset field is not an array")
	}

	inputDataDataset := make([]map[string]interface{}, len(datasetSlice))
	for i, item := range datasetSlice {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, nil, "", fmt.Errorf("dataset item %d is not a valid object", i+1)
		}
		inputDataDataset[i] = itemMap
	}

	var inputDataSchema map[string]interface{}
	if schemaInterface, exists := inputData["schema"]; exists {
		if schema, ok := schemaInterface.(map[string]interface{}); ok {
			inputDataSchema = schema
		}
	}

	var inputDataVersion string
	if versionInterface, exists := inputData["version"]; exists {
		if version, ok := versionInterface.(string); ok {
			inputDataVersion = version
		}
	}

	return inputDataDataset, inputDataSchema, inputDataVersion, nil
}

func UpdateMMDB(cfg CmdUpdateConfig) error {

	filesToCheck := []files.FilesListValidation{
		{FilePath: cfg.InputDataSet, ExpectedExtension: ".json", ShouldExist: true},
		{FilePath: cfg.InputDatabase, ExpectedExtension: ".mmdb", ShouldExist: true},
		{FilePath: cfg.OutputDatabase, ExpectedExtension: ".mmdb", ShouldExist: false},
	}

	if err := files.FilesValidation(filesToCheck); err != nil {
		return err
	}

	inputDataDataset, inputDataSchema, inputDataVersion, err := parseInputData(cfg.InputDataSet)
	if err != nil {
		return fmt.Errorf("error parsing input data: %w", err)
	}

	if inputDataVersion != "" {
		if inputDataVersion != "v1" {
			return fmt.Errorf("unsupported version: %s (supported: v1)", inputDataVersion)
		}
		fmt.Printf("[+] Dataset version: %s\n", inputDataVersion)
	}

	var useDefaultSchema bool = true
	if inputDataSchema != nil {
		fmt.Printf("[+] Dataset schema: %v\n", inputDataSchema)
		useDefaultSchema = false
	} else {
		fmt.Println("[-] No schema found in input data, using default schema")
	}

	var updatePosition int

	writer, err := mmdbwriter.Load(cfg.InputDatabase, mmdbwriter.Options{
		DisableIPv4Aliasing:     cfg.DisableIPv4Aliasing,
		IncludeReservedNetworks: cfg.IncludeReservedNetworks,
	})
	if err != nil {
		return fmt.Errorf("failed to load MMDB database: %w", err)
	}

	fmt.Println("[+] Starting update mmdb with dataset")

	for _, updateRequest := range inputDataDataset {
		updatePosition++

		_, networkExists := updateRequest["network"]
		if !networkExists {
			return fmt.Errorf("no 'network' found for record %d", updatePosition)
		}

		_, network, err := net.ParseCIDR(updateRequest["network"].(string))
		if err != nil {
			return fmt.Errorf("error parsing network for record %d (%s) - %w", updatePosition, updateRequest["network"], err)
		}

		_, dataExists := updateRequest["data"]
		if !dataExists {
			return fmt.Errorf("no 'data' found for record %d (network: %s)", updatePosition, network)
		}

		dynamicData, exists := updateRequest["data"].(map[string]interface{})
		if !exists {
			return fmt.Errorf("error parsing data for record %d (network: %s)", updatePosition, network)
		}

		dynamicMmdbData := mmdb.ConvertToMMDBTypeMap(dynamicData, useDefaultSchema, inputDataSchema)

		method, isMethodPresent := updateRequest["method"].(string)
		if !isMethodPresent {
			fmt.Printf("[!] No 'method' found for record %d, defaulting to 'deep_merge'\n", updatePosition)
			method = "deep_merge"
		}

		switch method {
		case "remove":
			if err := writer.InsertFunc(network, inserter.Remove); err != nil {
				return fmt.Errorf("error removing data for record %d (network: %s) - %w", updatePosition, network, err)
			}
		case "replace":
			if err := writer.InsertFunc(network, inserter.ReplaceWith(dynamicMmdbData)); err != nil {
				return fmt.Errorf("error replacing data for record %d (network: %s) - %w", updatePosition, network, err)
			}
		case "top_level_merge":
			if err := writer.InsertFunc(network, inserter.TopLevelMergeWith(dynamicMmdbData)); err != nil {
				return fmt.Errorf("error top level merging data for record %d (network: %s) - %w", updatePosition, network, err)
			}
		case "deep_merge":
			if err := writer.InsertFunc(network, inserter.DeepMergeWith(dynamicMmdbData)); err != nil {
				return fmt.Errorf("error deep merging data for record %d (network: %s) - %w", updatePosition, network, err)
			}
		default:
			return fmt.Errorf("unsupported method '%s' for record %d (supported: remove, replace, top_level_merge, deep_merge)", method, updatePosition)
		}

		if cfg.Verbose {
			fmt.Printf("[+] %d/%d dataset records processed - Data: %v\n", updatePosition, len(inputDataDataset), dynamicMmdbData)
		} else {
			fmt.Printf("\r[+] %d/%d dataset records processed", updatePosition, len(inputDataDataset))
		}
	}

	fmt.Printf("\r[+] %d Dataset records processed\n", updatePosition)

	fmt.Printf("[+] Writing updated MMDB to file")
	outputFile, err := os.Create(cfg.OutputDatabase)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if _, err = writer.WriteTo(outputFile); err != nil {
		return err
	}

	fileSize, err := files.CheckFileSizeMb(cfg.OutputDatabase)
	if err != nil {
		return fmt.Errorf("failed to check output file size: %w", err)
	}
	fmt.Printf("\r[+] %s file size: %.2f MB\n", cfg.OutputDatabase, fileSize)

	fmt.Println("[+] MMDB updated successfully")

	return nil
}

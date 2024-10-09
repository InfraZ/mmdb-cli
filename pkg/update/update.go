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
	"io"
	"log"
	"net"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"

	"github.com/InfraZ/mmdb-cli/pkg/mmdb"
)

type CmdUpdateConfig struct {
	InputDatabase  string
	InputDataSet   string
	OutputDatabase string
}

func readDataSet(inputDataSet string) ([]map[string]interface{}, error) {
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

	var dataset []map[string]interface{}
	if err := json.Unmarshal(byteValue, &dataset); err != nil {
		return nil, err
	}

	return dataset, nil
}

func UpdateMMDB(cfg CmdUpdateConfig) error {

	var dataset []map[string]interface{}
	dataset, err := readDataSet(cfg.InputDataSet)
	if err != nil {
		log.Fatalf("Error reading dataset: %v", err)
	}

	var (
		// updatedCount   int = len(dataset)
		updatePosition int
	)

	writer, err := mmdbwriter.Load(cfg.InputDatabase, mmdbwriter.Options{})
	if err != nil {
		log.Fatal(err)
	}

	for _, updateRequest := range dataset {
		updatePosition++

		// Check if network is present
		_, networkExists := updateRequest["network"]
		if !networkExists {
			log.Fatalf("[!] No 'network' found for record %d", updatePosition)
		}

		_, network, err := net.ParseCIDR(updateRequest["network"].(string))
		if err != nil {
			log.Fatalf("[!] Error parsing network for record %d (%s) - %v", updatePosition, updateRequest["network"], err)
		}

		// Check if data is present
		_, dataExists := updateRequest["data"]
		if !dataExists {
			log.Fatalf("[!] No 'data' found for record %d (network: %s)", updatePosition, network)
		}

		// Parse dynamic data
		dynamicData, exists := updateRequest["data"].(map[string]interface{})
		if !exists {
			log.Fatalf("[!] Error parsing data for record %d (network: %s) - %v", updatePosition, network, err)
		}

		dynamicMmdbData := mmdb.ConvertToMMDBTypeMap(dynamicData)

		// Switch to select the type of update
		method, isMethodPresent := updateRequest["method"].(string)
		if !isMethodPresent {
			log.Printf("[!] No 'method' found for record %d, defaulting to 'deep_merge'", updatePosition)
			method = "deep_merge"
		}

		switch method {
		case "remove":
			if err := writer.InsertFunc(network, inserter.Remove); err != nil {
				log.Fatalf("[!] Error removing data for record %d (network: %s) - %v", updatePosition, network, err)
			}
		case "replace":
			// Replace existing data with dynamic data
			if err := writer.InsertFunc(network, inserter.ReplaceWith(dynamicMmdbData)); err != nil {
				log.Fatalf("[!] Error replacing data for record %d (network: %s) - %v", updatePosition, network, err)
			}
		case "top_level_merge":
			// Merge top-level keys and values from new data
			if err := writer.InsertFunc(network, inserter.TopLevelMergeWith(dynamicMmdbData)); err != nil {
				log.Fatalf("[!] Error top level merging data for record %d (network: %s) - %v", updatePosition, network, err)
			}
		case "deep_merge":
			// Deep merge dynamic data with existing data
			if err := writer.InsertFunc(network, inserter.DeepMergeWith(dynamicMmdbData)); err != nil {
				log.Fatalf("[!] Error deep merging data for record %d (network: %s) - %v", updatePosition, network, err)
			}
		default:
			log.Fatalf("[!] Unsupported method %s for record %d (supported: remove, replace, top_level_merge, deep_merge)", method, updatePosition)
		}

	}

	outputFile, err := os.Create(cfg.OutputDatabase)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if _, err = writer.WriteTo(outputFile); err != nil {
		return err
	}

	return nil
}

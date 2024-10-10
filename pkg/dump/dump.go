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

package dump

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/oschwald/maxminddb-golang"
)

type CmdDumpConfig struct {
	InputDatabase string
	OutputFile    string
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

func DumpMMMDB(cfg *CmdDumpConfig) error {
	// Open the MMDB database file
	db, err := maxminddb.Open(cfg.InputDatabase)
	if err != nil {
		log.Fatalf("[!] Failed to open database: %s - %v", cfg.InputDatabase, err)
	}
	defer db.Close()

	// Check output file extension
	if len(cfg.OutputFile) < 5 || cfg.OutputFile[len(cfg.OutputFile)-5:] != ".json" {
		return fmt.Errorf("[!] Output file must have a .json extension")
	}

	// Create the output file
	outputFile, err := os.Create(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("[!] Failed to create output file: %s - %v", cfg.OutputFile, err)
	}
	defer outputFile.Close()

	// Prepare output data
	outputData := make(map[string]interface{})
	outputData["schema"] = "v1"
	outputData["metadata"] = db.Metadata

	// dump counter
	var dumpPosition int = 0

	// Init output data
	outputData["data"] = make([]map[string]interface{}, 0)

	// Get all available networks
	availableNetworks := db.Networks()

	// Iterate over all available networks
	for availableNetworks.Next() {

		dumpPosition++
		data := make(map[string]interface{})
		record := make(map[string]interface{})

		subnet, err := availableNetworks.Network(&record)
		if err != nil {
			return fmt.Errorf("failed to get record for next subnet: %w", err)
		}

		data["network"] = subnet.String()
		data["record"] = record

		outputData["data"] = append(outputData["data"].([]map[string]interface{}), data)

		fmt.Printf("\r[+] Dumped %d records", dumpPosition)

	}

	// Write the output data to the file
	encoder := json.NewEncoder(outputFile)
	err = encoder.Encode(outputData)
	if err != nil {
		return fmt.Errorf("[!] Failed to write output data to file: %s - %v", cfg.OutputFile, err)
	}

	fmt.Println("\n[+] MMDB dumped successfully")

	return nil
}

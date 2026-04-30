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
	"os"

	"github.com/InfraZ/mmdb-cli/internal/files"
	"github.com/InfraZ/mmdb-cli/pkg/jsonpath"
	"github.com/oschwald/maxminddb-golang"
)

type CmdDumpConfig struct {
	InputDatabase string
	OutputFile    string
	Verbose       bool
	JSONPath      string
}

/*
Structure of the dumped JSON dataset:

	{
		"version": "v1",
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
func DumpMMMDB(cfg *CmdDumpConfig) error {

	filesToCheck := []files.FilesListValidation{
		{FilePath: cfg.InputDatabase, ExpectedExtension: ".mmdb", ShouldExist: true},
		{FilePath: cfg.OutputFile, ExpectedExtension: ".json", ShouldExist: false},
	}

	if err := files.FilesValidation(filesToCheck); err != nil {
		return err
	}

	db, err := maxminddb.Open(cfg.InputDatabase)
	if err != nil {
		return fmt.Errorf("failed to open database: %s - %w", cfg.InputDatabase, err)
	}
	defer db.Close()

	if len(cfg.OutputFile) < 5 || cfg.OutputFile[len(cfg.OutputFile)-5:] != ".json" {
		return fmt.Errorf("output file must have a .json extension")
	}

	outputFile, err := os.Create(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %s - %w", cfg.OutputFile, err)
	}
	defer outputFile.Close()

	fmt.Printf("[+] Start dumping %s to %s\n", cfg.InputDatabase, cfg.OutputFile)

	if cfg.JSONPath != "" {
		if err := jsonpath.ValidateExpression(cfg.JSONPath); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	metadataJSON, err := json.Marshal(db.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	if _, err := fmt.Fprintf(outputFile, `{"version":"v1","metadata":%s,"dataset":[`, metadataJSON); err != nil {
		return fmt.Errorf("failed to write output header: %w", err)
	}

	encoder := json.NewEncoder(outputFile)
	var readPosition int
	var dumpPosition int
	firstRecord := true

	availableNetworks := db.Networks(
		maxminddb.SkipAliasedNetworks,
	)

	for availableNetworks.Next() {
		readPosition++
		record := make(map[string]interface{})

		subnet, err := availableNetworks.Network(&record)
		if err != nil {
			return fmt.Errorf("failed to get record for next subnet: %w", err)
		}

		if cfg.JSONPath != "" {
			match, err := jsonpath.MatchesRecord(cfg.JSONPath, record)
			if err != nil {
				return fmt.Errorf("failed to evaluate JSONPath for network %s: %w", subnet.String(), err)
			}
			if !match {
				if !cfg.Verbose {
					fmt.Printf("\r[-] Read records: %d, Matched records: %d", readPosition, dumpPosition)
				}
				continue
			}
		}

		dumpPosition++

		if !firstRecord {
			if _, err := outputFile.WriteString(","); err != nil {
				return fmt.Errorf("failed to write record separator: %w", err)
			}
		}
		firstRecord = false

		data := map[string]interface{}{
			"network": subnet.String(),
			"record":  record,
		}
		if err := encoder.Encode(data); err != nil {
			return fmt.Errorf("failed to encode record for network %s: %w", subnet.String(), err)
		}

		if cfg.Verbose {
			fmt.Printf("[-] Dumping record %d for network %s - data: %v\n", dumpPosition, subnet.String(), record)
		} else if cfg.JSONPath != "" {
			fmt.Printf("\r[-] Read records: %d, Matched records: %d", readPosition, dumpPosition)
		} else {
			fmt.Printf("\r[-] Dumped records: %d", dumpPosition)
		}
	}

	if _, err := outputFile.WriteString("]}"); err != nil {
		return fmt.Errorf("failed to write output footer: %w", err)
	}

	if cfg.JSONPath != "" {
		fmt.Printf("\r[+] Read %d records, matched %d records\n", readPosition, dumpPosition)
	} else {
		fmt.Printf("\r[+] Total %d records dumped successfully\n", dumpPosition)
	}

	outputFileStat, err := outputFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get output file stats: %s - %w", cfg.OutputFile, err)
	}

	outputFileSizeMB := float64(outputFileStat.Size()) / 1024 / 1024
	fmt.Printf("[+] %s file created with size: %.2f MB\n", cfg.OutputFile, outputFileSizeMB)

	fmt.Println("[+] MMDB Dumped successfully")

	return nil
}

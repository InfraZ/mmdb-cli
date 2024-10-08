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

package metadata

import (
	"encoding/json"
	"log"

	"github.com/oschwald/maxminddb-golang"
)

type CmdMetadataConfig struct {
	InputFile string
}

type DatabaseMetadata struct {
	Description              map[string]string `json:"description"`
	DatabaseType             string            `json:"database_type"`
	Languages                []string          `json:"languages"`
	BinaryFormatMajorVersion uint              `json:"binary_format_major_version"`
	BinaryFormatMinorVersion uint              `json:"binary_format_minor_version"`
	BuildEpoch               uint              `json:"build_epoch"`
	IPVersion                uint              `json:"ip_version"`
	NodeCount                uint              `json:"node_count"`
	RecordSize               uint              `json:"record_size"`
}

func MetadataMMDB(cfg CmdMetadataConfig) ([]byte, error) {

	db, err := maxminddb.Open(cfg.InputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mmdbMetadata := db.Metadata

	databaseMetadata := DatabaseMetadata{
		Description:              mmdbMetadata.Description,
		DatabaseType:             mmdbMetadata.DatabaseType,
		Languages:                mmdbMetadata.Languages,
		BinaryFormatMajorVersion: mmdbMetadata.BinaryFormatMajorVersion,
		BinaryFormatMinorVersion: mmdbMetadata.BinaryFormatMinorVersion,
		BuildEpoch:               mmdbMetadata.BuildEpoch,
		IPVersion:                mmdbMetadata.IPVersion,
		NodeCount:                mmdbMetadata.NodeCount,
		RecordSize:               mmdbMetadata.RecordSize,
	}

	jsonDatabaseMetadata, err := json.Marshal(databaseMetadata)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(jsonDatabaseMetadata))

	return jsonDatabaseMetadata, err
}

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

package mmdb

import (
	"log"

	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func ConvertToMMDBTypeMap(data map[string]interface{}) mmdbtype.Map {
	mmdbMap := mmdbtype.Map{}
	for key, value := range data {
		mmdbKey := mmdbtype.String(key)
		switch mmdbValue := value.(type) {
		case string:
			mmdbMap[mmdbKey] = mmdbtype.String(mmdbValue)
		case bool:
			mmdbMap[mmdbKey] = mmdbtype.Bool(mmdbValue)
		case float64:
			mmdbMap[mmdbKey] = mmdbtype.Float64(mmdbValue)
		case int:
			mmdbMap[mmdbKey] = mmdbtype.Int32(mmdbValue)
		case map[string]interface{}:
			// Recursively convert nested maps
			mmdbMap[mmdbKey] = ConvertToMMDBTypeMap(mmdbValue)
		default:
			log.Printf("Unsupported data type for key %s", key)
		}
	}
	return mmdbMap
}

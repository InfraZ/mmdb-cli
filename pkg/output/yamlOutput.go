// Copyright 2024 The MMDB CLI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func YamlOutput(data []byte, options OutputOptions) error {
	var jsonData interface{}

	// Unmarshal the JSON data into a Go data structure
	err := yaml.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}

	// Marshal the Go data structure into YAML format
	yamlData, err := yaml.Marshal(jsonData)
	if err != nil {
		return err
	}

	// Print the YAML data
	fmt.Println(string(yamlData))

	return nil
}

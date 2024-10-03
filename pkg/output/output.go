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

import "fmt"

type OutputOptions struct {
	Format     string
	JsonPretty bool
}

func Output(byteData []byte, options OutputOptions) error {

	if (options.JsonPretty) && (options.Format != "json") {
		return fmt.Errorf("Pretty print is only supported for JSON output")
	}

	if options.Format == "json-pretty" {
		options.Format = "json"
		options.JsonPretty = true
	}

	switch options.Format {
	case "json":
		return JsonOutput(byteData, options)
	case "yaml":
		return YamlOutput(byteData, options)
	default:
		return fmt.Errorf("Unsupported output format: %s", options.Format)
	}

}

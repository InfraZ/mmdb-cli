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
	"os"
	"testing"
)

func TestMetadataMMDB(t *testing.T) {
	// Skip test if running on GitHub Actions as the test requires the GeoLite2-Country.mmdb file
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test; running on GitHub Actions")
	}

	expected := `{"description":{"en":"GeoLite2 Country database"},"database_type":"GeoLite2-Country","languages":["de","en","es","fr","ja","pt-BR","ru","zh-CN"],"binary_format_major_version":2,"binary_format_minor_version":0,"build_epoch":1726249063,"ip_version":6,"node_count":1234259,"record_size":24}`
	testMMDBFile := "../../output/GeoLite2-Country.mmdb"

	result, err := MetadataMMDB(CmdMetadataConfig{InputFile: testMMDBFile})
	jsonResult := string(result)
	if err != nil {
		t.Errorf("MetadataMMDB() error = %v; want nil", err)
	}

	if string(result) != expected {
		t.Errorf("MetadataMMDB() = %v; want %v", jsonResult, expected)
	}
}

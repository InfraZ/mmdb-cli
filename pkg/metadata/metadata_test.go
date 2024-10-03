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
package metadata

import "testing"

func TestMetadataMMDB(t *testing.T) {
	expected := "" // TODO: Define expected value
	testMMDBFile := "../output/GeoLite2-Country.mmdb"

	result, err := MetadataMMDB(CmdMetadataConfig{InputFile: testMMDBFile})
	if err != nil {
		t.Errorf("MetadataMMDB() error = %v; want nil", err)
	}

	if result != nil {
		t.Errorf("MetadataMMDB() = %v; want %v", result, expected)
	}
}

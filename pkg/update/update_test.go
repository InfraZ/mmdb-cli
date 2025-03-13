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
	"os"
	"testing"
)

func TestReadDataSet(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "dataset-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write sample JSON data to the temporary file
	sampleData := `[{"network": "192.168.1.0/24", "data": {"country": "US"}}]`
	if _, err := tmpFile.Write([]byte(sampleData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	invalidJSONFile := createTempFile(t, "invalid json")
	defer os.Remove(invalidJSONFile)

	tests := []struct {
		name          string
		inputDataSet  string
		expectedError bool
		expectedData  []map[string]interface{}
	}{
		{
			name:          "Valid dataset file",
			inputDataSet:  tmpFile.Name(),
			expectedError: false,
			expectedData: []map[string]interface{}{
				{
					"network": "192.168.1.0/24",
					"data":    map[string]interface{}{"country": "US"},
				},
			},
		},
		{
			name:          "Invalid JSON format",
			inputDataSet:  invalidJSONFile,
			expectedError: true,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := readDataSet(tt.inputDataSet)
			if (err != nil) != tt.expectedError {
				t.Errorf("readDataSet() error = %v, expectedError %v", err, tt.expectedError)
			}
			if !tt.expectedError && data != nil {
				// Compare the actual data with expected data
				if len(data) != len(tt.expectedData) {
					t.Errorf("readDataSet() got %v entries, expected %v entries", len(data), len(tt.expectedData))
					return
				}
				// Basic structure validation for the first entry
				if len(data) > 0 {
					if data[0]["network"] != tt.expectedData[0]["network"] {
						t.Errorf("readDataSet() network = %v, expected %v", data[0]["network"], tt.expectedData[0]["network"])
					}
				}
			}
		})
	}
}

func createTempFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "dataset-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

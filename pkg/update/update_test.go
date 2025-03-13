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
	tmpFile.Close()

	tests := []struct {
		name          string
		inputDataSet  string
		expectedError bool
	}{
		{
			name:          "Valid dataset file",
			inputDataSet:  tmpFile.Name(),
			expectedError: false,
		},
		{
			name:          "Non-existent dataset file",
			inputDataSet:  "non_existent_file.json",
			expectedError: true,
		},
		{
			name:          "Invalid JSON format",
			inputDataSet:  createTempFile(t, "invalid json"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := readDataSet(tt.inputDataSet)
			if (err != nil) != tt.expectedError {
				t.Errorf("readDataSet() error = %v, expectedError %v", err, tt.expectedError)
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

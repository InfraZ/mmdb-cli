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
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/oschwald/maxminddb-golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testMMDB = "../../test/inspect.mmdb"

func writeTestFile(t *testing.T, dir, filename, content string) string {
	t.Helper()
	path := filepath.Join(dir, filename)
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func TestParseInputData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		content       string
		wantErr       bool
		errContains   string
		expectedLen   int
		expectSchema  bool
		expectVersion string
	}{
		{
			name:        "valid dataset",
			content:     `{"dataset": [{"network": "192.168.1.0/24", "data": {"country": "US"}}]}`,
			expectedLen: 1,
		},
		{
			name:        "valid with version and schema",
			content:     `{"version": "v1", "schema": {"country": "string"}, "dataset": [{"network": "1.0.0.0/8", "data": {"country": "AU"}}]}`,
			expectedLen: 1,
			expectSchema: true,
			expectVersion: "v1",
		},
		{
			name:        "empty dataset array",
			content:     `{"dataset": []}`,
			expectedLen: 0,
		},
		{
			name:        "missing dataset field",
			content:     `{"version": "v1"}`,
			wantErr:     true,
			errContains: "no 'dataset' field",
		},
		{
			name:        "dataset is not an array",
			content:     `{"dataset": "not_an_array"}`,
			wantErr:     true,
			errContains: "dataset field is not an array",
		},
		{
			name:        "dataset item is not an object",
			content:     `{"dataset": ["not_an_object"]}`,
			wantErr:     true,
			errContains: "dataset item 1 is not a valid object",
		},
		{
			name:        "invalid JSON",
			content:     `{invalid`,
			wantErr:     true,
			errContains: "error reading dataset",
		},
		{
			name:          "version without schema",
			content:       `{"version": "v1", "dataset": [{"network": "1.0.0.0/8", "data": {"test": "value"}}]}`,
			expectedLen:   1,
			expectVersion: "v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := writeTestFile(t, dir, "dataset.json", tt.content)

			data, schema, version, err := parseInputData(path)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}
			require.NoError(t, err)
			assert.Len(t, data, tt.expectedLen)

			if tt.expectSchema {
				assert.NotNil(t, schema)
			}
			if tt.expectVersion != "" {
				assert.Equal(t, tt.expectVersion, version)
			}
		})
	}

	t.Run("non-existent file", func(t *testing.T) {
		t.Parallel()
		_, _, _, err := parseInputData("/nonexistent/dataset.json")
		assert.Error(t, err)
	})
}

func TestReadJsonInput(t *testing.T) {
	t.Parallel()

	t.Run("valid JSON file", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := writeTestFile(t, dir, "data.json", `{"key": "value"}`)

		result, err := readJsonInput(path)
		require.NoError(t, err)
		assert.Equal(t, "value", result["key"])
	})

	t.Run("invalid JSON file", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := writeTestFile(t, dir, "bad.json", `{invalid}`)

		_, err := readJsonInput(path)
		assert.Error(t, err)
	})

	t.Run("non-existent file", func(t *testing.T) {
		t.Parallel()
		_, err := readJsonInput("/nonexistent/file.json")
		assert.Error(t, err)
	})
}

func TestUpdateMMDB(t *testing.T) {
	tests := []struct {
		name    string
		dataset string
		wantErr bool
		verify  func(t *testing.T, outputPath string)
	}{
		{
			name: "deep_merge update",
			dataset: `{
				"dataset": [
					{
						"network": "1.1.1.1/32",
						"method": "deep_merge",
						"data": {"extra_field": "new_value"}
					}
				]
			}`,
			verify: func(t *testing.T, outputPath string) {
				t.Helper()
				db, err := maxminddb.Open(outputPath)
				require.NoError(t, err)
				defer db.Close()

				var record map[string]interface{}
				err = db.Lookup(net.ParseIP("1.1.1.1"), &record)
				require.NoError(t, err)
				assert.Contains(t, record, "extra_field")
				assert.Contains(t, record, "registered_country")
			},
		},
		{
			name: "replace update",
			dataset: `{
				"dataset": [
					{
						"network": "1.1.1.1/32",
						"method": "replace",
						"data": {"replaced": "yes"}
					}
				]
			}`,
			verify: func(t *testing.T, outputPath string) {
				t.Helper()
				db, err := maxminddb.Open(outputPath)
				require.NoError(t, err)
				defer db.Close()

				var record map[string]interface{}
				err = db.Lookup(net.ParseIP("1.1.1.1"), &record)
				require.NoError(t, err)
				assert.Contains(t, record, "replaced")
				assert.NotContains(t, record, "registered_country")
			},
		},
		{
			name: "top_level_merge update",
			dataset: `{
				"dataset": [
					{
						"network": "1.1.1.1/32",
						"method": "top_level_merge",
						"data": {"top_level_new": "merged"}
					}
				]
			}`,
			verify: func(t *testing.T, outputPath string) {
				t.Helper()
				db, err := maxminddb.Open(outputPath)
				require.NoError(t, err)
				defer db.Close()

				var record map[string]interface{}
				err = db.Lookup(net.ParseIP("1.1.1.1"), &record)
				require.NoError(t, err)
				assert.Contains(t, record, "top_level_new")
			},
		},
		{
			name: "remove update",
			dataset: `{
				"dataset": [
					{
						"network": "1.1.1.1/32",
						"method": "remove",
						"data": {}
					}
				]
			}`,
			verify: func(t *testing.T, outputPath string) {
				t.Helper()
				db, err := maxminddb.Open(outputPath)
				require.NoError(t, err)
				defer db.Close()

				var record map[string]interface{}
				err = db.Lookup(net.ParseIP("1.1.1.1"), &record)
				require.NoError(t, err)
				assert.Empty(t, record)
			},
		},
		{
			name: "default method (deep_merge)",
			dataset: `{
				"dataset": [
					{
						"network": "1.1.1.1/32",
						"data": {"default_merge": "value"}
					}
				]
			}`,
			verify: func(t *testing.T, outputPath string) {
				t.Helper()
				db, err := maxminddb.Open(outputPath)
				require.NoError(t, err)
				defer db.Close()

				var record map[string]interface{}
				err = db.Lookup(net.ParseIP("1.1.1.1"), &record)
				require.NoError(t, err)
				assert.Contains(t, record, "default_merge")
			},
		},
		{
			name: "unsupported method",
			dataset: `{
				"dataset": [
					{
						"network": "1.1.1.1/32",
						"method": "invalid_method",
						"data": {"key": "value"}
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "missing network field",
			dataset: `{
				"dataset": [
					{
						"data": {"key": "value"}
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "missing data field",
			dataset: `{
				"dataset": [
					{
						"network": "1.1.1.1/32"
					}
				]
			}`,
			wantErr: true,
		},
		{
			name: "unsupported version",
			dataset: `{
				"version": "v99",
				"dataset": [
					{
						"network": "1.1.1.1/32",
						"data": {"key": "value"}
					}
				]
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			datasetPath := writeTestFile(t, dir, "update.json", tt.dataset)
			outputPath := filepath.Join(dir, "updated.mmdb")

			cfg := CmdUpdateConfig{
				InputDatabase:  testMMDB,
				InputDataSet:   datasetPath,
				OutputDatabase: outputPath,
			}

			err := UpdateMMDB(cfg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			_, statErr := os.Stat(outputPath)
			require.NoError(t, statErr)

			if tt.verify != nil {
				tt.verify(t, outputPath)
			}
		})
	}
}

func TestUpdateMMDBInvalidInputDB(t *testing.T) {
	dir := t.TempDir()
	datasetPath := writeTestFile(t, dir, "update.json", `{"dataset":[{"network":"1.0.0.0/8","data":{"k":"v"}}]}`)

	cfg := CmdUpdateConfig{
		InputDatabase:  "/nonexistent/file.mmdb",
		InputDataSet:   datasetPath,
		OutputDatabase: filepath.Join(dir, "out.mmdb"),
	}
	err := UpdateMMDB(cfg)
	assert.Error(t, err)
}

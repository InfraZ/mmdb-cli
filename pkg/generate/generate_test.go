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

package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oschwald/maxminddb-golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTestJSON(t *testing.T, dir, filename, content string) string {
	t.Helper()
	path := filepath.Join(dir, filename)
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func TestReadDataSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid JSON",
			content: `{"metadata":{"DatabaseType":"Test"},"dataset":[]}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			content: `{invalid`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := writeTestJSON(t, dir, "data.json", tt.content)

			result, err := readDataSet(path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, result)
		})
	}

	t.Run("non-existent file", func(t *testing.T) {
		t.Parallel()
		_, err := readDataSet("/nonexistent/file.json")
		assert.Error(t, err)
	})
}

func TestMmdbWriterOptions(t *testing.T) {
	t.Parallel()

	baseCfg := &CmdGenerateConfig{}

	tests := []struct {
		name     string
		metadata map[string]interface{}
		wantErr  bool
		verify   func(t *testing.T, err error)
	}{
		{
			name: "all fields present",
			metadata: map[string]interface{}{
				"DatabaseType": "GeoIP2-City",
				"Description":  map[string]interface{}{"en": "Test DB"},
				"IPVersion":    float64(6),
				"Languages":    []interface{}{"en", "de"},
				"RecordSize":   float64(28),
			},
			wantErr: false,
		},
		{
			name: "missing optional fields - defaults applied",
			metadata: map[string]interface{}{
				"DatabaseType": "GeoIP2-City",
				"Description":  map[string]interface{}{"en": "Test DB"},
			},
			wantErr: false,
		},
		{
			name: "missing required DatabaseType",
			metadata: map[string]interface{}{
				"Description": map[string]interface{}{"en": "Test DB"},
			},
			wantErr: true,
		},
		{
			name: "missing required Description",
			metadata: map[string]interface{}{
				"DatabaseType": "GeoIP2-City",
			},
			wantErr: true,
		},
		{
			name: "invalid IPVersion",
			metadata: map[string]interface{}{
				"DatabaseType": "GeoIP2-City",
				"Description":  map[string]interface{}{"en": "Test DB"},
				"IPVersion":    float64(5),
			},
			wantErr: true,
		},
		{
			name: "invalid RecordSize",
			metadata: map[string]interface{}{
				"DatabaseType": "GeoIP2-City",
				"Description":  map[string]interface{}{"en": "Test DB"},
				"RecordSize":   float64(16),
			},
			wantErr: true,
		},
		{
			name: "valid RecordSize 24",
			metadata: map[string]interface{}{
				"DatabaseType": "GeoIP2-City",
				"Description":  map[string]interface{}{"en": "Test DB"},
				"RecordSize":   float64(24),
			},
			wantErr: false,
		},
		{
			name: "valid RecordSize 32",
			metadata: map[string]interface{}{
				"DatabaseType": "GeoIP2-City",
				"Description":  map[string]interface{}{"en": "Test DB"},
				"RecordSize":   float64(32),
			},
			wantErr: false,
		},
		{
			name: "valid IPVersion 4",
			metadata: map[string]interface{}{
				"DatabaseType": "GeoIP2-City",
				"Description":  map[string]interface{}{"en": "Test DB"},
				"IPVersion":    float64(4),
			},
			wantErr: false,
		},
		{
			name: "BuildEpoch is ignored",
			metadata: map[string]interface{}{
				"DatabaseType": "GeoIP2-City",
				"Description":  map[string]interface{}{"en": "Test DB"},
				"BuildEpoch":   float64(1234567890),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opts, err := mmdbWriterOptions(baseCfg, tt.metadata)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, opts)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, opts)
		})
	}
}

func TestGenerateMMDB(t *testing.T) {
	t.Parallel()

	t.Run("end-to-end generation", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		inputJSON := `{
			"version": "v1",
			"metadata": {
				"DatabaseType": "Test-DB",
				"Description": {"en": "Test Database"},
				"IPVersion": 6,
				"Languages": ["en"],
				"RecordSize": 24
			},
			"dataset": [
				{
					"network": "1.1.1.0/24",
					"record": {
						"country": "AU",
						"city": "Sydney"
					}
				},
				{
					"network": "8.8.8.0/24",
					"record": {
						"country": "US",
						"city": "Mountain View"
					}
				}
			]
		}`
		inputPath := writeTestJSON(t, dir, "input.json", inputJSON)
		outputPath := filepath.Join(dir, "output.mmdb")

		cfg := &CmdGenerateConfig{
			InputDataset:   inputPath,
			OutputDatabase: outputPath,
		}

		err := GenerateMMDB(cfg)
		require.NoError(t, err)

		_, statErr := os.Stat(outputPath)
		require.NoError(t, statErr)

		db, openErr := maxminddb.Open(outputPath)
		require.NoError(t, openErr)
		defer db.Close()

		assert.NoError(t, db.Verify())
		assert.Equal(t, "Test-DB", db.Metadata.DatabaseType)
	})

	t.Run("generation with schema", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		inputJSON := `{
			"version": "v1",
			"schema": {
				"asn": "uint32",
				"name": "string"
			},
			"metadata": {
				"DatabaseType": "ASN-DB",
				"Description": {"en": "ASN Database"}
			},
			"dataset": [
				{
					"network": "1.0.0.0/24",
					"record": {
						"asn": 13335,
						"name": "Cloudflare"
					}
				}
			]
		}`
		inputPath := writeTestJSON(t, dir, "input.json", inputJSON)
		outputPath := filepath.Join(dir, "output.mmdb")

		cfg := &CmdGenerateConfig{
			InputDataset:   inputPath,
			OutputDatabase: outputPath,
		}

		err := GenerateMMDB(cfg)
		require.NoError(t, err)

		db, openErr := maxminddb.Open(outputPath)
		require.NoError(t, openErr)
		defer db.Close()
		assert.NoError(t, db.Verify())
	})

	t.Run("non-existent input file", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		cfg := &CmdGenerateConfig{
			InputDataset:   "/nonexistent/file.json",
			OutputDatabase: filepath.Join(dir, "output.mmdb"),
		}
		err := GenerateMMDB(cfg)
		assert.Error(t, err)
	})

	t.Run("invalid JSON input", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		inputPath := writeTestJSON(t, dir, "bad.json", `{invalid json}`)
		cfg := &CmdGenerateConfig{
			InputDataset:   inputPath,
			OutputDatabase: filepath.Join(dir, "output.mmdb"),
		}
		err := GenerateMMDB(cfg)
		assert.Error(t, err)
	})

	t.Run("unsupported version", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		inputJSON := `{
			"version": "v2",
			"metadata": {
				"DatabaseType": "Test-DB",
				"Description": {"en": "Test"}
			},
			"dataset": []
		}`
		inputPath := writeTestJSON(t, dir, "input.json", inputJSON)
		cfg := &CmdGenerateConfig{
			InputDataset:   inputPath,
			OutputDatabase: filepath.Join(dir, "output.mmdb"),
		}
		err := GenerateMMDB(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported dataset version")
	})

	t.Run("verbose mode", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		inputJSON := `{
			"metadata": {
				"DatabaseType": "Test-DB",
				"Description": {"en": "Test Database"}
			},
			"dataset": [
				{
					"network": "203.0.113.0/24",
					"record": {"name": "example"}
				}
			]
		}`
		inputPath := writeTestJSON(t, dir, "input.json", inputJSON)
		outputPath := filepath.Join(dir, "output.mmdb")

		cfg := &CmdGenerateConfig{
			InputDataset:            inputPath,
			OutputDatabase:          outputPath,
			Verbose:                 true,
			IncludeReservedNetworks: true,
		}

		err := GenerateMMDB(cfg)
		require.NoError(t, err)

		db, openErr := maxminddb.Open(outputPath)
		require.NoError(t, openErr)
		defer db.Close()
		assert.NoError(t, db.Verify())
	})
}

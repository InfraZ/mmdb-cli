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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testMMDB = "../../test/metadata.mmdb"

func TestMetadataMMDB(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     CmdMetadataConfig
		wantErr bool
		verify  func(t *testing.T, result []byte)
	}{
		{
			name: "valid MMDB file",
			cfg:  CmdMetadataConfig{InputFile: testMMDB},
			verify: func(t *testing.T, result []byte) {
				t.Helper()
				var meta DatabaseMetadata
				require.NoError(t, json.Unmarshal(result, &meta))

				assert.Equal(t, "Metadata Test", meta.DatabaseType)
				assert.Equal(t, map[string]string{"en": "MMDB CLI Metadata Test"}, meta.Description)
				assert.Contains(t, meta.Languages, "en")
				assert.Equal(t, uint(2), meta.BinaryFormatMajorVersion)
				assert.Equal(t, uint(0), meta.BinaryFormatMinorVersion)
				assert.Equal(t, uint(6), meta.IPVersion)
				assert.Equal(t, uint(24), meta.RecordSize)
				assert.Greater(t, meta.NodeCount, uint(0))
				assert.Greater(t, meta.BuildEpoch, uint(0))
			},
		},
		{
			name:    "non-existent file",
			cfg:     CmdMetadataConfig{InputFile: "/nonexistent/file.mmdb"},
			wantErr: true,
		},
		{
			name:    "empty path",
			cfg:     CmdMetadataConfig{InputFile: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := MetadataMMDB(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, result)
			if tt.verify != nil {
				tt.verify(t, result)
			}
		})
	}
}

func TestMetadataMMDBJsonStructure(t *testing.T) {
	t.Parallel()

	result, err := MetadataMMDB(CmdMetadataConfig{InputFile: testMMDB})
	require.NoError(t, err)

	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(result, &raw))

	expectedFields := []string{
		"description", "database_type", "languages",
		"binary_format_major_version", "binary_format_minor_version",
		"build_epoch", "ip_version", "node_count", "record_size",
	}
	for _, field := range expectedFields {
		assert.Contains(t, raw, field, "JSON output should contain field %q", field)
	}
}

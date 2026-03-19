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

package dump

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testMMDB = "../../test/inspect.mmdb"

func TestDumpMMMDB(t *testing.T) {
	tests := []struct {
		name    string
		cfg     func(t *testing.T) *CmdDumpConfig
		wantErr bool
		verify  func(t *testing.T, cfg *CmdDumpConfig)
	}{
		{
			name: "successful dump",
			cfg: func(t *testing.T) *CmdDumpConfig {
				t.Helper()
				outFile := filepath.Join(t.TempDir(), "output.json")
				return &CmdDumpConfig{
					InputDatabase: testMMDB,
					OutputFile:    outFile,
				}
			},
			wantErr: false,
			verify: func(t *testing.T, cfg *CmdDumpConfig) {
				t.Helper()
				data, err := os.ReadFile(cfg.OutputFile)
				require.NoError(t, err)

				var result map[string]interface{}
				require.NoError(t, json.Unmarshal(data, &result))

				assert.Equal(t, "v1", result["version"])
				assert.NotNil(t, result["metadata"])
				assert.NotNil(t, result["dataset"])

				dataset, ok := result["dataset"].([]interface{})
				require.True(t, ok)
				assert.Greater(t, len(dataset), 0)

				firstEntry, ok := dataset[0].(map[string]interface{})
				require.True(t, ok)
				assert.Contains(t, firstEntry, "network")
				assert.Contains(t, firstEntry, "record")
			},
		},
		{
			name: "successful dump with verbose",
			cfg: func(t *testing.T) *CmdDumpConfig {
				t.Helper()
				outFile := filepath.Join(t.TempDir(), "output.json")
				return &CmdDumpConfig{
					InputDatabase: testMMDB,
					OutputFile:    outFile,
					Verbose:       true,
				}
			},
			wantErr: false,
			verify: func(t *testing.T, cfg *CmdDumpConfig) {
				t.Helper()
				data, err := os.ReadFile(cfg.OutputFile)
				require.NoError(t, err)

				var result map[string]interface{}
				require.NoError(t, json.Unmarshal(data, &result))
				assert.Equal(t, "v1", result["version"])
			},
		},
		{
			name: "dump with JSONPath filter",
			cfg: func(t *testing.T) *CmdDumpConfig {
				t.Helper()
				outFile := filepath.Join(t.TempDir(), "filtered.json")
				return &CmdDumpConfig{
					InputDatabase: testMMDB,
					OutputFile:    outFile,
					JSONPath:      `{[?(@.registered_country.iso_code=="AU")]}`,
				}
			},
			wantErr: false,
			verify: func(t *testing.T, cfg *CmdDumpConfig) {
				t.Helper()
				data, err := os.ReadFile(cfg.OutputFile)
				require.NoError(t, err)

				var result map[string]interface{}
				require.NoError(t, json.Unmarshal(data, &result))

				dataset, ok := result["dataset"].([]interface{})
				require.True(t, ok)
				assert.Greater(t, len(dataset), 0)
			},
		},
		{
			name: "invalid input path",
			cfg: func(t *testing.T) *CmdDumpConfig {
				t.Helper()
				outFile := filepath.Join(t.TempDir(), "output.json")
				return &CmdDumpConfig{
					InputDatabase: "/nonexistent/path.mmdb",
					OutputFile:    outFile,
				}
			},
			wantErr: true,
		},
		{
			name: "invalid output extension",
			cfg: func(t *testing.T) *CmdDumpConfig {
				t.Helper()
				outFile := filepath.Join(t.TempDir(), "output.txt")
				return &CmdDumpConfig{
					InputDatabase: testMMDB,
					OutputFile:    outFile,
				}
			},
			wantErr: true,
		},
		{
			name: "invalid JSONPath expression",
			cfg: func(t *testing.T) *CmdDumpConfig {
				t.Helper()
				outFile := filepath.Join(t.TempDir(), "output.json")
				return &CmdDumpConfig{
					InputDatabase: testMMDB,
					OutputFile:    outFile,
					JSONPath:      "{[?(@.field==}",
				}
			},
			wantErr: true,
		},
		{
			name: "dump with non-matching JSONPath filter",
			cfg: func(t *testing.T) *CmdDumpConfig {
				t.Helper()
				outFile := filepath.Join(t.TempDir(), "empty.json")
				return &CmdDumpConfig{
					InputDatabase: testMMDB,
					OutputFile:    outFile,
					JSONPath:      `{[?(@.registered_country.iso_code=="ZZ")]}`,
				}
			},
			wantErr: false,
			verify: func(t *testing.T, cfg *CmdDumpConfig) {
				t.Helper()
				data, err := os.ReadFile(cfg.OutputFile)
				require.NoError(t, err)

				var result map[string]interface{}
				require.NoError(t, json.Unmarshal(data, &result))

				dataset, ok := result["dataset"].([]interface{})
				require.True(t, ok)
				assert.Empty(t, dataset)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.cfg(t)
			err := DumpMMMDB(cfg)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, cfg)
			}
		})
	}
}

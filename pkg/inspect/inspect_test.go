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

package inspect

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testMMDB = "../../test/inspect.mmdb"

func TestDetermineLookupNetwork(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "IPv4 address",
			input:    "192.168.1.1",
			expected: "192.168.1.1/32",
		},
		{
			name:     "IPv6 address",
			input:    "2001:db8::",
			expected: "2001:db8::/128",
		},
		{
			name:     "IPv4 CIDR",
			input:    "192.168.1.0/24",
			expected: "192.168.1.0/24",
		},
		{
			name:     "IPv6 CIDR",
			input:    "2001:db8::/32",
			expected: "2001:db8::/32",
		},
		{
			name:    "invalid input",
			input:   "invalid_input",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:     "IPv4 loopback",
			input:    "127.0.0.1",
			expected: "127.0.0.1/32",
		},
		{
			name:     "IPv6 loopback",
			input:    "::1",
			expected: "::1/128",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := determineLookupNetwork(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMmdbReader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid MMDB file",
			input:   testMMDB,
			wantErr: false,
		},
		{
			name:    "non-existent file",
			input:   "/nonexistent/file.mmdb",
			wantErr: true,
		},
		{
			name:    "empty path",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader, err := mmdbReader(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, reader)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, reader)
		})
	}
}

func TestMmdbLookup(t *testing.T) {
	reader, err := mmdbReader(testMMDB)
	require.NoError(t, err)

	tests := []struct {
		name     string
		query    string
		wantData bool
	}{
		{
			name:     "known IP returns data",
			query:    "1.1.1.1",
			wantData: true,
		},
		{
			name:     "IP in known network",
			query:    "1.0.0.1",
			wantData: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := net.ParseIP(tt.query)
			require.NotNil(t, query)

			record, err := mmdbLookup(reader, query)
			require.NoError(t, err)

			if tt.wantData {
				assert.NotNil(t, record)
				recordJSON, err := json.Marshal(record)
				require.NoError(t, err)
				assert.Contains(t, string(recordJSON), "registered_country")
			}
		})
	}
}

func TestMmdbNetworksWithin(t *testing.T) {
	reader, err := mmdbReader(testMMDB)
	require.NoError(t, err)

	tests := []struct {
		name         string
		cidr         string
		expectResult bool
	}{
		{
			name:         "CIDR containing known network",
			cidr:         "1.0.0.0/8",
			expectResult: true,
		},
		{
			name:         "exact match CIDR",
			cidr:         "1.1.1.1/32",
			expectResult: true,
		},
		{
			name:         "CIDR with no matching networks",
			cidr:         "192.168.0.0/16",
			expectResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, network, err := net.ParseCIDR(tt.cidr)
			require.NoError(t, err)

			networks := mmdbNetworksWithin(reader, network)
			assert.NotNil(t, networks)

			hasResults := networks.Next()
			assert.Equal(t, tt.expectResult, hasResults)
		})
	}
}

func TestInspectInMMDB(t *testing.T) {
	tests := []struct {
		name    string
		cfg     CmdInspectConfig
		wantErr bool
		verify  func(t *testing.T, result []byte)
	}{
		{
			name: "single IPv4 lookup",
			cfg: CmdInspectConfig{
				InputFile: testMMDB,
				Inputs:    []string{"1.1.1.1"},
			},
			wantErr: false,
			verify: func(t *testing.T, result []byte) {
				t.Helper()
				var parsed []map[string]interface{}
				require.NoError(t, json.Unmarshal(result, &parsed))
				assert.Len(t, parsed, 1)
				assert.Equal(t, "1.1.1.1", parsed[0]["query"])
				records, ok := parsed[0]["records"].([]interface{})
				require.True(t, ok)
				assert.Greater(t, len(records), 0)
			},
		},
		{
			name: "CIDR range lookup",
			cfg: CmdInspectConfig{
				InputFile: testMMDB,
				Inputs:    []string{"1.0.0.0/8"},
			},
			wantErr: false,
			verify: func(t *testing.T, result []byte) {
				t.Helper()
				var parsed []map[string]interface{}
				require.NoError(t, json.Unmarshal(result, &parsed))
				assert.Len(t, parsed, 1)
				records, ok := parsed[0]["records"].([]interface{})
				require.True(t, ok)
				assert.Greater(t, len(records), 0)
			},
		},
		{
			name: "multiple inputs",
			cfg: CmdInspectConfig{
				InputFile: testMMDB,
				Inputs:    []string{"1.1.1.1", "1.0.0.0/24"},
			},
			wantErr: false,
			verify: func(t *testing.T, result []byte) {
				t.Helper()
				var parsed []map[string]interface{}
				require.NoError(t, json.Unmarshal(result, &parsed))
				assert.Len(t, parsed, 2)
			},
		},
		{
			name: "with JSONPath filter - matching",
			cfg: CmdInspectConfig{
				InputFile: testMMDB,
				Inputs:    []string{"1.1.1.1"},
				JSONPath:  `{[?(@.registered_country.iso_code=="AU")]}`,
			},
			wantErr: false,
			verify: func(t *testing.T, result []byte) {
				t.Helper()
				var parsed []map[string]interface{}
				require.NoError(t, json.Unmarshal(result, &parsed))
				assert.Len(t, parsed, 1)
				records, ok := parsed[0]["records"].([]interface{})
				require.True(t, ok)
				assert.Greater(t, len(records), 0)
			},
		},
		{
			name: "with JSONPath filter - non-matching",
			cfg: CmdInspectConfig{
				InputFile: testMMDB,
				Inputs:    []string{"1.1.1.1"},
				JSONPath:  `{[?(@.registered_country.iso_code=="ZZ")]}`,
			},
			wantErr: false,
			verify: func(t *testing.T, result []byte) {
				t.Helper()
				var parsed []map[string]interface{}
				require.NoError(t, json.Unmarshal(result, &parsed))
				assert.Len(t, parsed, 1)
				records, ok := parsed[0]["records"].([]interface{})
				require.True(t, ok)
				assert.Empty(t, records)
			},
		},
		{
			name: "invalid MMDB file",
			cfg: CmdInspectConfig{
				InputFile: "/nonexistent/file.mmdb",
				Inputs:    []string{"1.1.1.1"},
			},
			wantErr: true,
		},
		{
			name: "invalid input IP",
			cfg: CmdInspectConfig{
				InputFile: testMMDB,
				Inputs:    []string{"not_an_ip"},
			},
			wantErr: true,
		},
		{
			name: "invalid JSONPath expression",
			cfg: CmdInspectConfig{
				InputFile: testMMDB,
				Inputs:    []string{"1.1.1.1"},
				JSONPath:  "{[?(@.field==}",
			},
			wantErr: true,
		},
		{
			name: "IP not in database",
			cfg: CmdInspectConfig{
				InputFile: testMMDB,
				Inputs:    []string{"192.168.1.1"},
			},
			wantErr: false,
			verify: func(t *testing.T, result []byte) {
				t.Helper()
				var parsed []map[string]interface{}
				require.NoError(t, json.Unmarshal(result, &parsed))
				assert.Len(t, parsed, 1)
				records, ok := parsed[0]["records"].([]interface{})
				require.True(t, ok)
				assert.Empty(t, records)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := InspectInMMDB(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
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

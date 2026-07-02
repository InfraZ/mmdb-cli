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

package diff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/InfraZ/mmdb-cli/pkg/generate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeTestMMDB writes jsonContent to a temp JSON file, generates an MMDB from it,
// and returns the path to the resulting .mmdb file.
func makeTestMMDB(t *testing.T, dir, name, jsonContent string) string {
	t.Helper()
	jsonPath := filepath.Join(dir, name+".json")
	mmdbPath := filepath.Join(dir, name+".mmdb")
	require.NoError(t, os.WriteFile(jsonPath, []byte(jsonContent), 0644))
	err := generate.GenerateMMDB(&generate.CmdGenerateConfig{
		InputDataset:   jsonPath,
		OutputDatabase: mmdbPath,
	})
	require.NoError(t, err)
	return mmdbPath
}

// makeDatasetJSON builds a minimal MMDB JSON dataset containing the given network entries.
func makeDatasetJSON(networks ...string) string {
	return `{
		"version": "v1",
		"metadata": {
			"DatabaseType": "Test-DB",
			"Description": {"en": "Test"},
			"IPVersion": 6,
			"Languages": ["en"],
			"RecordSize": 24
		},
		"dataset": [` + strings.Join(networks, ",") + `]
	}`
}

const (
	net1US = `{"network": "1.0.0.0/24", "record": {"country": "US"}}`
	net2AU = `{"network": "2.0.0.0/24", "record": {"country": "AU"}}`
	net3GB = `{"network": "3.0.0.0/24", "record": {"country": "GB"}}`
	net1DE = `{"network": "1.0.0.0/24", "record": {"country": "DE"}}`
)

func TestNormalizeLookupNetwork(t *testing.T) {
	tests := []struct {
		addr    string
		want    string
		wantErr bool
	}{
		{"1.1.1.1", "1.1.1.1/32", false},
		{"1.0.0.0/8", "1.0.0.0/8", false},
		{"2001:db8::1", "2001:db8::1/128", false},
		{"2001:db8::/32", "2001:db8::/32", false},
		{"notanip", "", true},
	}
	for _, tt := range tests {
		got, err := normalizeLookupNetwork(tt.addr)
		if tt.wantErr {
			assert.Error(t, err, "addr: %s", tt.addr)
		} else {
			require.NoError(t, err, "addr: %s", tt.addr)
			assert.Equal(t, tt.want, got)
		}
	}
}

func TestDiffMMDB_IdenticalFiles(t *testing.T) {
	dir := t.TempDir()
	mmdbA := makeTestMMDB(t, dir, "a", makeDatasetJSON(net1US, net2AU))
	mmdbB := makeTestMMDB(t, dir, "b", makeDatasetJSON(net1US, net2AU))

	result, err := DiffMMDB(CmdDiffConfig{InputA: mmdbA, InputB: mmdbB})
	require.NoError(t, err)
	assert.False(t, result.HasDiff)
	assert.Empty(t, result.Added)
	assert.Empty(t, result.Removed)
	assert.Empty(t, result.Modified)
}

func TestDiffMMDB_AddedNetwork(t *testing.T) {
	dir := t.TempDir()
	mmdbA := makeTestMMDB(t, dir, "a", makeDatasetJSON(net1US))
	mmdbB := makeTestMMDB(t, dir, "b", makeDatasetJSON(net1US, net2AU))

	result, err := DiffMMDB(CmdDiffConfig{InputA: mmdbA, InputB: mmdbB})
	require.NoError(t, err)
	assert.True(t, result.HasDiff)
	assert.Len(t, result.Added, 1)
	assert.Empty(t, result.Removed)
	assert.Empty(t, result.Modified)
	assert.Equal(t, "2.0.0.0/24", result.Added[0].Network)
}

func TestDiffMMDB_RemovedNetwork(t *testing.T) {
	dir := t.TempDir()
	mmdbA := makeTestMMDB(t, dir, "a", makeDatasetJSON(net1US, net2AU))
	mmdbB := makeTestMMDB(t, dir, "b", makeDatasetJSON(net1US))

	result, err := DiffMMDB(CmdDiffConfig{InputA: mmdbA, InputB: mmdbB})
	require.NoError(t, err)
	assert.True(t, result.HasDiff)
	assert.Empty(t, result.Added)
	assert.Len(t, result.Removed, 1)
	assert.Empty(t, result.Modified)
	assert.Equal(t, "2.0.0.0/24", result.Removed[0].Network)
}

func TestDiffMMDB_ModifiedNetwork(t *testing.T) {
	dir := t.TempDir()
	mmdbA := makeTestMMDB(t, dir, "a", makeDatasetJSON(net1US, net2AU))
	mmdbB := makeTestMMDB(t, dir, "b", makeDatasetJSON(net1DE, net2AU))

	result, err := DiffMMDB(CmdDiffConfig{InputA: mmdbA, InputB: mmdbB})
	require.NoError(t, err)
	assert.True(t, result.HasDiff)
	assert.Empty(t, result.Added)
	assert.Empty(t, result.Removed)
	assert.Len(t, result.Modified, 1)
	assert.Equal(t, "1.0.0.0/24", result.Modified[0].Network)
}

func TestDiffMMDB_ScopedFilter_WithChanges(t *testing.T) {
	dir := t.TempDir()
	mmdbA := makeTestMMDB(t, dir, "a", makeDatasetJSON(net1US, net2AU, net3GB))
	mmdbB := makeTestMMDB(t, dir, "b", makeDatasetJSON(net1DE, net2AU, net3GB))

	// Scope to 1.0.0.0/8 — only net1 changed.
	result, err := DiffMMDB(CmdDiffConfig{InputA: mmdbA, InputB: mmdbB, Filters: []string{"1.0.0.0/8"}})
	require.NoError(t, err)
	assert.True(t, result.HasDiff)
	assert.Len(t, result.Modified, 1)
	assert.Equal(t, "1.0.0.0/24", result.Modified[0].Network)
}

func TestDiffMMDB_ScopedFilter_NoChanges(t *testing.T) {
	dir := t.TempDir()
	mmdbA := makeTestMMDB(t, dir, "a", makeDatasetJSON(net1US, net2AU, net3GB))
	mmdbB := makeTestMMDB(t, dir, "b", makeDatasetJSON(net1DE, net2AU, net3GB))

	// Scope to 2.0.0.0/8 — net2 is unchanged.
	result, err := DiffMMDB(CmdDiffConfig{InputA: mmdbA, InputB: mmdbB, Filters: []string{"2.0.0.0/8"}})
	require.NoError(t, err)
	assert.False(t, result.HasDiff)
}

func TestDiffMMDB_ScopedFilter_IPAddress(t *testing.T) {
	dir := t.TempDir()
	mmdbA := makeTestMMDB(t, dir, "a", makeDatasetJSON(net1US, net2AU))
	mmdbB := makeTestMMDB(t, dir, "b", makeDatasetJSON(net1DE, net2AU))

	result, err := DiffMMDB(CmdDiffConfig{InputA: mmdbA, InputB: mmdbB, Filters: []string{"1.0.0.1"}})
	require.NoError(t, err)
	assert.True(t, result.HasDiff)
	assert.Len(t, result.Modified, 1)
}

func TestDiffMMDB_InvalidInputFile(t *testing.T) {
	_, err := DiffMMDB(CmdDiffConfig{InputA: "/nonexistent.mmdb", InputB: "/nonexistent.mmdb"})
	assert.Error(t, err)
}

func TestDiffMMDB_WrongExtension(t *testing.T) {
	_, err := DiffMMDB(CmdDiffConfig{InputA: "file.txt", InputB: "file.json"})
	assert.Error(t, err)
}

func TestDiffMMDB_InvalidFilter(t *testing.T) {
	dir := t.TempDir()
	mmdbA := makeTestMMDB(t, dir, "a", makeDatasetJSON(net1US))
	mmdbB := makeTestMMDB(t, dir, "b", makeDatasetJSON(net1US))

	_, err := DiffMMDB(CmdDiffConfig{InputA: mmdbA, InputB: mmdbB, Filters: []string{"notanip"}})
	assert.Error(t, err)
}

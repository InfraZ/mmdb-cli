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

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureAndExecute(t *testing.T, args ...string) (string, error) {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	execErr := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	var stdoutBuf bytes.Buffer
	stdoutBuf.ReadFrom(r)

	combined := stdoutBuf.String() + buf.String()
	return combined, execErr
}

func TestRootCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{})

	err := rootCmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "mmdb-cli")
}

func TestVersionCommand(t *testing.T) {
	output, err := captureAndExecute(t, "version")
	assert.NoError(t, err)
	assert.Contains(t, output, "Version:")
	assert.Contains(t, output, "Licence:")
	assert.Contains(t, output, "Documentation:")
	assert.Contains(t, output, "Maintainers:")
}

func TestVerifyCommandValid(t *testing.T) {
	output, err := captureAndExecute(t, "verify", "-i", "../test/verify-valid.mmdb")
	assert.NoError(t, err)
	assert.Contains(t, output, "valid")
}

func TestMetadataCommandJson(t *testing.T) {
	output, err := captureAndExecute(t, "metadata", "-i", "../test/metadata.mmdb", "-f", "json")
	assert.NoError(t, err)
	assert.Contains(t, output, "database_type")
	assert.Contains(t, output, "Metadata Test")
}

func TestMetadataCommandYaml(t *testing.T) {
	output, err := captureAndExecute(t, "metadata", "-i", "../test/metadata.mmdb", "-f", "yaml")
	assert.NoError(t, err)
	assert.Contains(t, output, "database_type")
}

func TestMetadataCommandJsonPretty(t *testing.T) {
	output, err := captureAndExecute(t, "metadata", "-i", "../test/metadata.mmdb", "-f", "json-pretty")
	assert.NoError(t, err)
	assert.Contains(t, output, "database_type")
	assert.Contains(t, output, "    ")
}

func TestInspectCommandIPv4(t *testing.T) {
	output, err := captureAndExecute(t, "inspect", "-i", "../test/inspect.mmdb", "-f", "json", "1.1.1.1")
	assert.NoError(t, err)
	assert.Contains(t, output, "1.1.1.1")
	assert.Contains(t, output, "registered_country")
	assert.Contains(t, output, "AU")
}

func TestInspectCommandMultipleIPs(t *testing.T) {
	output, err := captureAndExecute(t, "inspect", "-i", "../test/inspect.mmdb", "-f", "json", "1.1.1.1", "1.0.0.1")
	assert.NoError(t, err)
	assert.Contains(t, output, "1.1.1.1")
	assert.Contains(t, output, "1.0.0.1")
}

func TestInspectCommandCIDR(t *testing.T) {
	output, err := captureAndExecute(t, "inspect", "-i", "../test/inspect.mmdb", "-f", "json", "1.0.0.0/8")
	assert.NoError(t, err)
	assert.Contains(t, output, "1.0.0.0/8")
}

func TestInspectCommandMissingArgs(t *testing.T) {
	_, err := captureAndExecute(t, "inspect", "-i", "../test/inspect.mmdb")
	assert.Error(t, err)
}

func TestDumpCommand(t *testing.T) {
	dir := t.TempDir()
	outFile := filepath.Join(dir, "dump.json")

	output, err := captureAndExecute(t, "dump", "-i", "../test/inspect.mmdb", "-o", outFile)
	assert.NoError(t, err)
	assert.Contains(t, output, "MMDB Dumped successfully")

	_, statErr := os.Stat(outFile)
	assert.NoError(t, statErr)

	data, readErr := os.ReadFile(outFile)
	require.NoError(t, readErr)
	assert.Contains(t, string(data), "dataset")
	assert.Contains(t, string(data), "v1")
}

func TestDumpCommandWithJSONPath(t *testing.T) {
	dir := t.TempDir()
	outFile := filepath.Join(dir, "filtered.json")

	output, err := captureAndExecute(t, "dump", "-i", "../test/inspect.mmdb", "-o", outFile, "-j", `{[?(@.registered_country.iso_code=="AU")]}`)
	assert.NoError(t, err)
	assert.Contains(t, output, "matched")
}

func TestGenerateCommand(t *testing.T) {
	dir := t.TempDir()
	inputJSON := `{
		"version": "v1",
		"metadata": {
			"DatabaseType": "Cmd-Test",
			"Description": {"en": "Command test"}
		},
		"dataset": [
			{
				"network": "1.1.1.0/24",
				"record": {"country": "AU"}
			}
		]
	}`
	inputPath := filepath.Join(dir, "input.json")
	require.NoError(t, os.WriteFile(inputPath, []byte(inputJSON), 0644))
	outputPath := filepath.Join(dir, "output.mmdb")

	output, err := captureAndExecute(t, "generate", "-i", inputPath, "-o", outputPath)
	assert.NoError(t, err)
	assert.Contains(t, output, "Generated successfully")

	_, statErr := os.Stat(outputPath)
	assert.NoError(t, statErr)
}

func TestUpdateCommand(t *testing.T) {
	dir := t.TempDir()
	datasetJSON := `{
		"dataset": [
			{
				"network": "1.1.1.1/32",
				"method": "deep_merge",
				"data": {"extra": "field"}
			}
		]
	}`
	datasetPath := filepath.Join(dir, "update.json")
	require.NoError(t, os.WriteFile(datasetPath, []byte(datasetJSON), 0644))
	outputPath := filepath.Join(dir, "updated.mmdb")

	output, err := captureAndExecute(t, "update", "-i", "../test/inspect.mmdb", "-d", datasetPath, "-o", outputPath)
	assert.NoError(t, err)
	assert.Contains(t, output, "updated successfully")

	_, statErr := os.Stat(outputPath)
	assert.NoError(t, statErr)
}

func TestSubcommandRegistration(t *testing.T) {
	subcommands := []string{"version", "metadata", "inspect", "update", "dump", "generate", "verify"}
	registeredCmds := rootCmd.Commands()

	registeredNames := make(map[string]bool)
	for _, cmd := range registeredCmds {
		registeredNames[cmd.Name()] = true
	}

	for _, name := range subcommands {
		assert.True(t, registeredNames[name], "subcommand %q should be registered", name)
	}
}

func TestRequiredFlagsAreMarked(t *testing.T) {
	tests := []struct {
		cmdName       string
		requiredFlags []string
	}{
		{"verify", []string{"input"}},
		{"metadata", []string{"input"}},
		{"inspect", []string{"input"}},
		{"dump", []string{"input", "output"}},
		{"generate", []string{"input", "output"}},
		{"update", []string{"input", "dataset", "output"}},
	}

	for _, tt := range tests {
		t.Run(tt.cmdName, func(t *testing.T) {
			cmd, _, err := rootCmd.Find([]string{tt.cmdName})
			require.NoError(t, err)
			require.NotNil(t, cmd)

			for _, flagName := range tt.requiredFlags {
				flag := cmd.Flag(flagName)
				require.NotNil(t, flag, "flag %q should exist on command %q", flagName, tt.cmdName)

				annotations := flag.Annotations
				required, hasRequired := annotations["cobra_annotation_bash_completion_one_required_flag"]
				assert.True(t, hasRequired && len(required) > 0,
					"flag %q on command %q should be marked as required", flagName, tt.cmdName)
			}
		})
	}
}

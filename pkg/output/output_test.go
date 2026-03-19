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

package output

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return strings.TrimSpace(buf.String())
}

func TestOutput(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		options     OutputOptions
		wantErr     bool
		errContains string
	}{
		{
			name:    "JSON format",
			data:    []byte(`{"test":"data"}`),
			options: OutputOptions{Format: "json"},
		},
		{
			name:    "JSON pretty format via option",
			data:    []byte(`{"test":"data"}`),
			options: OutputOptions{Format: "json", JsonPretty: true},
		},
		{
			name:    "json-pretty format string alias",
			data:    []byte(`{"test":"data"}`),
			options: OutputOptions{Format: "json-pretty"},
		},
		{
			name:    "YAML format",
			data:    []byte(`{"test":"data"}`),
			options: OutputOptions{Format: "yaml"},
		},
		{
			name:        "unsupported format",
			data:        []byte(`{"test":"data"}`),
			options:     OutputOptions{Format: "xml"},
			wantErr:     true,
			errContains: "Unsupported output format",
		},
		{
			name:        "pretty with non-json format",
			data:        []byte(`{"test":"data"}`),
			options:     OutputOptions{Format: "yaml", JsonPretty: true},
			wantErr:     true,
			errContains: "Pretty print is only supported for JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(t, func() {
				err := Output(tt.data, tt.options)
				if tt.wantErr {
					assert.Error(t, err)
					if tt.errContains != "" {
						assert.Contains(t, err.Error(), tt.errContains)
					}
				} else {
					assert.NoError(t, err)
				}
			})

			if !tt.wantErr && output != "" {
				assert.NotEmpty(t, output)
			}
		})
	}
}

func TestOutputJsonPrettyContent(t *testing.T) {
	data := []byte(`{"name":"John","age":30}`)
	options := OutputOptions{Format: "json-pretty"}

	output := captureStdout(t, func() {
		err := Output(data, options)
		require.NoError(t, err)
	})

	assert.Contains(t, output, "\"name\": \"John\"")
	assert.Contains(t, output, "\"age\": 30")
	assert.Contains(t, output, "    ")
}

func TestOutputJsonCompactContent(t *testing.T) {
	data := []byte(`{"name":"John","age":30}`)
	options := OutputOptions{Format: "json"}

	output := captureStdout(t, func() {
		err := Output(data, options)
		require.NoError(t, err)
	})

	assert.Equal(t, `{"name":"John","age":30}`, output)
}

func TestOutputYamlContent(t *testing.T) {
	data := []byte(`{"name":"John","age":30}`)
	options := OutputOptions{Format: "yaml"}

	output := captureStdout(t, func() {
		err := Output(data, options)
		require.NoError(t, err)
	})

	assert.Contains(t, output, "age: 30")
	assert.Contains(t, output, "name: John")
}

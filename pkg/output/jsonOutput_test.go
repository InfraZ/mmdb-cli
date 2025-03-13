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
)

func TestJsonOutput(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		options  OutputOptions
		expected string
		wantErr  bool
	}{
		{
			name:     "Pretty JSON",
			data:     []byte(`{"name":"John","age":30}`),
			options:  OutputOptions{JsonPretty: true},
			expected: "{\n    \"name\": \"John\",\n    \"age\": 30\n}",
			wantErr:  false,
		},
		{
			name:     "Compact JSON",
			data:     []byte(`{"name":"John","age":30}`),
			options:  OutputOptions{JsonPretty: false},
			expected: `{"name":"John","age":30}`,
			wantErr:  false,
		},
		{
			name:     "Invalid JSON",
			data:     []byte(`{"name":"John","age":30`),
			options:  OutputOptions{JsonPretty: true},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the function
			err := JsonOutput(tt.data, tt.options)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			got := strings.TrimSpace(buf.String())

			if (err != nil) != tt.wantErr {
				t.Errorf("JsonOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("JsonOutput() = %v, want %v", got, tt.expected)
			}
		})
	}
}

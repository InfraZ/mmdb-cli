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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXmlOutput(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		wantErr     bool
		wantContain []string
	}{
		{
			name:        "Simple object",
			data:        []byte(`{"name":"John","age":30}`),
			wantContain: []string{"<root>", "<name>John</name>", "<age>30</age>", "</root>"},
		},
		{
			name:        "Nested object",
			data:        []byte(`{"country":{"iso_code":"US","name":"United States"}}`),
			wantContain: []string{"<country>", "<iso_code>US</iso_code>", "<name>United States</name>", "</country>"},
		},
		{
			name:        "Array at top level",
			data:        []byte(`[{"ip":"1.0.0.1"},{"ip":"2.0.0.1"}]`),
			wantContain: []string{"<item>", "<ip>1.0.0.1</ip>", "<ip>2.0.0.1</ip>", "</item>"},
		},
		{
			name:        "Boolean and float fields",
			data:        []byte(`{"active":true,"score":3.14}`),
			wantContain: []string{"<active>true</active>", "<score>3.14</score>"},
		},
		{
			name:        "Null field",
			data:        []byte(`{"empty":null}`),
			wantContain: []string{"<empty></empty>"},
		},
		{
			name:        "XML special characters escaped",
			data:        []byte(`{"desc":"a<b>&c"}`),
			wantContain: []string{"<desc>a&lt;b&gt;&amp;c</desc>"},
		},
		{
			name:    "Invalid JSON",
			data:    []byte(`{"broken":`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(t, func() {
				err := XmlOutput(tt.data, OutputOptions{})
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			})

			for _, s := range tt.wantContain {
				assert.Contains(t, output, s)
			}
		})
	}
}

func TestOutputXmlContent(t *testing.T) {
	data := []byte(`{"name":"John","age":30}`)
	options := OutputOptions{Format: "xml"}

	output := captureStdout(t, func() {
		err := Output(data, options)
		require.NoError(t, err)
	})

	assert.Contains(t, output, "<?xml version=")
	assert.Contains(t, output, "<root>")
	assert.Contains(t, output, "<name>John</name>")
	assert.Contains(t, output, "<age>30</age>")
	assert.Contains(t, output, "</root>")
}

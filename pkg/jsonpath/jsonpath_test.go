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

package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateExpression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		expression string
		wantErr    bool
	}{
		{
			name:       "valid simple field access",
			expression: "{.country}",
			wantErr:    false,
		},
		{
			name:       "valid nested field access",
			expression: "{.country.iso_code}",
			wantErr:    false,
		},
		{
			name:       "valid filter expression",
			expression: `{[?(@.country.iso_code=="US")]}`,
			wantErr:    false,
		},
		{
			name:       "valid wildcard",
			expression: "{.names.*}",
			wantErr:    false,
		},
		{
			name:       "invalid syntax - unclosed brace",
			expression: "{.country",
			wantErr:    true,
		},
		{
			name:       "invalid syntax - bad filter",
			expression: "{[?(@.country==}",
			wantErr:    true,
		},
		{
			name:       "empty string is valid",
			expression: "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateExpression(tt.expression)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMatchesRecord(t *testing.T) {
	t.Parallel()

	record := map[string]interface{}{
		"country": map[string]interface{}{
			"iso_code":   "US",
			"geoname_id": float64(6252001),
			"names": map[string]interface{}{
				"en": "United States",
				"de": "Vereinigte Staaten",
			},
		},
		"continent": map[string]interface{}{
			"code": "NA",
		},
	}

	tests := []struct {
		name       string
		expression string
		record     map[string]interface{}
		want       bool
		wantErr    bool
	}{
		{
			name:       "matching filter - iso_code US",
			expression: `{[?(@.country.iso_code=="US")]}`,
			record:     record,
			want:       true,
			wantErr:    false,
		},
		{
			name:       "non-matching filter - iso_code DE",
			expression: `{[?(@.country.iso_code=="DE")]}`,
			record:     record,
			want:       false,
			wantErr:    false,
		},
		{
			name:       "matching nested field",
			expression: `{[?(@.continent.code=="NA")]}`,
			record:     record,
			want:       true,
			wantErr:    false,
		},
		{
			name:       "missing field - allowed by AllowMissingKeys",
			expression: `{.nonexistent}`,
			record:     record,
			want:       false,
			wantErr:    false,
		},
		{
			name:       "empty record - no match",
			expression: `{[?(@.country.iso_code=="US")]}`,
			record:     map[string]interface{}{},
			want:       false,
			wantErr:    false,
		},
		{
			name:       "simple field access on wrapped record",
			expression: `{.country.iso_code}`,
			record:     record,
			want:       false,
			wantErr:    false,
		},
		{
			name:       "invalid expression",
			expression: `{[?(@.country==}`,
			record:     record,
			want:       false,
			wantErr:    true,
		},
		{
			name:       "nil record fields",
			expression: `{.country.iso_code}`,
			record:     map[string]interface{}{"country": nil},
			want:       false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := MatchesRecord(tt.expression, tt.record)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

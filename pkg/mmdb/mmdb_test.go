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

package mmdb

import (
	"testing"

	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func TestConvertToMMDBTypeMap(t *testing.T) {
	tests := []struct {
		name string
		data map[string]interface{}
		want mmdbtype.Map
	}{
		{
			name: "simple map",
			data: map[string]interface{}{
				"string": "value",
				"bool":   true,
				"float":  1.23,
				"int":    42,
			},
			want: mmdbtype.Map{
				mmdbtype.String("string"): mmdbtype.String("value"),
				mmdbtype.String("bool"):   mmdbtype.Bool(true),
				mmdbtype.String("float"):  mmdbtype.Float64(1.23),
				mmdbtype.String("int"):    mmdbtype.Int32(42),
			},
		},
		{
			name: "nested map",
			data: map[string]interface{}{
				"nested": map[string]interface{}{
					"string": "nestedValue",
				},
			},
			want: mmdbtype.Map{
				mmdbtype.String("nested"): mmdbtype.Map{
					mmdbtype.String("string"): mmdbtype.String("nestedValue"),
				},
			},
		},
		{
			name: "unsupported type",
			data: map[string]interface{}{
				"unsupported": []string{"value1", "value2"},
			},
			want: mmdbtype.Map{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToMMDBTypeMap(tt.data)
			if !compareMMDBTypeMaps(got, tt.want) {
				t.Errorf("ConvertToMMDBTypeMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareMMDBTypeMaps(a, b mmdbtype.Map) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || !compareMMDBTypes(v, bv) {
			return false
		}
	}
	return true
}

func compareMMDBTypes(a, b mmdbtype.DataType) bool {
	switch a := a.(type) {
	case mmdbtype.String:
		if b, ok := b.(mmdbtype.String); ok {
			return a == b
		}
	case mmdbtype.Bool:
		if b, ok := b.(mmdbtype.Bool); ok {
			return a == b
		}
	case mmdbtype.Float64:
		if b, ok := b.(mmdbtype.Float64); ok {
			return a == b
		}
	case mmdbtype.Int32:
		if b, ok := b.(mmdbtype.Int32); ok {
			return a == b
		}
	case mmdbtype.Map:
		if b, ok := b.(mmdbtype.Map); ok {
			return compareMMDBTypeMaps(a, b)
		}
	}
	return false
}

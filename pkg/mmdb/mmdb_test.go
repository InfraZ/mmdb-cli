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
	"github.com/stretchr/testify/assert"
)

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
	case mmdbtype.Uint16:
		if b, ok := b.(mmdbtype.Uint16); ok {
			return a == b
		}
	case mmdbtype.Uint32:
		if b, ok := b.(mmdbtype.Uint32); ok {
			return a == b
		}
	case mmdbtype.Map:
		if b, ok := b.(mmdbtype.Map); ok {
			return compareMMDBTypeMaps(a, b)
		}
	}
	return false
}

func TestConvertToMMDBTypeMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data map[string]interface{}
		want mmdbtype.Map
	}{
		{
			name: "simple map with all default types",
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
			name: "float64 values",
			data: map[string]interface{}{
				"autonomous_system_number": float64(12345),
				"negative_float64_int":     float64(-12345),
				"real_float":               1.23,
			},
			want: mmdbtype.Map{
				mmdbtype.String("autonomous_system_number"): mmdbtype.Float64(12345),
				mmdbtype.String("negative_float64_int"):     mmdbtype.Float64(-12345),
				mmdbtype.String("real_float"):               mmdbtype.Float64(1.23),
			},
		},
		{
			name: "negative integers",
			data: map[string]interface{}{
				"negative_int": -42,
			},
			want: mmdbtype.Map{
				mmdbtype.String("negative_int"): mmdbtype.Int32(-42),
			},
		},
		{
			name: "unsupported type is skipped",
			data: map[string]interface{}{
				"unsupported": []string{"value1", "value2"},
			},
			want: mmdbtype.Map{},
		},
		{
			name: "empty map",
			data: map[string]interface{}{},
			want: mmdbtype.Map{},
		},
		{
			name: "deeply nested map",
			data: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"value": "deep",
					},
				},
			},
			want: mmdbtype.Map{
				mmdbtype.String("level1"): mmdbtype.Map{
					mmdbtype.String("level2"): mmdbtype.Map{
						mmdbtype.String("value"): mmdbtype.String("deep"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ConvertToMMDBTypeMap(tt.data, true, nil)
			assert.True(t, compareMMDBTypeMaps(got, tt.want), "ConvertToMMDBTypeMap() = %v, want %v", got, tt.want)
		})
	}
}

func TestConvertWithSchema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		data   map[string]interface{}
		schema map[string]interface{}
		want   mmdbtype.Map
	}{
		{
			name: "string schema type",
			data: map[string]interface{}{
				"name": "test",
			},
			schema: map[string]interface{}{
				"name": "string",
			},
			want: mmdbtype.Map{
				mmdbtype.String("name"): mmdbtype.String("test"),
			},
		},
		{
			name: "bool schema type",
			data: map[string]interface{}{
				"active": true,
			},
			schema: map[string]interface{}{
				"active": "bool",
			},
			want: mmdbtype.Map{
				mmdbtype.String("active"): mmdbtype.Bool(true),
			},
		},
		{
			name: "boolean schema type alias",
			data: map[string]interface{}{
				"active": false,
			},
			schema: map[string]interface{}{
				"active": "boolean",
			},
			want: mmdbtype.Map{
				mmdbtype.String("active"): mmdbtype.Bool(false),
			},
		},
		{
			name: "float64 schema type",
			data: map[string]interface{}{
				"score": float64(3.14),
			},
			schema: map[string]interface{}{
				"score": "float64",
			},
			want: mmdbtype.Map{
				mmdbtype.String("score"): mmdbtype.Float64(3.14),
			},
		},
		{
			name: "float schema type alias",
			data: map[string]interface{}{
				"score": float64(2.71),
			},
			schema: map[string]interface{}{
				"score": "float",
			},
			want: mmdbtype.Map{
				mmdbtype.String("score"): mmdbtype.Float64(2.71),
			},
		},
		{
			name: "uint16 from float64",
			data: map[string]interface{}{
				"port": float64(8080),
			},
			schema: map[string]interface{}{
				"port": "uint16",
			},
			want: mmdbtype.Map{
				mmdbtype.String("port"): mmdbtype.Uint16(8080),
			},
		},
		{
			name: "uint16 from int",
			data: map[string]interface{}{
				"port": 443,
			},
			schema: map[string]interface{}{
				"port": "uint16",
			},
			want: mmdbtype.Map{
				mmdbtype.String("port"): mmdbtype.Uint16(443),
			},
		},
		{
			name: "int32 from float64",
			data: map[string]interface{}{
				"geoname_id": float64(12345),
			},
			schema: map[string]interface{}{
				"geoname_id": "int32",
			},
			want: mmdbtype.Map{
				mmdbtype.String("geoname_id"): mmdbtype.Int32(12345),
			},
		},
		{
			name: "int from int",
			data: map[string]interface{}{
				"count": 100,
			},
			schema: map[string]interface{}{
				"count": "int",
			},
			want: mmdbtype.Map{
				mmdbtype.String("count"): mmdbtype.Int32(100),
			},
		},
		{
			name: "uint32 from float64",
			data: map[string]interface{}{
				"asn": float64(13335),
			},
			schema: map[string]interface{}{
				"asn": "uint32",
			},
			want: mmdbtype.Map{
				mmdbtype.String("asn"): mmdbtype.Uint32(13335),
			},
		},
		{
			name: "uint from int",
			data: map[string]interface{}{
				"asn": 13335,
			},
			schema: map[string]interface{}{
				"asn": "uint",
			},
			want: mmdbtype.Map{
				mmdbtype.String("asn"): mmdbtype.Uint32(13335),
			},
		},
		{
			name: "nested schema",
			data: map[string]interface{}{
				"country": map[string]interface{}{
					"iso_code":   "US",
					"geoname_id": float64(6252001),
				},
			},
			schema: map[string]interface{}{
				"country": map[string]interface{}{
					"iso_code":   "string",
					"geoname_id": "uint32",
				},
			},
			want: mmdbtype.Map{
				mmdbtype.String("country"): mmdbtype.Map{
					mmdbtype.String("iso_code"):   mmdbtype.String("US"),
					mmdbtype.String("geoname_id"): mmdbtype.Uint32(6252001),
				},
			},
		},
		{
			name: "field without schema falls back to default",
			data: map[string]interface{}{
				"known":   "value",
				"unknown": "fallback",
			},
			schema: map[string]interface{}{
				"known": "string",
			},
			want: mmdbtype.Map{
				mmdbtype.String("known"):   mmdbtype.String("value"),
				mmdbtype.String("unknown"): mmdbtype.String("fallback"),
			},
		},
		{
			name: "type mismatch falls back to default value",
			data: map[string]interface{}{
				"name": 12345,
			},
			schema: map[string]interface{}{
				"name": "string",
			},
			want: mmdbtype.Map{
				mmdbtype.String("name"): mmdbtype.String(""),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ConvertToMMDBTypeMap(tt.data, false, tt.schema)
			assert.True(t, compareMMDBTypeMaps(got, tt.want), "ConvertToMMDBTypeMap() with schema = %v, want %v", got, tt.want)
		})
	}
}

func TestConvertValueDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value interface{}
		want  mmdbtype.DataType
	}{
		{
			name:  "string value",
			value: "hello",
			want:  mmdbtype.String("hello"),
		},
		{
			name:  "bool value",
			value: true,
			want:  mmdbtype.Bool(true),
		},
		{
			name:  "float64 value",
			value: 3.14,
			want:  mmdbtype.Float64(3.14),
		},
		{
			name:  "int value",
			value: 42,
			want:  mmdbtype.Int32(42),
		},
		{
			name:  "unsupported type returns empty string",
			value: []string{"a", "b"},
			want:  mmdbtype.String(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := convertValueDefault(tt.value, "test_key")
			assert.True(t, compareMMDBTypes(got, tt.want), "convertValueDefault() = %v, want %v", got, tt.want)
		})
	}
}

func TestConvertValueWithType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		value        interface{}
		expectedType string
		want         mmdbtype.DataType
	}{
		{
			name:         "string type with string value",
			value:        "hello",
			expectedType: "string",
			want:         mmdbtype.String("hello"),
		},
		{
			name:         "string type with wrong value type",
			value:        42,
			expectedType: "string",
			want:         mmdbtype.String(""),
		},
		{
			name:         "bool type",
			value:        true,
			expectedType: "bool",
			want:         mmdbtype.Bool(true),
		},
		{
			name:         "boolean type alias",
			value:        false,
			expectedType: "boolean",
			want:         mmdbtype.Bool(false),
		},
		{
			name:         "bool type with wrong value type",
			value:        "true",
			expectedType: "bool",
			want:         mmdbtype.Bool(false),
		},
		{
			name:         "float64 type",
			value:        3.14,
			expectedType: "float64",
			want:         mmdbtype.Float64(3.14),
		},
		{
			name:         "float type alias",
			value:        2.71,
			expectedType: "float",
			want:         mmdbtype.Float64(2.71),
		},
		{
			name:         "float64 type with wrong value type",
			value:        "3.14",
			expectedType: "float64",
			want:         mmdbtype.Float64(0),
		},
		{
			name:         "uint16 from float64",
			value:        float64(1024),
			expectedType: "uint16",
			want:         mmdbtype.Uint16(1024),
		},
		{
			name:         "uint16 from int",
			value:        1024,
			expectedType: "uint16",
			want:         mmdbtype.Uint16(1024),
		},
		{
			name:         "uint16 with wrong type",
			value:        "1024",
			expectedType: "uint16",
			want:         mmdbtype.Uint16(0),
		},
		{
			name:         "int32 from float64",
			value:        float64(-100),
			expectedType: "int32",
			want:         mmdbtype.Int32(-100),
		},
		{
			name:         "int type from int",
			value:        -100,
			expectedType: "int",
			want:         mmdbtype.Int32(-100),
		},
		{
			name:         "int32 with wrong type",
			value:        "not_an_int",
			expectedType: "int32",
			want:         mmdbtype.Int32(0),
		},
		{
			name:         "uint32 from float64",
			value:        float64(65535),
			expectedType: "uint32",
			want:         mmdbtype.Uint32(65535),
		},
		{
			name:         "uint type from int",
			value:        65535,
			expectedType: "uint",
			want:         mmdbtype.Uint32(65535),
		},
		{
			name:         "uint32 with wrong type",
			value:        "not_a_uint",
			expectedType: "uint32",
			want:         mmdbtype.Uint32(0),
		},
		{
			name:         "unknown schema type falls back to default",
			value:        "hello",
			expectedType: "unknown_type",
			want:         mmdbtype.String("hello"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := convertValueWithType(tt.value, tt.expectedType, "test_key")
			assert.True(t, compareMMDBTypes(got, tt.want), "convertValueWithType() = %v, want %v", got, tt.want)
		})
	}
}

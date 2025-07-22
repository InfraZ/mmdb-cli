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
	"log"

	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func ConvertToMMDBTypeMap(data map[string]interface{}, useDefaultSchema bool, schema map[string]interface{}) mmdbtype.Map {
	mmdbMap := mmdbtype.Map{}

	if !useDefaultSchema && schema != nil {
		// Use user-defined schema
		return convertWithSchema(data, schema)
	}

	// Use default schema (original switch case logic)
	for key, value := range data {
		mmdbKey := mmdbtype.String(key)
		switch mmdbValue := value.(type) {
		case string:
			mmdbMap[mmdbKey] = mmdbtype.String(mmdbValue)
		case bool:
			mmdbMap[mmdbKey] = mmdbtype.Bool(mmdbValue)
		case float64:
			mmdbMap[mmdbKey] = mmdbtype.Float64(mmdbValue)
		case int:
			mmdbMap[mmdbKey] = mmdbtype.Int32(mmdbValue)
		case map[string]interface{}:
			// Recursively convert nested maps
			mmdbMap[mmdbKey] = ConvertToMMDBTypeMap(mmdbValue, useDefaultSchema, schema)
		default:
			log.Printf("Unsupported data type for key %v", key)
		}
	}
	return mmdbMap
}

func convertWithSchema(data map[string]interface{}, schema map[string]interface{}) mmdbtype.Map {
	mmdbMap := mmdbtype.Map{}

	for key, value := range data {
		mmdbKey := mmdbtype.String(key)

		// Get the schema for this key
		var schemaForKey interface{}
		var hasSchema bool
		if schema != nil {
			schemaForKey, hasSchema = schema[key]
		}

		if hasSchema {
			// Handle schema-defined conversion
			switch schemaValue := schemaForKey.(type) {
			case string:
				// Simple type definition
				mmdbMap[mmdbKey] = convertValueWithType(value, schemaValue, key)
			case map[string]interface{}:
				// Nested object with its own schema
				if nestedData, ok := value.(map[string]interface{}); ok {
					mmdbMap[mmdbKey] = convertWithSchema(nestedData, schemaValue)
				} else {
					log.Printf("Expected nested object for key %s, got %T", key, value)
				}
			default:
				// Schema value is not string or map, fall back to default
				mmdbMap[mmdbKey] = convertValueDefault(value, key)
			}
		} else {
			// No schema for this key, use default conversion
			mmdbMap[mmdbKey] = convertValueDefault(value, key)
		}
	}

	return mmdbMap
}

func convertValueWithType(value interface{}, expectedType string, key string) mmdbtype.DataType {
	switch expectedType {
	case "string":
		if str, ok := value.(string); ok {
			return mmdbtype.String(str)
		} else {
			log.Printf("Expected string for key %s, got %T", key, value)
			return mmdbtype.String("")
		}
	case "bool", "boolean":
		if b, ok := value.(bool); ok {
			return mmdbtype.Bool(b)
		} else {
			log.Printf("Expected bool for key %s, got %T", key, value)
			return mmdbtype.Bool(false)
		}
	case "float", "float64":
		if f, ok := value.(float64); ok {
			return mmdbtype.Float64(f)
		} else {
			log.Printf("Expected float64 for key %s, got %T", key, value)
			return mmdbtype.Float64(0)
		}
	case "int", "int32":
		if i, ok := value.(int); ok {
			return mmdbtype.Int32(i)
		} else if f, ok := value.(float64); ok {
			return mmdbtype.Int32(int(f))
		} else {
			log.Printf("Expected int for key %s, got %T", key, value)
			return mmdbtype.Int32(0)
		}
	case "uint", "uint32":
		if i, ok := value.(int); ok && i >= 0 {
			return mmdbtype.Uint32(uint32(i))
		} else if f, ok := value.(float64); ok && f >= 0 {
			return mmdbtype.Uint32(uint32(f))
		} else {
			log.Printf("Expected uint for key %s, got %T", key, value)
			return mmdbtype.Uint32(0)
		}
	default:
		// Unknown type in schema, fall back to default
		return convertValueDefault(value, key)
	}
}

func convertValueDefault(value interface{}, key string) mmdbtype.DataType {
	switch v := value.(type) {
	case string:
		return mmdbtype.String(v)
	case bool:
		return mmdbtype.Bool(v)
	case float64:
		return mmdbtype.Float64(v)
	case int:
		return mmdbtype.Int32(v)
	case map[string]interface{}:
		// For nested maps without schema, use default conversion
		return ConvertToMMDBTypeMap(v, true, nil)
	default:
		log.Printf("Unsupported data type for key %s: %T", key, value)
		return mmdbtype.String("")
	}
}

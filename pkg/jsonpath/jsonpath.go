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
	"bytes"
	"fmt"

	k8sjsonpath "k8s.io/client-go/util/jsonpath"
)

// ValidateExpression parses the expression and returns an error if it is
// syntactically invalid. Call this once up-front to give users a clear error
// before starting any expensive iteration.
func ValidateExpression(expression string) error {
	j := k8sjsonpath.New("validate")
	if err := j.Parse(expression); err != nil {
		return fmt.Errorf("invalid jsonpath expression: %w", err)
	}
	return nil
}

// MatchesRecord evaluates a kubectl-style JSONPath filter expression against an
// MMDB record. The expression is applied to a single-element slice that wraps
// the record, so @ refers directly to the record's fields:
//
//	{[?(@.country.iso_code=="US")]}
//
// Returns true when the expression produces non-empty output (i.e. the filter
// matched), false when the output is empty (no match), and an error when the
// expression itself is invalid.
func MatchesRecord(expression string, record map[string]interface{}) (bool, error) {
	j := k8sjsonpath.New("filter").AllowMissingKeys(true)
	if err := j.Parse(expression); err != nil {
		return false, fmt.Errorf("invalid jsonpath expression: %w", err)
	}
	var buf bytes.Buffer
	// Wrap record in a slice so filter syntax [?(@.field==...)] can iterate.
	if err := j.Execute(&buf, []interface{}{record}); err != nil {
		return false, fmt.Errorf("jsonpath execution error: %w", err)
	}
	return buf.Len() > 0, nil
}

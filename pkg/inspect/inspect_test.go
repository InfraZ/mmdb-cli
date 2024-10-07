// Copyright 2024 The MMDB CLI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inspect

import (
	"testing"
)

func TestDetermineLookupNetwork(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"192.168.1.1", "192.168.1.1/32"},
		{"2001:db8::", "2001:db8::/128"},
		{"192.168.1.0/24", "192.168.1.0/24"},
		{"2001:db8::/32", "2001:db8::/32"},
	}

	for _, test := range tests {
		result := determineLookupNetwork(test.input)
		if result != test.expected {
			t.Errorf("determineLookupNetwork(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestDetermineLookupNetworkInvalidInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("determineLookupNetwork did not panic on invalid input")
		}
	}()

	determineLookupNetwork("invalid_input")
}

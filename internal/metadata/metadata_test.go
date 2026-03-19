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

package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadataConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
	}{
		{"Version is set", Version},
		{"License is set", License},
		{"ShortDescription is set", ShortDescription},
		{"Homepage is set", Homepage},
		{"DocumentationURL is set", DocumentationURL},
		{"Organization is set", Organization},
		{"OrganizationWebsite is set", OrganizationWebsite},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.NotEmpty(t, tt.value)
		})
	}
}

func TestMaintainers(t *testing.T) {
	t.Parallel()
	assert.NotEmpty(t, Maintainers, "Maintainers should have at least one entry")
	for _, m := range Maintainers {
		assert.NotEmpty(t, m, "Each maintainer entry should be non-empty")
	}
}

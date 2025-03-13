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

package verify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyMMDB(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		want      bool
		wantErr   bool
	}{
		{
			name:      "Valid MMDB file",
			inputFile: "../../test/verify-valid.mmdb",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "Invalid MMDB file",
			inputFile: "../../test/verify-invalid.mmdb",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "Non-existent file",
			inputFile: "testdata/nonexistent.mmdb",
			want:      false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := CmdVerifyConfig{InputFile: tt.inputFile}
			got, err := VerifyMMDB(cfg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

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

import "testing"

func TestOutput(t *testing.T) {
	testData := []byte("test data")
	testOptions := OutputOptions{Format: "json", JsonPretty: false}

	err := Output(testData, testOptions)
	if err != nil {
		t.Errorf("Output() error = %v; want nil", err)
	}
}

func TestOutputUnsupportedFormat(t *testing.T) {
	testData := []byte("test data")
	testOptions := OutputOptions{Format: "unsupported", JsonPretty: false}

	err := Output(testData, testOptions)
	if err == nil {
		t.Errorf("Output() error = nil; want error")
	}
}

func TestOutputJson(t *testing.T) {
	testData := []byte("{\"test\": \"data\"}")
	testOptions := OutputOptions{Format: "json", JsonPretty: false}

	err := Output(testData, testOptions)
	if err != nil {
		t.Errorf("Output() error = %v; want nil", err)
	}
}

func TestOutputJsonPretty(t *testing.T) {
	testData := []byte("{\"test\": \"data\"}")
	testOptions := OutputOptions{Format: "json", JsonPretty: true}

	err := Output(testData, testOptions)
	if err != nil {
		t.Errorf("Output() error = %v; want nil", err)
	}
}

func TestOutputYaml(t *testing.T) {
	testData := []byte("test: data")
	testOptions := OutputOptions{Format: "yaml", JsonPretty: false}

	err := Output(testData, testOptions)
	if err != nil {
		t.Errorf("Output() error = %v; want nil", err)
	}
}

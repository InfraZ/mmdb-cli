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

package files

import (
	"fmt"
	"os"
	"strings"
)

type FilesListValidation struct {
	FilePath          string
	ExpectedExtension string
	ShouldExist       bool
}

// CheckFileExists checks if a file exists
func CheckFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// CheckFileSizeMb checks the size of a file in MB
func CheckFileSizeMb(filePath string) (float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return float64(fileInfo.Size()) / 1024 / 1024, nil
}

// CheckFileExtension checks if a file has a specific extension
func CheckFileExtension(filePath string, extension string) bool {
	return strings.HasSuffix(filePath, extension)
}

func FilesValidation(filesList []FilesListValidation) error {
	for _, file := range filesList {
		if file.ShouldExist && !CheckFileExists(file.FilePath) {
			return fmt.Errorf("[!] File %s does not exist", file.FilePath)
		}

		if !CheckFileExtension(file.FilePath, file.ExpectedExtension) {
			return fmt.Errorf("[!] File %s must have a %s extension", file.FilePath, file.ExpectedExtension)
		}
	}

	return nil
}

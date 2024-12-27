package files

import (
	"os"
	"testing"
)

func TestCheckFileExists(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Test case: File exists
	if !CheckFileExists(tmpFile.Name()) {
		t.Errorf("Expected file to exist, but it does not")
	}

	// Test case: File does not exist
	nonExistentFile := "/path/to/nonexistent/file"
	if CheckFileExists(nonExistentFile) {
		t.Errorf("Expected file to not exist, but it does")
	}
}

func TestCheckFileSizeMb(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some data to the file
	data := []byte("Hello, World!")
	if _, err := tmpFile.Write(data); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	// Test case: Check file size
	expectedSizeMb := float64(len(data)) / 1024 / 1024
	sizeMb, err := CheckFileSizeMb(tmpFile.Name())
	if err != nil {
		t.Errorf("Failed to check file size: %v", err)
	}
	if sizeMb != expectedSizeMb {
		t.Errorf("Expected file size to be %f MB, but got %f MB", expectedSizeMb, sizeMb)
	}

	// Test case: File does not exist
	nonExistentFile := "/path/to/nonexistent/file"
	_, err = CheckFileSizeMb(nonExistentFile)
	if err == nil {
		t.Errorf("Expected error for non-existent file, but got none")
	}
}

func TestCheckFileExtension(t *testing.T) {
	// Test case: File has the correct extension
	filePath := "example.txt"
	extension := ".txt"
	if !CheckFileExtension(filePath, extension) {
		t.Errorf("Expected file %s to have extension %s, but it does not", filePath, extension)
	}

	// Test case: File does not have the correct extension
	filePath = "example.txt"
	extension = ".pdf"
	if CheckFileExtension(filePath, extension) {
		t.Errorf("Expected file %s to not have extension %s, but it does", filePath, extension)
	}

	// Test case: File has no extension
	filePath = "example"
	extension = ".txt"
	if CheckFileExtension(filePath, extension) {
		t.Errorf("Expected file %s to not have extension %s, but it does", filePath, extension)
	}

	// Test case: Empty file path
	filePath = ""
	extension = ".txt"
	if CheckFileExtension(filePath, extension) {
		t.Errorf("Expected empty file path to not have extension %s, but it does", extension)
	}

	// Test case: Empty extension
	filePath = "example.txt"
	extension = ""
	if !CheckFileExtension(filePath, extension) {
		t.Errorf("Expected file %s to have empty extension, but it does not", filePath)
	}
}

func TestFilesValidation(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name      string
		filesList []FilesListValidation
		wantErr   bool
	}{
		{
			name: "All files exist with correct extensions",
			filesList: []FilesListValidation{
				{FilePath: tmpFile.Name(), ExpectedExtension: ".txt", ShouldExist: true},
			},
			wantErr: false,
		},
		{
			name: "File does not exist",
			filesList: []FilesListValidation{
				{FilePath: "/path/to/nonexistent/file.txt", ExpectedExtension: ".txt", ShouldExist: true},
			},
			wantErr: true,
		},
		{
			name: "File exists but has incorrect extension",
			filesList: []FilesListValidation{
				{FilePath: tmpFile.Name(), ExpectedExtension: ".pdf", ShouldExist: true},
			},
			wantErr: true,
		},
		{
			name: "File does not exist and has incorrect extension",
			filesList: []FilesListValidation{
				{FilePath: "/path/to/nonexistent/file.pdf", ExpectedExtension: ".pdf", ShouldExist: true},
			},
			wantErr: true,
		},
		{
			name: "File exists and should not exist",
			filesList: []FilesListValidation{
				{FilePath: tmpFile.Name(), ExpectedExtension: ".txt", ShouldExist: false},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FilesValidation(tt.filesList)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesValidation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package filescaner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_dir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{
		"file1.txt",
		"file2.ghtml",
		"file3.txt",
		"subdir/file4.ghtml",
		"subdir/file5.txt",
	}
	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		file, err := create(filePath)
		_, err = file.Write([]byte("test data"))
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
	}

	// Call the function being tested
	data, errors := readFiles(tempDir, true)

	// Assert the expected number of files and errors
	expectedFiles := 2
	expectedErrors := 0
	if len(data) != expectedFiles {
		t.Errorf("Unexpected number of files.\nExpected: %d\nActual: %d", expectedFiles, len(data))
	}
	if len(errors) != expectedErrors {
		t.Errorf("Unexpected number of errors.\nExpected: %d\nActual: %d", expectedErrors, len(errors))
	}
}

func TestReadFiles_NoMatchingFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_dir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Call the function being tested
	data, errors := readFiles(tempDir, true)

	// Assert there are no files and no errors
	expectedFiles := 0
	expectedErrors := 0
	if len(data) != expectedFiles {
		t.Errorf("Unexpected number of files.\nExpected: %d\nActual: %d", expectedFiles, len(data))
	}
	if len(errors) != expectedErrors {
		t.Errorf("Unexpected number of errors.\nExpected: %d\nActual: %d", expectedErrors, len(errors))
	}
}

func TestReadFiles_ReadAllFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_dir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{
		"file1.txt",
		"file2.ghtml",
		"file3.txt",
		"subdir/file4.ghtml",
		"subdir/file5.txt",
	}
	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		file, err := create(filePath)
		_, err = file.Write([]byte("test data"))
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
	}

	// Call the function being tested
	data, errors := readFiles(tempDir, false)

	// Assert the expected number of files and no errors
	expectedFiles := 5
	expectedErrors := 0
	if len(data) != expectedFiles {
		t.Errorf("Unexpected number of files.\nExpected: %d\nActual: %d", expectedFiles, len(data))
	}
	if len(errors) != expectedErrors {
		t.Errorf("Unexpected number of errors.\nExpected: %d\nActual: %d", expectedErrors, len(errors))
	}
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

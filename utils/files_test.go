package utils

import (
	"os"
	"testing"
)

func TestWrite2FileLineByLine(t *testing.T) {
	// Create a temporary test file
	file, err := os.CreateTemp("", "test_file.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Define test data
	lines := []string{"Line 1", "Line 2", "Line 3"}

	// Call the function being tested
	err = write2FileLineByLine(file, lines)
	if err != nil {
		t.Errorf("Error writing lines to file: %v", err)
	}

	// Close the file to ensure all writes are flushed
	err = file.Close()
	if err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}

	// Read the contents of the written file
	content, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	// Assert the content matches the expected lines
	expected := "Line 1\nLine 2\nLine 3\n"
	if string(content) != expected {
		t.Errorf("Unexpected file content.\nExpected: %s\nActual: %s", expected, string(content))
	}
}

func TestRemoveHTMLComment(t *testing.T) {
	input := `Hello <!-- Comment --> World!`
	expected := "Hello  World!"

	output := removeHtmlComment(input)

	if output != expected {
		t.Errorf("Unexpected result.\nExpected: %s\nActual: %s", expected, output)
	}
}

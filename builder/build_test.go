package builder

import (
	"testing"
)

func TestRemoveSpaceAndTab(t *testing.T) {
	input := "   Hello    World!	"
	expected := "HelloWorld!"

	output := removeSpaceAndTab(input)

	if output != expected {
		t.Errorf("Unexpected result.\nExpected: %s\nActual: %s", expected, output)
	}
}

func TestGetFileNameOnly(t *testing.T) {
	input := "/path/to/file.txt"
	expected := "file"

	output := getFileNameOnly(input)

	if output != expected {
		t.Errorf("Unexpected result.\nExpected: %s\nActual: %s", expected, output)
	}
}

func TestGetFileNameOnly_NoExtension(t *testing.T) {
	input := "/path/to/file"
	expected := "file"

	output := getFileNameOnly(input)

	if output != expected {
		t.Errorf("Unexpected result.\nExpected: %s\nActual: %s", expected, output)
	}
}

func TestGetFileNameOnly_EmptyPath(t *testing.T) {
	input := ""
	expected := ""

	output := getFileNameOnly(input)

	if output != expected {
		t.Errorf("Unexpected result.\nExpected: %s\nActual: %s", expected, output)
	}
}

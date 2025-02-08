package fsPatrol

import (
	"reflect"
	"testing"
	
)

func Test_getState(t *testing.T) {
	// Mock ScanFS function
	mockScanFS := func(path string) (map[string]string, []error) {
		// Return mock snapshot and no errors
		return map[string]string{
			"file1": "content1",
			"file2": "content2",
		}, nil
	}

	// Mock ScanGhtmlFilesOnly function
	mockScanGhtmlFilesOnly := func(path string) (map[string]string, []error) {
		// Return mock snapshot and no errors
		return map[string]string{
			"file3": "content3",
			"file4": "content4",
		}, nil
	}

	expectedSnap :=  FsSnap(map[string]string{
			"file1": "content1",
			"file2": "content2",
			"file3": "content3",
			"file4": "content4",
	})

	snap, errs := getState(mockScanFS, mockScanGhtmlFilesOnly)

	// Check if the returned snapshot matches the expected snapshot
	if !reflect.DeepEqual(snap, expectedSnap) {
		t.Errorf("GetState returned unexpected snapshot:\nExpected: %+v\nGot: %+v", expectedSnap, snap)
	}

	// Check if there are no errors
	if len(errs) != 0 {
		t.Errorf("GetState returned unexpected errors: %+v", errs)
	}
}


func TestCombiningSnap(t *testing.T) {
	mockSnaps := []FsSnap{
		 map[string]string{"file1": "content1", "file2": "content2"},
		 map[string]string{"file3": "content3", "file4": "content4"},
	}

	expectedSnap := FsSnap(map[string]string{"file1": "content1", "file2": "content2", "file3": "content3", "file4": "content4"})

	snap := combiningSnap(mockSnaps)



	if !reflect.DeepEqual(snap, expectedSnap) {
		t.Errorf("Expected snap FS %v, but got %v", expectedSnap, snap)
	}
}

func TestIsDiffState(t *testing.T) {
	filesSourceState := map[string]string{
		"file1": "content1",
		"file2": "content2",
		"file3": "content3",
	}

	filesCurrentState := map[string]string{
		"file1": "content1",
		"file2": "modified",
		"file4": "content4",
	}

	expectedDiffState := true

	diffState := IsDiffState(filesSourceState, filesCurrentState)

	if diffState != expectedDiffState {
		t.Errorf("Expected diff state %t, but got %t", expectedDiffState, diffState)
	}
}

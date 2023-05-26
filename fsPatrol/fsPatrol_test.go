package fsPatrol

import (
	"reflect"
	"testing"
	"time"
)

func TestCombiningSnap(t *testing.T) {
	mockSnaps := []FsSnap{
		{SnapTime: 1234567890, SnapFS: map[string]string{"file1": "content1", "file2": "content2"}},
		{SnapTime: 1234567890, SnapFS: map[string]string{"file3": "content3", "file4": "content4"}},
	}

	expectedSnap := FsSnap{
		SnapTime: time.Now().Unix(),
		SnapFS:   map[string]string{"file1": "content1", "file2": "content2", "file3": "content3", "file4": "content4"},
	}

	snap := combiningSnap(mockSnaps)

	if snap.SnapTime != expectedSnap.SnapTime {
		t.Errorf("Expected snap time %d, but got %d", expectedSnap.SnapTime, snap.SnapTime)
	}

	if !reflect.DeepEqual(snap.SnapFS, expectedSnap.SnapFS) {
		t.Errorf("Expected snap FS %v, but got %v", expectedSnap.SnapFS, snap.SnapFS)
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

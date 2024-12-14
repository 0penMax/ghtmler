package fsPatrol

import (
	"goHtmlBuilder/filescaner"
)

const (
	ghtmlFolder      = "ghtml"
	componentsFolder = "components"
	staticFolder     = "static"
)

// FsSnap represents a snapshot of file paths and their hashes.
type FsSnap map[string]string //[path]hash

// GetGhtmlFiles filters the snapshot for only GHTML files.
func (s FsSnap) GetGhtmlFiles() []string {
	var result []string
	for path := range s {
		if filescaner.IsGhtmlFile(path) {
			result = append(result, path)
		}
	}
	return result
}

// GetState captures the current state of the filesystem for predefined folders.
func GetState() (FsSnap, []error) {
	return collectState(filescaner.ScanFS, filescaner.ScanGhtmlFilesOnly)
}

// collectState aggregates the snapshots of specified folders and files.
func collectState(
	scanFS, scanGhtmlFilesOnly func(path string) (map[string]string, []error),
) (FsSnap, []error) {
	snapshots := []FsSnap{}
	folders := []string{componentsFolder, staticFolder}

	// Scan folders
	for _, folder := range folders {
		snap, errs := scanFS(folder)
		if len(errs) > 0 {
			return nil, errs
		}
		snapshots = append(snapshots, snap)
	}

	// Scan GHTML folder
	ghtmlSnap, errs := scanGhtmlFilesOnly(ghtmlFolder)
	if len(errs) > 0 {
		return nil, errs
	}
	snapshots = append(snapshots, ghtmlSnap)

	// Combine all snapshots
	return mergeSnapshots(snapshots), nil
}

// mergeSnapshots merges multiple FsSnap instances into one.
func mergeSnapshots(snaps []FsSnap) FsSnap {
	merged := FsSnap{}
	for _, snap := range snaps {
		for path, hash := range snap {
			merged[path] = hash
		}
	}
	return merged
}

// IsDiffState compares two snapshots and determines if there are differences.
func IsDiffState(sourceState, currentState FsSnap) bool {
	// Check for new or modified files
	for path, currentHash := range currentState {
		if sourceHash, exists := sourceState[path]; !exists || sourceHash != currentHash {
			return true
		}
	}

	// Check for deleted files
	for path := range sourceState {
		if _, exists := currentState[path]; !exists {
			return true
		}
	}
	return false
}

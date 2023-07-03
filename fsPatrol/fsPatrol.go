// fsPatrol project fsPatrol.go
package fsPatrol

import (
	"goHtmlBuilder/filescaner"
)

type FsSnap map[string]string //[path]hash

func (s FsSnap) GetGhtmlFiles() (result []string) {
	for p := range s {
		if filescaner.IsGhtmlFile(p) {
			result = append(result, p)
		}
	}

	return
}

func GetState() (FsSnap, []error) {
	return getState(filescaner.ScanFS, filescaner.ScanGhtmlFilesOnly)
}

func getState(fScanFs, fScanGhtmlFilesOnly func(path string) (map[string]string, []error)) (FsSnap, []error) {
	var snaps []FsSnap

	for _, v := range []string{"components", "static"} {
		snap, errs := fScanFs(v)
		if len(errs) != 0 {
			return FsSnap{}, errs
		}
		snaps = append(snaps, snap)
	}

	snap, errs := fScanGhtmlFilesOnly(".")
	if len(errs) != 0 {
		return FsSnap{}, errs
	}
	snaps = append(snaps, snap)

	oneSnap := combiningSnap(snaps)
	return oneSnap, nil
}

func combiningSnap(snaps []FsSnap) FsSnap {
	oneSnap := make(map[string]string)
	for _, snap := range snaps {
		for k, v := range snap {
			oneSnap[k] = v
		}
	}
	return oneSnap
}

func IsDiffState(filesSourceState map[string]string, filesCurrentState map[string]string) bool {
	//Сравнивам данные из снимка с текущими файлами и их состоянием
	for k, v := range filesCurrentState {
		hs, ok := filesSourceState[k]
		if !ok {
			return true //new file
		}
		if hs != v {
			return true //change file
		}
	}
	for k := range filesSourceState {
		_, ok := filesCurrentState[k]
		if !ok {
			return true //deleted file
		}
	}
	return false
}

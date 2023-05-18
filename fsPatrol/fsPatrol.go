// fsPatrol project fsPatrol.go
package fsPatrol

import (
	"goHtmlBuilder/filescaner"
	"time"
)



type FsSnap struct {
	SnapTime int64
	SnapFS   map[string]string
}


func GetState() (FsSnap, []error) {
	var snaps []FsSnap

	for _, v := range []string{"components", "static"} {
		snap, errs := filescaner.ScanFS(v)
		if len(errs) != 0 {
			return FsSnap{}, errs
		}
		snaps = append(snaps, FsSnap{SnapTime: 0, SnapFS: snap})
	}

	snap, errs := filescaner.ScanGhtmlFilesOnly(".")
	if len(errs) != 0 {
		return FsSnap{}, errs
	}
	snaps = append(snaps, FsSnap{SnapTime: 0, SnapFS: snap})

	oneSnap := combiningSnap(snaps)
	return oneSnap, nil
}

func combiningSnap(snaps []FsSnap) FsSnap {
	oneSnap := make(map[string]string)
	for _, snap := range snaps {
		for k, v := range snap.SnapFS {
			oneSnap[k] = v
		}
	}
	return FsSnap{SnapTime: time.Now().Unix(), SnapFS: oneSnap}
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
	for k, _ := range filesSourceState {
		_, ok := filesCurrentState[k]
		if !ok {
			return true //delete file
		}
	}
	return false
}



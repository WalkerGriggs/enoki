package enoki

import (
	"io/ioutil"
	"path/filepath"
)

type DiskManager struct {
	Root string
}

var (
	// TODO(walker): Optimize extension lookup
	SupportedContainerMap = map[string]bool{
		"mp4": true,
		"mkv": true,
		"mov": true,
		"avi": true,
		"ogg": true,
		"flv": true,
	}
)

// SourceFiles is used to list all single, source video files with a supported
// file extension
//
// TODO(walker): Checking the file extension is a shortcut for checking the file
// extension itself. This should be updated to use ffprod or some similar
// utility in the future.
func (m *DiskManager) SourceFiles() ([]string, error) {
	files, err := m.readDir(m.Root)
	if err != nil {
		return files, err
	}

	ret := make([]string, 0)
	for _, file := range files {
		ext := filepath.Ext(file)
		if _, ok := SupportedContainerMap[ext]; ok {
			ret = append(ret, file)
		}
	}
	return ret, nil
}

// readDir is used to list all files in a given directory. Note, this filters
// out directories.
func (m *DiskManager) readDir(dir string) ([]string, error) {
	ret := make([]string, 0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return ret, err
	}

	for _, file := range files {
		if !file.IsDir() {
			ret = append(ret, file.Name())
		}
	}
	return ret, nil
}

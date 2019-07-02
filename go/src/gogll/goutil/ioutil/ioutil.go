package ioutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Permission given to all created files and directories
const FilePermission = 0731

func MkdirAll(path string) error {
	if path == "" {
		return nil
	}
	return os.MkdirAll(path, FilePermission)
}

func WriteFile(path string, data []byte) error {
	dir, _ := filepath.Split(path)
	if err := MkdirAll(dir); err != nil {
		return fmt.Errorf("Error creating directory %s: %s", dir, err)
	}
	if err := ioutil.WriteFile(path, data, FilePermission); err != nil {
		return fmt.Errorf("Error writing file %s: %s\n", path, err)
	}
	return nil
}

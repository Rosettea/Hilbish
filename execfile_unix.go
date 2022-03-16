// +build linux darwin

package main

import (
	"path/filepath"
	"os"
)

func findExecutable(path string) (error, string) {
	f, err := os.Stat(path)
	if err != nil {
		return err, ""
	}
	if m := f.Mode(); !m.IsDir() && m & 0111 != 0 {
		return nil, filepath.Base(path)
	}
	return errNotExec, ""
}

// +build windows

package main

import (
	"path/filepath"
	"os"
	"syscall"
)

var bgProcAttr *syscall.SysProcAttr = &syscall.SysProcAttr{
	CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
}

func findExecutable(path string, inPath, dirs bool) error {
	nameExt := filepath.Ext(path)
	pathExts := filepath.SplitList(os.Getenv("PATHEXT"))
	if inPath {
		if nameExt == "" {
			for _, ext := range pathExts {
				_, err := os.Stat(path + ext)
				if err == nil {
					return nil
				}
			}
		} else {
			_, err := os.Stat(path)
			if err == nil {
				if contains(pathExts, nameExt) { return nil }
				return errNotExec
			}
		}
	} else {
		_, err := os.Stat(path)
		if err == nil {
			if contains(pathExts, nameExt) { return nil }
			return errNotExec
		}
	}

	return os.ErrNotExist
}

// +build windows

package main

import (
	"fmt"
	"path/filepath"
	"os"
)

func findExecutable(path string) error {
	nameExt := filepath.Ext(path)

	if nameExt == "" {
		for _, ext := range filepath.SplitList(os.Getenv("PATHEXT")) {
			_, err := os.Stat(path + ext)
			if err == nil {
				return nil
			}
		}
	} else {
		_, err := os.Stat(path)
		return err
	}

	return errNotExec
}

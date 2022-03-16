// +build linux darwin

package main

import (
	"os"
)

func findExecutable(path string) error {
	f, err := os.Stat(path)
	if err != nil {
		return err
	}
	if m := f.Mode(); !m.IsDir() && m & 0111 != 0 {
		return nil
	}
	return errNotExec
}

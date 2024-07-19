//go:build unix

package main

import (
	"os"
	"syscall"
)

var bgProcAttr *syscall.SysProcAttr = &syscall.SysProcAttr{
	Setpgid: true,
}

func findExecutable(path string, inPath, dirs bool) error {
	f, err := os.Stat(path)
	if err != nil {
		return err
	}
	if dirs {
		if m := f.Mode(); m & 0111 != 0 {
			return nil
		}
	} else {
		if m := f.Mode(); !m.IsDir() && m & 0111 != 0 {
			return nil
		}
	}
	return errNotExec
}

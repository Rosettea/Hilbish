// +build windows

package main

import "golang.org/x/sys/windows"

func init() {
	var mode uint32
	windows.GetConsoleMode(windows.Stdout, &mode)
	windows.SetConsoleMode(windows.Stdout, mode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING | windows.ENABLE_VIRTUAL_TERMINAL_INPUT)
}

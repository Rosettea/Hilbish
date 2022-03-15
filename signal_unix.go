// +build darwin linux

package main

import (
	"syscall"
	"os"
	"os/signal"
)

func handleSignals() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGWINCH, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGQUIT)

	for s := range c {
		switch s {
		case os.Interrupt:
			hooks.Em.Emit("signal.sigint")
			if !running && interactive {
				lr.ClearInput()
			}
		case syscall.SIGWINCH:
			hooks.Em.Emit("signal.resize")
			if !running && interactive {
				lr.Resize()
			}
		case syscall.SIGUSR1: hooks.Em.Emit("signal.sigusr1")
		case syscall.SIGUSR2: hooks.Em.Emit("signal.sigusr2")
		}
	}
}

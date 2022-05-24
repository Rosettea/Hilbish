// +build darwin linux

package main

import (
	"syscall"
	"os"
	"os/signal"
)

func handleSignals() {
	c := make(chan os.Signal)
	signal.Ignore(syscall.SIGTTOU, syscall.SIGTTIN, syscall.SIGTSTP)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGWINCH, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGQUIT)

	for s := range c {
		switch s {
		case os.Interrupt: hooks.Em.Emit("signal.sigint")
		case syscall.SIGTERM: exit(0)
		case syscall.SIGWINCH: hooks.Em.Emit("signal.resize")
		case syscall.SIGUSR1: hooks.Em.Emit("signal.sigusr1")
		case syscall.SIGUSR2: hooks.Em.Emit("signal.sigusr2")
		}
	}
}

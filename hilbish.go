// Here is the core api for the hilbi shell itself
// Basically, stuff about the shell itself and other functions
// go here.
package main

type Hilbish struct {
	Version string `luar:"version"` // Hilbish's version
	User string `luar:"user"` // Name of the user
	Hostname string `luar:"hostname"`
}


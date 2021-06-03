// +build hilbiline

// Here we define a generic interface for readline and hilbiline,
// making them interchangable during build time
// this is hilbiline's, as is obvious by the filename
package main

import "github.com/Rosettea/Hilbiline"

type LineReader struct {
	hl *hilbiline.HilbilineState
}

// other gophers might hate this naming but this is local, shut up
func NewLineReader(prompt string) *LineReader {
	hl := hilbiline.New(prompt)

	return &LineReader{
		&hl,
	}
}

func (lr *LineReader) Read() (string, error) {
	return lr.hl.Read()
}

func (lr *LineReader) SetPrompt(prompt string) {
	lr.hl.SetPrompt(prompt)
}

func (lr *LineReader) AddHistory(cmd string) {
	return
}


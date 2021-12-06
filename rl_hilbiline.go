// +build hilbiline

package main

// Here we define a generic interface for readline and hilbiline,
// making them interchangable during build time
// this is hilbiline's, as is obvious by the filename

import "github.com/Rosettea/Hilbiline"

type lineReader struct {
	hl *hilbiline.HilbilineState
}

// other gophers might hate this naming but this is local, shut up
func newLineReader(prompt string) *lineReader {
	hl := hilbiline.New(prompt)

	return &lineReader{
		&hl,
	}
}

func (lr *lineReader) Read() (string, error) {
	return lr.hl.Read()
}

func (lr *lineReader) SetPrompt(prompt string) {
	lr.hl.SetPrompt(prompt)
}

func (lr *lineReader) AddHistory(cmd string) {
	return
}

func (lr *lineReader) ClearInput() {
	return
}

func (lr *lineReader) Resize() {
	return
}


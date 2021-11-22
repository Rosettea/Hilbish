//go:build !hilbiline
// +build !hilbiline

package main

// Here we define a generic interface for readline and hilbiline,
// making them interchangable during build time
// this is normal readline

import "github.com/bobappleyard/readline"

type LineReader struct {
	Prompt string
}

// other gophers might hate this naming but this is local, shut up
func NewLineReader(prompt string) *LineReader {
	readline.Completer = readline.FilenameCompleter
	readline.LoadHistory(defaultHistPath)

	return &LineReader{
		Prompt: prompt,
	}
}

func (lr *LineReader) Read() (string, error) {
	hooks.Em.Emit("command.precmd", nil)
	return readline.String(lr.Prompt)
}

func (lr *LineReader) SetPrompt(prompt string) {
	lr.Prompt = prompt
}

func (lr *LineReader) AddHistory(cmd string) {
	readline.AddHistory(cmd)
	readline.SaveHistory(defaultHistPath)
}

func (lr *LineReader) ClearInput() {
	readline.ReplaceLine("", 0)
	readline.RefreshLine()
}

func (lr *LineReader) Resize() {
	readline.Resize()
}

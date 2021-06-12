package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
func TestRunInputSh(t *testing.T) {
	LuaInit()
	cmd := "echo 'hello'"
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	RunInput(cmd)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	if out != "hello\n" {
		t.Fatalf("Expected 'hello', found %s", out)
	}
}

func TestRunInputLua(t *testing.T) {
	LuaInit()
	cmd := "print('hello')"
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	RunInput(cmd)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	if out != "hello\n" {
		t.Fatalf("Expected 'hello', found %s", out)
	}
}

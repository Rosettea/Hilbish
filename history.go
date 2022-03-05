package main

import (
	"errors"
	"io/fs"
	"os"
	"strings"
)

type fileHistory struct {
	items []string
	f *os.File
}

func newFileHistory() (*fileHistory, error) {
	data, err := os.ReadFile(defaultHistPath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}
	
	var itms []string
	for _, l := range strings.Split(string(data), "\n") {
		itms = append(itms, l)
	}
	f, err := os.OpenFile(defaultHistPath, os.O_RDWR | os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	
	fh := &fileHistory{
		items: itms,
		f: f,
	}
	
	return fh, nil
}

func (h *fileHistory) Write(line string) (int, error) {
	_, err := h.f.WriteString(line + "\n")
	if err != nil {
		return 0, err
	}
	h.f.Sync()
	
	h.items = append(h.items, line)
	return len(h.items), nil
}

func (h *fileHistory) GetLine(idx int) (string, error) {
	if len(h.items) == 0 {
		return "", nil
	}
	if idx == -1 { // this should be fixed readline side
		return "", nil
	}
	return h.items[idx], nil
}

func (h *fileHistory) Len() int {
	return len(h.items)
}

func (h *fileHistory) Dump() interface{} {
	return nil
}

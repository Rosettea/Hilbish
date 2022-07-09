package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type fileHistory struct {
	items []string
	f *os.File
}

func newFileHistory(path string) *fileHistory {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			panic(err)
		}
	}

	itms := []string{""}
	lines := strings.Split(string(data), "\n")
	for i, l := range lines {
		if i == len(lines) - 1 {
			continue
		}
		itms = append(itms, l)
	}
	f, err := os.OpenFile(path, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	fh := &fileHistory{
		items: itms,
		f: f,
	}

	return fh
}

func (h *fileHistory) Write(line string) (int, error) {
	if line == "" {
		return len(h.items), nil
	}

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
	return h.items
}

func (h *fileHistory) clear() {
	h.items = []string{}
	h.f.Truncate(0)
	h.f.Sync()
}

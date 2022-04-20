package main

import (
	"path/filepath"
	"strings"
	"os"
)

func fileComplete(query, ctx string, fields []string) ([]string, string) {
	return matchPath(query)
}

func binaryComplete(query, ctx string, fields []string) ([]string, string) {
	var completions []string

	prefixes := []string{"./", "../", "/", "~/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(query, prefix) {
			fileCompletions, filePref := matchPath(query)
			if len(fileCompletions) != 0 {
				for _, f := range fileCompletions {
					fullPath, _ := filepath.Abs(expandHome(query + strings.TrimPrefix(f, filePref)))
					if err := findExecutable(fullPath, false, true); err != nil {
						continue
					}
					completions = append(completions, f)
				}
			}
			return completions, filePref
		}
	}

	// filter out executables, but in path
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		// print dir to stderr for debugging
		// search for an executable which matches our query string
		if matches, err := filepath.Glob(filepath.Join(dir, query + "*")); err == nil {
			// get basename from matches
			for _, match := range matches {
				// check if we have execute permissions for our match
				err := findExecutable(match, true, false)
				if err != nil {
					continue
				}
				// get basename from match
				name := filepath.Base(match)
				// add basename to completions
				completions = append(completions, name)
			}
		}
	}

	// add lua registered commands to completions
	for cmdName := range commands {
		if strings.HasPrefix(cmdName, query) {
			completions = append(completions, cmdName)
		}
	}

	completions = removeDupes(completions)

	return completions, query
}

func matchPath(query string) ([]string, string) {
	var entries []string
	var baseName string

	path, _ := filepath.Abs(expandHome(filepath.Dir(query)))
	if string(query) == "" {
		// filepath base below would give us "."
		// which would cause a match of only dotfiles
		path, _ = filepath.Abs(".")
	} else if !strings.HasSuffix(query, string(os.PathSeparator)) {
		baseName = filepath.Base(query)
	}

	files, _ := os.ReadDir(path)
	for _, file := range files {
		if strings.HasPrefix(strings.ToLower(file.Name()), strings.ToLower(baseName)) {
			entry := file.Name()
			if file.IsDir() {
				entry = entry + string(os.PathSeparator)
			}
			entry = escapeFilename(entry)
			entries = append(entries, entry)
		}
	}

	return entries, baseName
}

func escapeFilename(fname string) string {
	args := []string{
		"\"", "\\\"",
		"'", "\\'",
		"`", "\\`",
		" ", "\\ ",
		"(", "\\(",
		")", "\\)",
		"[", "\\[",
		"]", "\\]",
		"$", "\\$",
		"&", "\\&",
		"*", "\\*",
		">", "\\>",
		"<", "\\<",
		"|", "\\|",
	}

	r := strings.NewReplacer(args...)
	return r.Replace(fname)
}


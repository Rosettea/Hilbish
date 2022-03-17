package main

import (
	"path/filepath"
	"strings"
	"os"
)

func fileComplete(query, ctx string, fields []string) []string {
	var completions []string

	prefixes := []string{"./", "../", "/", "~/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(query, prefix) {
			completions, _ = matchPath(strings.Replace(query, "~", curuser.HomeDir, 1), query)
		}
	}

	if len(completions) == 0 && len(fields) > 1 {
		completions, _ = matchPath("./" + query, query)
	}

	return completions
}

func binaryComplete(query, ctx string, fields []string) ([]string, string) {
	var completions []string

	prefixes := []string{"./", "../", "/", "~/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(query, prefix) {
			fileCompletions := fileComplete(query, ctx, fields)
			if len(fileCompletions) != 0 {
				for _, f := range fileCompletions {
					name := strings.Replace(query + f, "~", curuser.HomeDir, 1)
					if info, err := os.Stat(name); err == nil && info.Mode().Perm() & 0100 == 0 {
						continue
					}
					completions = append(completions, f)
				}
			}
			return completions, ""
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
				if info, err := os.Stat(match); err == nil && info.Mode().Perm() & 0100 == 0 {
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

func matchPath(path, pref string) ([]string, error) {
	var entries []string
	matches, err := filepath.Glob(path + "*")
	if err == nil {
		args := []string{
			"\"", "\\\"",
			"'", "\\'",
			"`", "\\`",
			" ", "\\ ",
			"(", "\\(",
			")", "\\)",
			"[", "\\[",
			"]", "\\]",
		}

		r := strings.NewReplacer(args...)
		for _, match := range matches {
			name := filepath.Base(match)
			p := filepath.Base(pref)
			if pref == "" || pref == "./" {
				p = ""
			}
			name = strings.TrimPrefix(name, p)
			matchFull, _ := filepath.Abs(match)
			if info, err := os.Stat(matchFull); err == nil && info.IsDir() {
				name = name + string(os.PathSeparator)
			}
			name = r.Replace(name)
			entries = append(entries, name)
		}
	}

	return entries, err
}

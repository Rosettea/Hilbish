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

func matchPath(path, pref string) ([]string, error) {
	var entries []string
	matches, err := filepath.Glob(path + "*")
	if err == nil {
		args := []string{
			"\"", "\\\"",
			"'", "\\'",
			"`", "\\`",
			" ", "\\ ",
		}

		r := strings.NewReplacer(args...)
		for _, match := range matches {
			name := filepath.Base(match)
			p := filepath.Base(pref)
			if pref == "" {
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

package main

import (
	"fmt"
	"path/filepath"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
	"os"
)

// feel free to clean this up
// it works, dont really care about the code
func main() {
	fset := token.NewFileSet()

	dirs := []string{"./"}
	filepath.Walk("golibs/", func (path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		dirs = append(dirs, "./" + path)
		return nil
	})

	pkgs := make(map[string]*ast.Package)
	for _, path := range dirs {
		d, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Println(err)
			return
		}
		for k, v := range d {
			pkgs[k] = v
		}
	}

	prefix := map[string]string{
		"main": "hsh",
		"hilbish": "hl",
		"fs": "f",
		"commander": "c",
		"bait": "b",
		"terminal": "term",
	}
	docs := make(map[string][]string)

	for l, f := range pkgs {
		p := doc.New(f, "./", doc.AllDecls)
		for _, t := range p.Funcs {
			mod := l
			if strings.HasPrefix(t.Name, "hl") { mod = "hilbish" }
			if !strings.HasPrefix(t.Name, prefix[mod]) || t.Name == "Loader" { continue }
			parts := strings.Split(t.Doc, "\n")
			funcsig := parts[0]
			doc := parts[1:]

			docs[mod] = append(docs[mod], funcsig + " > " + strings.Join(doc, "\n"))
		}
		for _, t := range p.Types {
			for _, m := range t.Methods {
				if !strings.HasPrefix(m.Name, prefix[l]) || m.Name == "Loader" { continue }
				parts := strings.Split(m.Doc, "\n")
				funcsig := parts[0]
				doc := parts[1:]

				docs[l] = append(docs[l], funcsig + " > " + strings.Join(doc, "\n"))
			}
		}
	}

	for mod, v := range docs {
		if mod == "main" { mod = "global" }
		os.Mkdir("docs", 0777)
		f, _ := os.Create("docs/" + mod + ".txt")
		f.WriteString(strings.Join(v, "\n") + "\n")
	}
}

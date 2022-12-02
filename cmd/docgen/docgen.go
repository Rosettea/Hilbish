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

var header = `---
name: Module %s
description: %s
layout: apidoc
---

`

type emmyPiece struct {
	FuncName string
	Docs []string
	Params []string // we only need to know param name to put in function
}

type module struct {
	Docs []docPiece
	ShortDescription string
	Description string
}
type docPiece struct {
	Doc []string
	FuncSig string
	FuncName string
}

// feel free to clean this up
// it works, dont really care about the code
func main() {
	fset := token.NewFileSet()
	os.Mkdir("docs", 0777)
	os.Mkdir("docs/api", 0777)
	os.Mkdir("emmyLuaDocs", 0777)

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
		"hilbish": "hl",
		"fs": "f",
		"commander": "c",
		"bait": "b",
		"terminal": "term",
	}
	docs := make(map[string]module)
	emmyDocs := make(map[string][]emmyPiece)

	for l, f := range pkgs {
		p := doc.New(f, "./", doc.AllDecls)
		pieces := []docPiece{}
		for _, t := range p.Funcs {
			mod := l
			if strings.HasPrefix(t.Name, "hl") { mod = "hilbish" }
			if !strings.HasPrefix(t.Name, prefix[mod]) || t.Name == "Loader" { continue }
			parts := strings.Split(strings.TrimSpace(t.Doc), "\n")
			funcsig := parts[0]
			doc := parts[1:]
			funcdoc := []string{}
			em := emmyPiece{FuncName: strings.TrimPrefix(t.Name, prefix[mod])}
			for _, d := range doc {
				if strings.HasPrefix(d, "---") {
					emmyLine := strings.TrimSpace(strings.TrimPrefix(d, "---"))
					emmyLinePieces := strings.Split(emmyLine, " ")
					emmyType := emmyLinePieces[0]
					if emmyType == "@param" {
						em.Params = append(em.Params, emmyLinePieces[1])
					}
					if emmyType == "@vararg" {
						em.Params = append(em.Params, "...") // add vararg
					}
					em.Docs = append(em.Docs, d)
				} else {
					funcdoc = append(funcdoc, d)
				}
			}
			
			dps := docPiece{
				Doc: funcdoc,
				FuncSig: funcsig,
				FuncName: strings.TrimPrefix(t.Name, prefix[mod]),
			}
			
			pieces = append(pieces, dps)
			emmyDocs[mod] = append(emmyDocs[mod], em)
		}
		for _, t := range p.Types {
			for _, m := range t.Methods {
				if !strings.HasPrefix(m.Name, prefix[l]) || m.Name == "Loader" { continue }
				parts := strings.Split(strings.TrimSpace(m.Doc), "\n")
				funcsig := parts[0]
				doc := parts[1:]
				funcdoc := []string{}
				em := emmyPiece{FuncName: strings.TrimPrefix(m.Name, prefix[l])}
				for _, d := range doc {
					if strings.HasPrefix(d, "---") {
						emmyLine := strings.TrimSpace(strings.TrimPrefix(d, "---"))
						emmyLinePieces := strings.Split(emmyLine, " ")
						emmyType := emmyLinePieces[0]
						if emmyType == "@param" {
							em.Params = append(em.Params, emmyLinePieces[1])
						}
						if emmyType == "@vararg" {
							em.Params = append(em.Params, "...") // add vararg
						}
						em.Docs = append(em.Docs, d)
					} else {
						funcdoc = append(funcdoc, d)
					}
				}
				dps := docPiece{
					Doc: funcdoc,
					FuncSig: funcsig,
					FuncName: strings.TrimPrefix(m.Name, prefix[l]),
				}

				pieces = append(pieces, dps)
				emmyDocs[l] = append(emmyDocs[l], em)
			}
		}

		descParts := strings.Split(strings.TrimSpace(p.Doc), "\n")
		shortDesc := descParts[0]
		desc := descParts[1:]
		docs[l] = module{
			Docs: pieces,
			ShortDescription: shortDesc,
			Description: strings.Join(desc, "\n"),
		}
	}

	for mod, v := range docs {
		if mod == "main" { continue }
		f, _ := os.Create("docs/api/" + mod + ".md")
		f.WriteString(fmt.Sprintf(header, mod, v.ShortDescription))
		f.WriteString(fmt.Sprintf("## Introduction\n%s\n\n## Functions\n", v.Description))
		for _, dps := range v.Docs {
			f.WriteString(fmt.Sprintf("### %s\n", dps.FuncSig))
			for _, doc := range dps.Doc {
				if !strings.HasPrefix(doc, "---") {
					f.WriteString(doc + "\n")
				}
			}
			f.WriteString("\n")
		}
	}
	
	for mod, v := range emmyDocs {
		if mod == "main" { continue }
		f, _ := os.Create("emmyLuaDocs/" + mod + ".lua")
		f.WriteString("--- @meta\n\nlocal " + mod + " = {}\n\n")
		for _, em := range v {
			var funcdocs []string
			for _, dps := range docs[mod].Docs {
				if dps.FuncName == em.FuncName {
					funcdocs = dps.Doc
				}
			}
			f.WriteString("--- " + strings.Join(funcdocs, "\n--- ") + "\n")
			if len(em.Docs) != 0 {
				f.WriteString(strings.Join(em.Docs, "\n") + "\n")
			}
			f.WriteString("function " + mod + "." + em.FuncName + "(" + strings.Join(em.Params, ", ") + ") end\n\n")
		}
		f.WriteString("return " + mod + "\n")
	}
}

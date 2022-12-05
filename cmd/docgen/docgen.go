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
	"sync"
)

var header = `---
name: %s %s
description: %s
layout: apidoc
---

`

type emmyPiece struct {
	DocPiece *docPiece
	Annotations []string
	Params []string // we only need to know param name to put in function
	FuncName string
}

type module struct {
	Docs []docPiece
	Properties []docPiece
	ShortDescription string
	Description string
	ParentModule string
	HasInterfaces bool
}

type docPiece struct {
	Doc []string
	FuncSig string
	FuncName string
	Interfacing string
	ParentModule string
	GoFuncName string
	IsInterface bool
	IsMember bool
	Properties []docPiece
}

type tag struct {
	id string
	fields []string
}

var docs = make(map[string]module)
var interfaceDocs = make(map[string]module)
var emmyDocs = make(map[string][]emmyPiece)
var prefix = map[string]string{
	"main": "hl",
	"hilbish": "hl",
	"fs": "f",
	"commander": "c",
	"bait": "b",
	"terminal": "term",
}

func setupDoc(mod string, fun *doc.Func) *docPiece {
	docs := strings.TrimSpace(fun.Doc)
	inInterface := strings.HasPrefix(docs, "#interface")
	if (!strings.HasPrefix(fun.Name, prefix[mod]) && !inInterface) || (strings.ToLower(fun.Name) == "loader" && !inInterface) {
		return nil
	}

	pts := strings.Split(docs, "\n")
	parts := []string{}
	tags := make(map[string][]tag)
	for _, part := range pts {
		if strings.HasPrefix(part, "#") {
			tagParts := strings.Split(strings.TrimPrefix(part, "#"), " ")
			if tags[tagParts[0]] == nil {
				var id string
				if len(tagParts) > 1 {
					id = tagParts[1]
				}
				tags[tagParts[0]] = []tag{
					{id: id},
				}
				if len(tagParts) >= 2 {
					tags[tagParts[0]][0].fields = tagParts[2:]
				}
			} else {
				fleds := []string{}
				if len(tagParts) >= 2 {
					fleds = tagParts[2:]
				}
				tags[tagParts[0]] = append(tags[tagParts[0]], tag{
					id: tagParts[1],
					fields: fleds,
				})
			}
		} else {
			parts = append(parts, part)
		}
	}

	var interfaces string
	funcsig := parts[0]
	doc := parts[1:]
	funcName := strings.TrimPrefix(fun.Name, prefix[mod])
	funcdoc := []string{}

	if inInterface {
		interfaces = tags["interface"][0].id
		funcName = interfaces + "." + strings.Split(funcsig, "(")[0]
	}
	em := emmyPiece{FuncName: funcName}

	// manage properties
	properties := []docPiece{}
	for _, tag := range tags["property"] {
		properties = append(properties, docPiece{
			FuncName: tag.id,
			Doc: tag.fields,
		})
	}

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
			em.Annotations = append(em.Annotations, d)
		} else {
			funcdoc = append(funcdoc, d)
		}
	}

	var isMember bool
	if tags["member"] != nil {
		isMember = true
	}
	var parentMod string
	if inInterface {
		parentMod = mod
	}
	dps := &docPiece{
		Doc: funcdoc,
		FuncSig: funcsig,
		FuncName: funcName,
		Interfacing: interfaces,
		GoFuncName: strings.ToLower(fun.Name),
		IsInterface: inInterface,
		IsMember: isMember,
		ParentModule: parentMod,
		Properties: properties,
	}
	if strings.HasSuffix(dps.GoFuncName, strings.ToLower("loader")) {
		dps.Doc = parts
	}
	em.DocPiece = dps
			
	emmyDocs[mod] = append(emmyDocs[mod], em)
	return dps
}

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

	interfaceModules := make(map[string]*module)
	for l, f := range pkgs {
		p := doc.New(f, "./", doc.AllDecls)
		pieces := []docPiece{}
		mod := l
		if mod == "main" {
			mod = "hilbish"
		}
		var hasInterfaces bool
		for _, t := range p.Funcs {
			piece := setupDoc(mod, t)
			if piece == nil {
				continue
			}

			pieces = append(pieces, *piece)
			if piece.IsInterface {
				hasInterfaces = true
			}
		}
		for _, t := range p.Types {
			for _, m := range t.Methods {
				piece := setupDoc(mod, m)
				if piece == nil {
					continue
				}

				pieces = append(pieces, *piece)
				if piece.IsInterface {
					hasInterfaces = true
				}
			}
		}

		descParts := strings.Split(strings.TrimSpace(p.Doc), "\n")
		shortDesc := descParts[0]
		desc := descParts[1:]
		filteredPieces := []docPiece{}
		for _, piece := range pieces {
			if !piece.IsInterface {
				filteredPieces = append(filteredPieces, piece)
				continue
			}

			modname := piece.ParentModule + "." + piece.Interfacing
			if interfaceModules[modname] == nil {
				interfaceModules[modname] = &module{
					ParentModule: piece.ParentModule,
				}
			}

			if strings.HasSuffix(piece.GoFuncName, strings.ToLower("loader")) {
				shortDesc := piece.Doc[0]
				desc := piece.Doc[1:]
				interfaceModules[modname].ShortDescription = shortDesc
				interfaceModules[modname].Description = strings.Join(desc, "\n")
				interfaceModules[modname].Properties = piece.Properties
				continue
			}
			interfaceModules[modname].Docs = append(interfaceModules[modname].Docs, piece)
		}

		docs[mod] = module{
			Docs: filteredPieces,
			ShortDescription: shortDesc,
			Description: strings.Join(desc, "\n"),
			HasInterfaces: hasInterfaces,
		}
	}

	for key, mod := range interfaceModules {
		docs[key] = *mod
	}

	var wg sync.WaitGroup
	wg.Add(len(docs) * 2)

	for mod, v := range docs {
		docPath := "docs/api/" + mod + ".md"
		if v.HasInterfaces {
			os.Mkdir("docs/api/" + mod, 0777)
			os.Remove(docPath) // remove old doc path if it exists
			docPath = "docs/api/" + mod + "/_index.md"
		}
		if v.ParentModule != "" {
			docPath = "docs/api/" + v.ParentModule + "/" + mod + ".md"
		}

		go func(modname, docPath string, modu module) {
			defer wg.Done()
			modOrIface := "Module"
			if modu.ParentModule != "" {
				modOrIface = "Interface"
			}

			f, _ := os.Create(docPath)
			f.WriteString(fmt.Sprintf(header, modOrIface, modname, modu.ShortDescription))
			f.WriteString(fmt.Sprintf("## Introduction\n%s\n\n", modu.Description))
			if len(modu.Properties) != 0 {
				f.WriteString("## Properties\n")
				for _, dps := range modu.Properties {
					f.WriteString(fmt.Sprintf("- `%s`: ", dps.FuncName))
					f.WriteString(strings.Join(dps.Doc, " "))
					f.WriteString("\n")
				}
			}
			if len(modu.Docs) != 0 {
				f.WriteString("## Functions\n")
			}
			for _, dps := range modu.Docs {
				f.WriteString(fmt.Sprintf("### %s\n", dps.FuncSig))
				for _, doc := range dps.Doc {
					if !strings.HasPrefix(doc, "---") {
						f.WriteString(doc + "\n")
					}
				}
				f.WriteString("\n")
			}
		}(mod, docPath, v)

		go func(md, modname string, modu module) {
			defer wg.Done()

			if modu.ParentModule != "" {
				return
			}

			ff, _ := os.Create("emmyLuaDocs/" + modname + ".lua")
			ff.WriteString("--- @meta\n\nlocal " + modname + " = {}\n\n")
			for _, em := range emmyDocs[modname] {
				if strings.HasSuffix(em.DocPiece.GoFuncName, strings.ToLower("loader")) {
					continue
				}

				dps := em.DocPiece
				funcdocs := dps.Doc
				ff.WriteString("--- " + strings.Join(funcdocs, "\n--- ") + "\n")
				if len(em.Annotations) != 0 {
					ff.WriteString(strings.Join(em.Annotations, "\n") + "\n")
				}
				accessor := "."
				if dps.IsMember {
					accessor = ":"
				}
				signature := strings.Split(dps.FuncSig, " ->")[0]
				var intrface string
				if dps.IsInterface {
					intrface = "." + dps.Interfacing
				}
				ff.WriteString("function " + modname + intrface + accessor + signature + " end\n\n")
			}
			ff.WriteString("return " + modname + "\n")
		}(mod, mod, v)
	}
	wg.Wait()
}

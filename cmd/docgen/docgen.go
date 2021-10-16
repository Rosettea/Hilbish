package main

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()

	d, err := parser.ParseDir(fset, "./", nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, f := range d {
		p := doc.New(f, "./", 0)

		for _, t := range p.Funcs {
			fmt.Println("	type", t.Name)
			fmt.Println("		docs:", t.Doc)
		}
	}
}

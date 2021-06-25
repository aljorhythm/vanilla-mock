package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

type ParsedFile struct {
	file    *ast.File
	iface   *types.Interface
	fileSet *token.FileSet
}

func (p ParsedFile) Find(name string) (*types.Interface, error) {
	return nil, nil
}

func (p ParsedFile) PrintAst() {
	ast.Print(p.fileSet, p.file)
}

type InterfaceCollector struct {
	collector func(p *types.Interface, name string)
}

func (i *InterfaceCollector) collect(p *types.Interface, name string) {
	if i.collector != nil {
		i.collector(p, name)
	}
}

func (i *InterfaceCollector) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.TypeSpec:
		switch n.Type.(type) {
		case *ast.InterfaceType:
			interfaceType := n.Type.(*ast.InterfaceType)
			log.Printf("ctype %s %#v", n.Name, interfaceType)
		}
	}
	return i
}

func Parse(filePath string) (*ParsedFile, error) {
	parsed := ParsedFile{}
	fset := token.NewFileSet()
	parsed.fileSet = fset
	var err error
	parsed.file, err = parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	ast.Walk(&InterfaceCollector{}, parsed.file)
	return &parsed, nil
}

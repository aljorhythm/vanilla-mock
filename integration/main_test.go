package integration

import (
	"context"
	"github.com/aljorhythm/vanilla-mock/generator"
	"github.com/aljorhythm/vanilla-mock/parser"
	"github.com/aljorhythm/vanilla-mock/test"
	"go/types"
	"golang.org/x/tools/go/packages"
	"testing"
)

func Test_ParseWalkGenerate(t *testing.T) {
	filePath := test.GetAbsPath("test/vaniller.go")
	t.Logf(filePath)
	name := "Vaniller"
	parsed, err := parser.Parse(filePath)
	if err != nil {
		t.Error(err)
		return
	}
	parsed.Find(name)
}

func Test_TestMain(t *testing.T) {
	interfaceName := "Vaniller"
	dir := test.GetTestRoot()
	t.Logf("dir %s", dir)
	pkgs, err := packages.Load(&packages.Config{
		Context: context.Background(),
		Dir:     dir,
		Mode:    packages.NeedFiles | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypes,
	})

	t.Logf("pkgs %#v", len(pkgs))

	if err != nil {
		t.Error(err)
		return
	}

	pkg := pkgs[0]
	t.Logf("pkg %#v", pkg)

	scope := pkg.Types.Scope()
	t.Logf("pkg scope %#v", scope)
	t.Logf("scope names %#v", scope.Names())

	obj := scope.Lookup(interfaceName)

	if "1" == "1" {
		t.Logf("obj %#v", obj)

		//return
	}

	typ, _ := obj.Type().(*types.Named)

	iface, _ := typ.Underlying().(*types.Interface)

	v, err := generator.GenerateVanillaMock(iface, interfaceName)
	t.Logf(v.Output())
}

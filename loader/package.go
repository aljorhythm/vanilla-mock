package loader

import (
	"context"
	"errors"
	"go/types"
	"golang.org/x/tools/go/packages"
)

func LoadInterface(ctx context.Context, dir string, ifaceName string) (*types.Interface, error) {
	pkgs, err := packages.Load(&packages.Config{
		Context: ctx,
		Dir:     dir,
		Mode:    packages.NeedFiles | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypes,
	})

	if err != nil {
		return nil, err
	}

	pkg := pkgs[0]

	scope := pkg.Types.Scope()

	obj := scope.Lookup(ifaceName)

	if obj == nil {
		return nil, errors.New("interface_not_found")
	}

	objType := obj.Type()

	if objType == nil {
		return nil, errors.New("obj_type_not_found")
	}

	typ, ok := objType.(*types.Named)
	if !ok {
		return nil, errors.New("not_named_type")
	}
	iface, ok := typ.Underlying().(*types.Interface)

	if !ok {
		return nil, errors.New("not_interface")
	}

	return iface, nil
}

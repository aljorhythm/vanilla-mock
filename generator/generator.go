package generator

import (
	"fmt"
	"go/types"
	"strings"
)

func GenerateVanillaMock(iface *types.Interface, ifaceName string) (*VanillaMockStructOutput, error) {
	v := VanillaMockStructOutput{
		iface:             iface,
		ifaceName:         ifaceName,
		localizationCache: map[string]string{},
		packagePathToName: map[string]string{},
		nameToPackagePath: map[string]string{},
	}

	for i := 0; i < iface.NumMethods(); i++ {
		m := iface.Method(i)
		v.parseMethod(m)
	}

	return &v, nil
}

type VanillaMockStructOutput struct {
	iface             *types.Interface
	ifaceName         string
	pkg               string
	fields            []string
	impls             []string
	packagePathToName map[string]string
	nameToPackagePath map[string]string
	localizationCache map[string]string
	packageRoots      []string
}

func (v *VanillaMockStructOutput) mockStructName() string {
	return v.ifaceName + "VMock"
}

func (v *VanillaMockStructOutput) ifaceFirstLetter() string {
	return strings.ToLower(string(v.ifaceName[0]))
}

func (v *VanillaMockStructOutput) addField(fn string) {
	v.fields = append(v.fields, fn)
}

func (v *VanillaMockStructOutput) addImpl(impl string) {
	v.impls = append(v.impls, impl)
}

func (v *VanillaMockStructOutput) Output() string {
	template := `package mock

type %s struct %s
%s
`
	name := v.mockStructName()
	fnsStr := strings.Join(append([]string{"{"}, v.fields...), "\n\t")
	implsStr := strings.Join(append([]string{"}"}, v.impls...), "\n\n")
	return fmt.Sprintf(template, name, fnsStr, implsStr)
}

func FuncSig(fn *types.Func) *types.Signature {
	sig, ok := fn.Type().(*types.Signature)
	if !ok {
		panic("failed type assert into types.Signature")
	}
	return sig
}

func (v *VanillaMockStructOutput) parseMethod(m *types.Func) {
	sig := FuncSig(m)

	params := sig.Params()

	newFnParams := []string{}
	mockFnInputs := []string{}
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		pName := param.Name()
		var pType string

		if i == params.Len()-1 && sig.Variadic() {
			switch t := param.Type().(type) {
			case *types.Slice:
				pType = "..." + t.Elem().String()
			default:
				panic("bad variadic type!")
			}
		} else {
			pType = v.renderType(param.Type())
		}

		if pName == "" {
			pName = fmt.Sprintf("%s%d", strings.ToLower(string(pType[0])), i)
		}
		newParam := fmt.Sprintf("%s %s", pName, pType)
		newFnParams = append(newFnParams, newParam)

		input := pName
		if i == params.Len()-1 && sig.Variadic() {
			input = input + "..."
		}

		mockFnInputs = append(mockFnInputs, input)
	}

	newSigParams := strings.Join(newFnParams, ", ")
	fName := m.Name() + "Fn"
	returns := []string{}

	for i := 0; i < sig.Results().Len(); i++ {
		r := sig.Results().At(i)
		returns = append(returns, v.renderType(r.Type()))
	}

	rets := fmt.Sprintf("(%s)", strings.Join(returns, ", "))

	v.addField(fmt.Sprintf("%s func(%s) %s", fName, newSigParams, rets))

	passArgs := strings.Join(mockFnInputs, ", ")
	receiver := v.mockStructName()
	iLetter := v.ifaceFirstLetter()

	ret := "return "
	if sig.Results().Len() == 0 {
		ret = ""
	}

	impl := fmt.Sprintf(`func (%s %s) %s(%s) %s {
	%s%s.%s(%s)
}`, iLetter, receiver, m.Name(), newSigParams, rets, ret, iLetter, fName, passArgs)
	v.addImpl(impl)
}

func (v *VanillaMockStructOutput) addPackageImport(pkg *types.Package) string {
	return v.addPackageImportWithName(pkg.Path(), pkg.Name())
}

func (vo *VanillaMockStructOutput) importNameExists(name string) bool {
	_, nameExists := vo.nameToPackagePath[name]
	return nameExists
}

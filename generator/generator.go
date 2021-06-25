package generator

import (
	"errors"
	"fmt"
	"github.com/aljorhythm/vanilla-mock/parser"
	"go/types"
	"strings"
)

func GenerateVanillaMockFromFile(interfacePath string, interfaceName string) (*VanillaMockStructOutput, error) {
	parsed, err := parser.Parse(interfacePath)
	iface, err := parsed.Find(interfaceName)

	if err != nil {
		return nil, err
	}

	if iface == nil {
		return nil, errors.New("iface is nil")
	}

	return GenerateVanillaMock(iface, interfaceName)
}

func GenerateVanillaMock(iface *types.Interface, ifaceName string) (*VanillaMockStructOutput, error) {
	v := VanillaMockStructOutput{
		iface:     iface,
		ifaceName: ifaceName,
	}

	for i := 0; i < iface.NumMethods(); i++ {
		m := iface.Method(i)
		v.parseMethod(m)
	}

	return &v, nil
}

type VanillaMockStructOutput struct {
	iface     *types.Interface
	ifaceName string
	pkg       string
	fields    []string
	impls     []string
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
	template := `type %s struct %s
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
	sigString := sig.String()
	fName := m.Name() + "Fn"
	v.addField(fmt.Sprintf("%s %s", fName, sigString))

	params := sig.Params()

	newFnParams := []string{}
	mockFnInputs := []string{}
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		pName := param.Name()
		pType := param.Type().String()

		if i == params.Len()-1 && sig.Variadic() {
			switch t := param.Type().(type) {
			case *types.Slice:
				pType = "..." + t.Elem().String()
			default:
				panic("bad variadic type!")
			}
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

	newSig := strings.Join(newFnParams, ", ")
	passArgs := strings.Join(mockFnInputs, ", ")
	rets := sig.Results().String()
	receiver := v.mockStructName()
	iLetter := v.ifaceFirstLetter()
	ret := "return "
	if sig.Results().Len() == 0 {
		ret = ""
	}

	impl := fmt.Sprintf(`func (%s %s) %s(%s) %s {
	%s%s.%s(%s)
}`, iLetter, receiver, m.Name(), newSig, rets, ret, iLetter, fName, passArgs)
	v.addImpl(impl)
}

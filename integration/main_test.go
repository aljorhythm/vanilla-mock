package integration

import (
	"context"
	"github.com/aljorhythm/vanilla-mock/generator"
	"github.com/aljorhythm/vanilla-mock/loader"
	"github.com/aljorhythm/vanilla-mock/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LoadInterfaceAndGenerate(t *testing.T) {
	interfaceName := "Vaniller"
	dir := test.GetTestRoot()
	iface, _ := loader.LoadInterface(context.Background(), dir, interfaceName)
	v, err := generator.GenerateVanillaMock(iface, interfaceName)
	if err != nil {
		t.Error(err)
		return
	}

	actual := v.Output()

	expected := `package mock

type VanillerVMock struct {
	CombinationFn func(i0 int64) (string, error)
	ExternalFn func(c difflib.UnifiedDiff) ()
	IntValueFn func() (int64)
	StringParamFn func(s0 string) ()
	VariadicFn func(abc string, more ...string) (string)
	WithNameFn func(abc int) ()
}

func (v VanillerVMock) Combination(i0 int64) (string, error) {
	return v.CombinationFn(i0)
}

func (v VanillerVMock) External(c difflib.UnifiedDiff) () {
	v.ExternalFn(c)
}

func (v VanillerVMock) IntValue() (int64) {
	return v.IntValueFn()
}

func (v VanillerVMock) StringParam(s0 string) () {
	v.StringParamFn(s0)
}

func (v VanillerVMock) Variadic(abc string, more ...string) (string) {
	return v.VariadicFn(abc, more...)
}

func (v VanillerVMock) WithName(abc int) () {
	v.WithNameFn(abc)
}
`

	assert.Equal(t, expected, actual)
}

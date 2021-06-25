# vanilla-mock

Generate vanilla style mocks say from

```
type Vaniller interface {
	IntValue() int64
	StringParam(string)
	WithName(abc int)
	Combination(int64) (string, error)
	Variadic(abc string, more ...string) string
}
```

to

```
type VanillerVMock struct {
	CombinationFn func(int64) (string, error)
	IntValueFn func() int64
	StringParamFn func(string)
	VariadicFn func(abc string, more ...string) string
	WithNameFn func(abc int)
}

func (v VanillerVMock) Combination(i0 int64) (string, error) {
	return v.CombinationFn(i0)
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
```

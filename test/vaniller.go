package test

// Vaniller is an interface used for testing
type Vaniller interface {
	IntValue() int64
	StringParam(string)
	WithName(abc int)
	Combination(int64) (string, error)
	Variadic(abc string, more ...string) string
}

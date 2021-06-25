package test

import "github.com/pmezard/go-difflib/difflib"

// Vaniller is an interface used for testing
type Vaniller interface {
	IntValue() int64
	StringParam(string)
	WithName(abc int)
	Combination(int64) (string, error)
	Variadic(abc string, more ...string) string
	External(c difflib.UnifiedDiff)
}

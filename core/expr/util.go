package expr

import (
	"github.com/xnacly/sophia/core/serror"
	"github.com/xnacly/sophia/core/token"
)

// fastpath for casting bool, reduces memory allocation by skipping allocation
func MustBool(in any, t *token.Token) bool {
	switch v := in.(type) {
	case bool:
		return v
	default:
		serror.Add(t, "Type error", "Expected value of type bool, got %s", token.TOKEN_NAME_MAP[t.Type])
		serror.Panic()
	}
	// technically unreachable
	return false
}

// fastpath for casting float64, reduces memory allocation by skipping allocation
func MustFloat(in any, t *token.Token) float64 {
	switch v := in.(type) {
	case float64:
		return v
	default:
		serror.Add(t, "Type error", "Expected value of type float, got %s", token.TOKEN_NAME_MAP[t.Type])
		serror.Panic()
	}
	// technically unreachable
	return 0
}

// attempts to cast `in` to `T`, returns `in` cast to `T` if successful. If
// cast fails, panics.
func MustT[T any](in any, t *token.Token) T {
	val, ok := in.(T)
	if !ok {
		var e T
		serror.Add(t, "Type error", "Expected value of type %T, got %T", e, in)
		serror.Panic()
	}
	return val
}

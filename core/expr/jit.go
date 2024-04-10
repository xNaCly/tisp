package expr

// TODO: I wish to just put this in a different package

var JIT *Jit

// Jit is the just in time compiler for the sophia language using the go plugin
// api
//
// Downsides: no windows support and requires the go compiler on the users system
type Jit struct{}

// Compile compiles
func (j *Jit) Compile(ast *Func) (func(any) any, error) {
	// TODO: go code gen
	// TODO: go compiler invocation
	// TODO: opening the go plugin
	// TODO: look up generated function
	return nil, nil
}

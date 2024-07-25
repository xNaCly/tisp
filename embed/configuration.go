package embed

import "github.com/xnacly/sophia/core/types"

type Configuration struct {
	// enables the just in time compiler, this should only be enabled if you have the go toolchain on the system
	EnableJit bool
	// tells the sophia runtime to link the go standard library modules
	EnableGoStd bool // TODO:
	// expose functions written in go into the sophia runtime
	Functions map[string]types.KnownFunctionInterface
	// enable debug logs
	Debug bool
}

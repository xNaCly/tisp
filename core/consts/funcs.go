package consts

type Return struct {
	HasValue bool
	Value    any
}

var RETURN = Return{}

var FUNC_TABLE = make(map[uint32]any, 64)

// JIT_CONSTANT determines
var JIT_CONSTANT uint64 = 1_000

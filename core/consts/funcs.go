package consts

type Return struct {
	HasValue bool
	Value    any
}

var FUNC_TABLE = make(map[uint32]any, 16)
var RETURN = Return{}

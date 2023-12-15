package expr

import (
	"sophia/core/consts"
	"sophia/core/serror"
	"sophia/core/token"
	"sophia/core/types"
)

// function definition
type For struct {
	Token    *token.Token
	Params   types.Node
	LoopOver types.Node
	Body     []types.Node
}

func (f *For) GetChildren() []types.Node {
	return f.Body
}

func (n *For) SetChildren(c []types.Node) {
	n.Body = c
}

func (f *For) GetToken() *token.Token {
	return f.Token
}

func (f *For) Eval() any {
	params := f.Params.(*Params).Children
	if len(params) < 1 {
		serror.Add(f.Token, "Not enough arguments", "Expected at least %d parameters for loop, got %d.", 1, len(params))
		serror.Panic()
	}
	element := castPanicIfNotType[*Ident](params[0], params[0].GetToken())
	oldValue, foundOldValue := consts.SYMBOL_TABLE[element.Key]

	v := f.LoopOver.Eval()
	switch v.(type) {
	case []interface{}:
		loopOver := castPanicIfNotType[[]interface{}](v, f.LoopOver.GetToken())

		for _, el := range loopOver {
			consts.SYMBOL_TABLE[element.Key] = el
			for _, stmt := range f.Body {
				stmt.Eval()
			}
		}
	case float64:
		con := v.(float64)
		for i := 0.0; i < con; i++ {
			consts.SYMBOL_TABLE[element.Key] = i
			for _, stmt := range f.Body {
				stmt.Eval()
			}
		}
	default:
		t := f.LoopOver.GetToken()
		serror.Add(t, "Invalid iterator", "expected container or upper bound for iteration, got: %T\n", v)
		serror.Panic()
	}

	if foundOldValue {
		consts.SYMBOL_TABLE[element.Key] = oldValue
	}
	return nil
}

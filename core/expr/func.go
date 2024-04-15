package expr

import (
	"github.com/xnacly/sophia/core/consts"
	"github.com/xnacly/sophia/core/token"
	"github.com/xnacly/sophia/core/types"
)

// function definition
type Func struct {
	Token             *token.Token
	Name              types.Node
	Params            *Array
	Body              []types.Node
	Jited             func(...any) any
	ArgumentDataTypes []string
	WasJitted         bool
}

func (f *Func) GetChildren() []types.Node {
	return f.Body
}

func (n *Func) SetChildren(c []types.Node) {
	n.Body = c
}

func (f *Func) GetToken() *token.Token {
	return f.Token
}

func (f *Func) Eval() any {
	ident := f.Name.(*Ident)
	consts.FUNC_TABLE[ident.Key] = f
	return nil
}

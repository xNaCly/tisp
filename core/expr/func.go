package expr

import (
	"sophia/core/consts"
	"sophia/core/token"
	"sophia/core/types"
)

// function definition
type Func struct {
	Token  *token.Token
	Name   types.Node
	Params types.Node
	Body   []types.Node
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

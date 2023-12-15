package expr

import (
	"sophia/core/token"
	"sophia/core/types"
)

type Root struct {
	Children []types.Node
}

func (r *Root) GetChildren() []types.Node {
	return r.Children
}

func (r *Root) SetChildren(c []types.Node) {
	r.Children = c
}

func (r *Root) GetToken() *token.Token {
	return nil
}

func (r *Root) Eval() any {
	return nil
}

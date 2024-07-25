package expr

import (
	"github.com/xnacly/sophia/core/token"
	"github.com/xnacly/sophia/core/types"
)

type Or struct {
	Token    *token.Token
	Children []types.Node
}

func (o *Or) GetChildren() []types.Node {
	return o.Children
}

func (n *Or) SetChildren(c []types.Node) {
	n.Children = c
}

func (o *Or) GetToken() *token.Token {
	return o.Token
}

func (o *Or) Eval() any {
	if len(o.Children) == 2 {
		f := o.Children[0]
		s := o.Children[1]
		return MustBool(f.Eval(), f.GetToken()) || MustBool(s.Eval(), s.GetToken())
	}
	for _, c := range o.Children {
		if MustBool(c.Eval(), c.GetToken()) {
			return true
		}
	}
	return false
}

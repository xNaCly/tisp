package expr

import (
	"github.com/xnacly/sophia/core/token"
	"github.com/xnacly/sophia/core/types"
)

type And struct {
	Token    *token.Token
	Children []types.Node
}

func (a *And) GetChildren() []types.Node {
	return a.Children
}

func (n *And) SetChildren(c []types.Node) {
	n.Children = c
}

func (a *And) GetToken() *token.Token {
	return a.Token
}

func (a *And) Eval() any {
	// fastpaths
	if len(a.Children) == 2 {
		f := a.Children[0]
		s := a.Children[1]
		return MustBool(f.Eval(), f.GetToken()) && MustBool(s.Eval(), s.GetToken())
	}

	for _, c := range a.Children {
		v := MustBool(c.Eval(), a.Token)
		if !v {
			return false
		}
	}
	return true
}

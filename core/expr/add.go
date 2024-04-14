package expr

import (
	"github.com/xnacly/sophia/core/token"
	"github.com/xnacly/sophia/core/types"
)

type Add struct {
	Token    *token.Token
	Children []types.Node
}

func (a *Add) GetChildren() []types.Node {
	return a.Children
}

func (n *Add) SetChildren(c []types.Node) {
	n.Children = c
}

func (a *Add) GetToken() *token.Token {
	return a.Token
}

func (a *Add) Eval() any {
	if len(a.Children) == 2 {
		// fastpath for two children
		f := a.Children[0]
		s := a.Children[1]
		return MustFloat(f.Eval(), f.GetToken()) + MustFloat(s.Eval(), s.GetToken())
	}

	res := 0.0
	for i, c := range a.Children {
		if i == 0 {
			res = MustFloat(c.Eval(), c.GetToken())
		} else {
			res += MustFloat(c.Eval(), c.GetToken())
		}
	}
	return res
}

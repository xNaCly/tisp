package expr

import (
	"github.com/xnacly/sophia/core/token"
	"github.com/xnacly/sophia/core/types"
)

type Mul struct {
	Token    *token.Token
	Children []types.Node
}

func (m *Mul) GetChildren() []types.Node {
	return m.Children
}

func (n *Mul) SetChildren(c []types.Node) {
	n.Children = c
}

func (m *Mul) GetToken() *token.Token {
	return m.Token
}

func (m *Mul) Eval() any {
	if len(m.Children) == 2 {
		// fastpath for two children
		f := m.Children[0]
		s := m.Children[1]
		return MustFloat(f.Eval(), f.GetToken()) * MustFloat(s.Eval(), s.GetToken())
	}

	res := 0.0
	for i, c := range m.Children {
		if i == 0 {
			res = MustFloat(c.Eval(), c.GetToken())
		} else {
			res *= MustFloat(c.Eval(), c.GetToken())
		}
	}
	return res
}

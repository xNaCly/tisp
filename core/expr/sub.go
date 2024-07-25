package expr

import (
	"github.com/xnacly/sophia/core/token"
	"github.com/xnacly/sophia/core/types"
)

type Sub struct {
	Token    *token.Token
	Children []types.Node
}

func (s *Sub) GetChildren() []types.Node {
	return s.Children
}

func (n *Sub) SetChildren(c []types.Node) {
	n.Children = c
}

func (s *Sub) GetToken() *token.Token {
	return s.Token
}

func (s *Sub) Eval() any {
	if len(s.Children) == 2 {
		// fastpath for two children
		f := s.Children[0]
		s := s.Children[1]
		return MustFloat(f.Eval(), f.GetToken()) - MustFloat(s.Eval(), s.GetToken())
	}

	res := 0.0
	for i, c := range s.Children {
		if i == 0 {
			res = MustFloat(c.Eval(), c.GetToken())
		} else {
			res -= MustFloat(c.Eval(), c.GetToken())
		}
	}
	return res
}

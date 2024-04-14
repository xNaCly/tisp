package expr

import (
	"github.com/xnacly/sophia/core/token"
	"github.com/xnacly/sophia/core/types"
)

type Lt struct {
	Token    *token.Token
	Children []types.Node
}

func (l *Lt) GetChildren() []types.Node {
	return l.Children
}

func (n *Lt) SetChildren(c []types.Node) {
	n.Children = c
}

func (l *Lt) GetToken() *token.Token {
	return l.Token
}

func (l *Lt) Eval() any {
	return MustFloat(l.Children[0].Eval(), l.Children[0].GetToken()) < MustFloat(l.Children[1].Eval(), l.Children[1].GetToken())
}

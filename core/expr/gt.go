package expr

import (
	"github.com/xnacly/sophia/core/token"
	"github.com/xnacly/sophia/core/types"
)

type Gt struct {
	Token    *token.Token
	Children []types.Node
}

func (g *Gt) GetChildren() []types.Node {
	return g.Children
}

func (n *Gt) SetChildren(c []types.Node) {
	n.Children = c
}

func (g *Gt) GetToken() *token.Token {
	return g.Token
}

func (g *Gt) Eval() any {
	return MustFloat(g.Children[0].Eval(), g.Children[0].GetToken()) > MustFloat(g.Children[1].Eval(), g.Children[1].GetToken())
}

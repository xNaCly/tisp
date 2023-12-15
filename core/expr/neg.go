package expr

import (
	"sophia/core/serror"
	"sophia/core/token"
	"sophia/core/types"
)

type Neg struct {
	Token    *token.Token
	Children types.Node
}

func (n *Neg) GetChildren() []types.Node {
	return []types.Node{n.Children}
}

func (n *Neg) SetChildren(c []types.Node) {
	if len(c) == 0 {
		return
	}
	n.Children = c[0]
}

func (n *Neg) GetToken() *token.Token {
	return n.Token
}

func (n *Neg) Eval() any {
	child := n.Children.Eval()
	var r any
	switch v := child.(type) {
	case float64:
		r = v * -1
	case bool:
		r = !v
	default:
		t := n.Children.GetToken()
		serror.Add(t, "Type Error", "Expected float64 or bool, got %T", child)
		serror.Panic()
	}
	return r
}

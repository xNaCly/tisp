package expr

import (
	"sophia/core/token"
	"strings"
)

type Or struct {
	Token    token.Token
	Children []Node
}

func (o *Or) GetToken() token.Token {
	return o.Token
}

func (o *Or) Eval() any {
	for _, c := range o.Children {
		if castPanicIfNotType[bool](c.Eval(), token.OR) {
			return true
		}
	}
	return false
}
func (n *Or) CompileJs(b *strings.Builder) {}

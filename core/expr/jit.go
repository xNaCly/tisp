package expr

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/xnacly/sophia/core/debug"
	"github.com/xnacly/sophia/core/types"
)

// TODO: I wish to just put this in a different package

var JIT *Jit

// Jit is the just in time compiler for the sophia language using the go plugin
// api
//
// Downsides: no windows support and requires the go compiler on the users system
type Jit struct {
	mutex sync.Mutex
}

// Compile compiles
func (j *Jit) Compile(ast *Func) (func(any) any, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	buf := &bytes.Buffer{}
	buf.WriteString("package main;func ")
	buf.WriteString(ast.Name.(*Ident).Name)
	buf.WriteRune('(')
	for i, arg := range ast.Params.Children {
		a := arg.(*Ident)
		buf.WriteString(a.Name)
		buf.WriteString(" any")
		if i+1 != len(ast.Params.Children) {
			buf.WriteRune(',')
		}
	}
	buf.WriteRune(')')
	buf.WriteRune('{')
	err := codeGen(buf, ast.Body)
	if err != nil {
		return nil, err
	}
	buf.WriteRune('}')

	debug.Log(buf.String())

	// TODO: go compiler invocation
	// TODO: opening the go plugin
	// TODO: look up generated function

	debug.Logf("Done compiling %q\n", ast.Name.GetToken().Raw)
	return nil, nil
}

func codeGen(b *bytes.Buffer, node []types.Node) error {
	for _, n := range node {
		switch t := n.(type) {
		case *String:
			b.WriteRune('"')
			b.WriteString(t.Token.Raw)
			b.WriteRune('"')
		case *Float:
			b.WriteString(t.Token.Raw)
		case *Ident:
			b.WriteString(t.Name)
		case *Var:
			if len(t.Value) > 1 {
				return fmt.Errorf("Variables containing more than 1 constant not supported by the JIT: %T", t)
			}
			b.WriteString(t.Ident.Name)
			b.WriteString(" := ")
			err := codeGen(b, t.Value)
			if err != nil {
				return err
			}
			b.WriteRune(';')
		default:
			return fmt.Errorf("Expression %T not yet supported by the JIT", t)
		}
	}
	return nil
}

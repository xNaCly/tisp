package expr

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"plugin"
	"strings"
	"sync"

	"github.com/xnacly/sophia/core"
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

	functionName := strings.ToTitle(ast.Name.(*Ident).Name)
	buf := &bytes.Buffer{}
	buf.WriteString("package main;func ")
	buf.WriteString(functionName)
	buf.WriteRune('(')
	for i, arg := range ast.Params.Children {
		a := arg.(*Ident)
		buf.WriteString(a.Name)
		buf.WriteString(" any")
		if i+1 != len(ast.Params.Children) {
			buf.WriteRune(',')
		}
	}
	buf.WriteString(")any{")
	err := codeGen(buf, ast.Body, true)
	if err != nil {
		return nil, err
	}
	buf.WriteRune('}')

	debug.Log("[JIT] generated code:", buf.String())

	name := "jit_" + functionName
	file, err := os.Create(name + ".go")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())

	_, err = buf.WriteTo(file)
	if err != nil {
		return nil, err
	}

	sName := name + ".so"
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", sName, file.Name())
	if core.CONF.Debug {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	defer os.Remove(sName)

	plug, err := plugin.Open(sName)
	if err != nil {
		return nil, err
	}

	symbol, err := plug.Lookup(functionName)
	if err != nil {
		return nil, err
	}

	function, ok := symbol.(func(any) any)
	if !ok {
		return nil, errors.New("Failed to cast the function to the given type")
	}

	debug.Logf("[JIT] Done compiling %q\n", ast.Name.GetToken().Raw)
	return function, nil
}

func codeGen(b *bytes.Buffer, node []types.Node, final bool) error {
	for i, n := range node {
		if final && i+1 == len(node) {
			b.WriteString("return ")
		}
		switch t := n.(type) {
		case *String:
			b.WriteRune('"')
			b.WriteString(t.Token.Raw)
			b.WriteRune('"')
		case *Float:
			b.WriteString(t.Token.Raw)
		case *Boolean:
			if t.Value {
				b.WriteString("true")
			} else {
				b.WriteString("false")
			}
		case *Ident:
			b.WriteString(t.Name)
		case *Var:
			if len(t.Value) > 1 {
				return fmt.Errorf("Variables containing more than 1 constant not supported by the JIT: %T", t)
			}
			b.WriteString(t.Ident.Name)
			b.WriteString(" := ")
			err := codeGen(b, t.Value, false)
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

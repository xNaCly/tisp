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
func (j *Jit) Compile(ast *Func) (func(...any) any, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	functionName := strings.ToTitle(ast.Name.(*Ident).Name)
	buf := &bytes.Buffer{}
	buf.WriteString(`package main;func `)
	buf.WriteString(functionName)
	buf.WriteString("(args ...any)any{")
	// regenerating params with type assertions for params
	for i, arg := range ast.Params.Children {
		a := arg.(*Ident)
		buf.WriteString(a.Name)
		buf.WriteString(" := ")
		buf.WriteString("args[")
		buf.WriteRune(rune(i + 48))
		buf.WriteString("].(")
		buf.WriteString(ast.ArgumentDataTypes[i])
		buf.WriteString(");")
	}

	err := codeGen(buf, true, ast.Body...)
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

	function, ok := symbol.(func(...any) any)
	if !ok {
		return nil, errors.New("Failed to cast the function to the given type")
	}

	debug.Logf("[JIT] Done compiling %q\n", ast.Name.GetToken().Raw)
	return function, nil
}

func genBinary(b *bytes.Buffer, node types.Node, op rune) error {
	c := node.GetChildren()
	if len(c) != 2 {
		return fmt.Errorf("%T with more than 2 children not supported, got %d", node, len(c))
	}

	err := codeGen(b, false, c[0])
	if err != nil {
		return err
	}
	b.WriteRune(op)
	return codeGen(b, false, c[1])
}

func codeGen(b *bytes.Buffer, final bool, node ...types.Node) error {
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
		case *Array:
			b.WriteString("[]any{")
			for i, c := range t.Children {
				codeGen(b, false, c)
				if i+1 != len(t.Children) {
					b.WriteRune(',')
				}
			}
			b.WriteRune('}')
		case *Var:
			if len(t.Value) > 1 {
				return fmt.Errorf("Variables containing more than 1 constant not supported by the JIT: %T", t)
			}
			b.WriteString(t.Ident.Name)
			b.WriteString(" := ")
			err := codeGen(b, false, t.Value...)
			if err != nil {
				return err
			}
			b.WriteRune(';')
		case *Add:
			err := genBinary(b, t, '+')
			if err != nil {
				return err
			}
		case *Sub:
			err := genBinary(b, t, '-')
			if err != nil {
				return err
			}
		case *Div:
			err := genBinary(b, t, '/')
			if err != nil {
				return err
			}
		case *Mul:
			err := genBinary(b, t, '*')
			if err != nil {
				return err
			}

			// BUG: disabled because requires "math" import
			// case *Mod:
			// 	if len(t.Children) != 2 {
			// 		return fmt.Errorf("%T with more than 2 children not supported, got %d", t, len(t.Children))
			// 	}
			// 	b.WriteString("math.Mod(")
			// 	err := codeGen(b, false, t.Children[0])
			// 	if err != nil {
			// 		return err
			// 	}
			// 	b.WriteRune(',')
			// 	err = codeGen(b, false, t.Children[1])
			// 	if err != nil {
			// 		return err
			// 	}
			// 	b.WriteString(")")
		case *Neg:
			b.WriteRune('!')
			err := codeGen(b, false, t.Children)
			if err != nil {
				return err
			}
		case *Gt:
			err := codeGen(b, false, t.Children[0])
			if err != nil {
				return err
			}
			b.WriteRune('>')
			err = codeGen(b, false, t.Children[1])
			if err != nil {
				return err
			}
		case *Lt:
			err := codeGen(b, false, t.Children[0])
			if err != nil {
				return err
			}
			b.WriteRune('<')
			err = codeGen(b, false, t.Children[1])
			if err != nil {
				return err
			}
		case *Equal:
			err := codeGen(b, false, t.Children[0])
			if err != nil {
				return err
			}
			b.WriteString("==")
			err = codeGen(b, false, t.Children[1])
			if err != nil {
				return err
			}
		case *And:
			err := codeGen(b, false, t.Children[0])
			if err != nil {
				return err
			}
			b.WriteString("&&")
			err = codeGen(b, false, t.Children[1])
			if err != nil {
				return err
			}
		case *Or:
			err := codeGen(b, false, t.Children[0])
			if err != nil {
				return err
			}
			b.WriteString("||")
			err = codeGen(b, false, t.Children[1])
			if err != nil {
				return err
			}
		case *If:
			b.WriteString("if ")
			err := codeGen(b, false, t.Condition)
			if err != nil {
				return err
			}
			b.WriteString(" {")
			err = codeGen(b, false, t.Body...)
			if err != nil {
				return err
			}
			b.WriteString("};")
		case *Return:
			b.WriteString("return ")
			err := codeGen(b, false, t.Child)
			if err != nil {
				return err
			}
		case *For:
		default:
			return fmt.Errorf("Expression %T not yet supported by the JIT", t)
		}
	}
	return nil
}

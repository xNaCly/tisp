package expr

import (
	"fmt"

	"github.com/xnacly/sophia/core/consts"
	"github.com/xnacly/sophia/core/debug"
	"github.com/xnacly/sophia/core/serror"
	"github.com/xnacly/sophia/core/token"
	"github.com/xnacly/sophia/core/types"
)

type Call struct {
	Token *token.Token
	Key   uint32
	Args  []types.Node
	Calls uint64
}

func (c *Call) GetChildren() []types.Node {
	return c.Args
}

func (n *Call) SetChildren(c []types.Node) {
	n.Args = c
}

func (c *Call) GetToken() *token.Token {
	return c.Token
}

func typeName(tn any) (r string, err error) {
	switch t := tn.(type) {
	case float64:
		r = "float64"
	case string:
		r = "string"
	case bool:
		r = "bool"
	default:
		err = fmt.Errorf("Unsupported type %T for argument, skipping compilation", t)
		r = ""
	}
	return
}

func (c *Call) Eval() any {
	storedFunc, ok := consts.FUNC_TABLE[c.Key]
	if !ok {
		serror.Add(c.Token, "Undefined function", "Function %q not defined", c.Token.Raw)
		serror.Panic()
	}

	def, ok := storedFunc.(*Func)
	// fastpath for built-in function support, see core/builtin/builtin.go
	if !ok {
		// this branch is hit if a function is not of type *Func which only
		// happens for built ins, thus the cast can not fail
		function, _ := storedFunc.(types.KnownFunctionInterface)
		return function(c.Token, c.Args...)
	}

	if !def.WasJitted && c.Calls >= consts.JIT_CONSTANT && def.Jited == nil && JIT != nil {
		def.WasJitted = true
		def.ArgumentDataTypes = make([]string, len(def.Params.Children))
		for i, a := range c.Args {
			if t, err := typeName(a.Eval()); err != nil {
				debug.Logf("[JIT] Skipping compilation of %q: %s\n", c.Token.Raw, err)
				def.ArgumentDataTypes = []string{}
				break
			} else {
				def.ArgumentDataTypes[i] = t
			}
		}
		if len(def.ArgumentDataTypes) != 0 {
			go func(functionDef *Func) {
				debug.Logf("[JIT] Attempting to compile function %q\n", c.Token.Raw)
				fun, err := JIT.Compile(def)
				if err != nil {
					debug.Logf("[JIT] Failed to compile function %q: %s, bailing out to the interpreter\n", c.Token.Raw, err)
					return
				}
				def.Jited = fun
			}(def)
		}
	} else if def.Jited != nil {
		args := make([]any, len(c.Args))
		for i, arg := range c.Args {
			args[i] = arg.Eval()
		}
		return def.Jited(args...)
	} else {
		c.Calls++
	}

	return callFunction(c.Token, def.Body, def.Params, c.Args)
}

func callFunction(tok *token.Token, body []types.Node, params *Array, args []types.Node) any {
	if len(params.Children) != len(args) {
		argLen := len(args)
		if len(params.Children) < argLen {
			serror.Add(tok, "Too many arguments", "Too many arguments for %q, wanted %d, got %d", tok.Raw, len(params.Children), len(args))
			serror.Panic()
		} else if len(params.Children) > argLen {
			serror.Add(tok, "Not enough arguments", "Not enough arguments for %q, wanted %d, got %d", tok.Raw, len(params.Children), len(args))
			serror.Panic()
		}
	}

	// store variable values from before entering the function scope
	for i, arg := range args {
		identifier := params.Children[i].(*Ident)
		if val, ok := consts.SYMBOL_TABLE[identifier.Key]; ok {
			consts.SCOPE_TABLE[identifier.Key] = val
		}
		consts.SYMBOL_TABLE[identifier.Key] = arg.Eval()
	}

	var ret any

	for i, stmt := range body {
		// enabling early returns
		if consts.RETURN.HasValue {
			ret = consts.RETURN.Value
			consts.RETURN.HasValue = false
			consts.RETURN.Value = nil
			break
		}
		if i+1 == len(body) {
			ret = stmt.Eval()
			break
		}
		stmt.Eval()
	}

	// if last line was a return
	if consts.RETURN.HasValue {
		ret = consts.RETURN.Value
		consts.RETURN.HasValue = false
		consts.RETURN.Value = nil
	}

	defer func() {
		// going out of scope, therefore we restore variables used in the
		// function scope to their previous value stored in the local scope table
		for k, v := range consts.SCOPE_TABLE {
			consts.SYMBOL_TABLE[k] = v
			delete(consts.SCOPE_TABLE, k)
		}
	}()

	return ret

}

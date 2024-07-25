package run

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/xnacly/sophia/core"
	_ "github.com/xnacly/sophia/core/builtin"
	"github.com/xnacly/sophia/core/debug"
	"github.com/xnacly/sophia/core/eval"
	"github.com/xnacly/sophia/core/lexer"
	"github.com/xnacly/sophia/core/parser"
	"github.com/xnacly/sophia/core/serror"
)

// runtime execution starting point
func Run(r io.Reader, filename string) (s []string, e error) {
	defer func() {
		if core.CONF.Debug {
			return
		}
		if err := recover(); err != nil {
			serror.Display()
			if err, ok := err.(error); ok {
				// catch all for panics
				if !strings.Contains(err.Error(), "sophia: ") {
					fmt.Println(err)
				}
			}
			return
		}
	}()

	debug.Log("starting lexer")
	l := lexer.New(r)
	tokens := l.Lex()
	if serror.HasErrors() {
		serror.Display()
		e = errors.New("Syntax errors found, skipping remaining interpreter stages. (parsing and evaluation)")
		return
	}
	debug.Log("lexed", len(tokens), "token")

	debug.Log(debug.Token(tokens))

	debug.Log("starting parser")
	p := parser.New(tokens, filename)
	ast := p.Parse()
	if serror.HasErrors() {
		serror.Display()
		e = errors.New("Semantic errors found, skipping remaining interpreter stages. (evaluation)")
		return
	}

	if core.CONF.Debug {
		out, _ := json.MarshalIndent(ast, "", "  ")
		debug.Log("ast:", string(out))
	}

	debug.Log("done parsing - starting eval")

	if len(ast) == 0 {
		return
	}

	start := time.Now()
	s = eval.Eval(filename, ast)
	fmt.Println("Took: ", time.Since(start))
	debug.Log("done evaling")

	return
}

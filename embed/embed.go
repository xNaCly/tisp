// embed provides helpers and abstractions over the github.com/xnacly/sophia runtime for easily
// embedding the runtime and registering functions via the known function
// interface.
package embed

import (
	"bytes"
	"io"
	"os"

	"github.com/xnacly/sophia/core"
	"github.com/xnacly/sophia/core/alloc"
	"github.com/xnacly/sophia/core/consts"
	"github.com/xnacly/sophia/core/expr"
	"github.com/xnacly/sophia/core/run"
	"github.com/xnacly/sophia/core/serror"
)

// calls into the runtime to apply the given configuration
func Embed(config Configuration) {
	if config.EnableGoStd {
		panic("Embedding error: Linking the go standard library via modules is currently not implemented")
	}

	if config.EnableJit {
		expr.JIT = &expr.Jit{}
	}

	core.CONF.Debug = config.Debug

	for name, function := range config.Functions {
		consts.FUNC_TABLE[alloc.NewFunc(name)] = function
	}
}

// starts the runtime, returns error op on error occurrence, writes prints and
// errors to w, is w nil, os.Stdout is used
func Execute(file *os.File, w io.Writer) error {
	buf := &bytes.Buffer{}
	r := io.TeeReader(file, buf)
	buf.ReadFrom(r)
	serror.SetDefault(serror.NewFormatter(&core.CONF, buf.String(), file.Name(), nil))
	_, err := run.Run(buf, file.Name())
	return err
}

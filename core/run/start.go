package run

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/xnacly/sophia/core"
	"github.com/xnacly/sophia/core/debug"
	"github.com/xnacly/sophia/core/serror"
)

// entry point for sophia cli
func Start() {
	execute := flag.String("exp", "", "specifiy expression to execute")
	dbg := flag.Bool("dbg", false, "enable debug logs")
	allErrors := flag.Bool("all-errors", false, "display all found errors")
	jit := flag.Bool("jit", true, "disabled or enable jit compiler")
	flag.Parse()
	core.CONF = core.Config{
		Debug:     *dbg,
		AllErrors: *allErrors,
		JIT:       *jit,
	}

	if *dbg {
		log.SetFlags(log.Ltime | log.Lmicroseconds)
	} else {
		log.SetFlags(0)
	}

	stdinInf, err := os.Stdin.Stat()
	// check if stdin is readable and the process is in a pipe
	if err == nil && !(stdinInf.Mode()&os.ModeNamedPipe == 0) {
		debug.Log("got stdin content, running...")
		buf := bytes.Buffer{}
		buf.ReadFrom(os.Stdin)
		serror.SetDefault(serror.NewFormatter(&core.CONF, buf.String(), "stdin", nil))
		_, err = Run(bytes.NewReader(buf.Bytes()), "stdin")
		if err != nil {
			log.Fatalln(err)
		}
	} else if len(*execute) != 0 {
		debug.Log("got -exp flag, running...")
		serror.SetDefault(serror.NewFormatter(&core.CONF, *execute, "cli", nil))
		_, err := Run(strings.NewReader(*execute), "cli")
		if err != nil {
			log.Fatalln(err)
		}
	} else if len(flag.Args()) == 1 {
		debug.Log("got file, running...")
		file := flag.Args()[0]
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("Failed to open file: %s\n", err)
		}
		defer f.Close()
		buf := &bytes.Buffer{}
		r := io.TeeReader(f, buf)
		buf.ReadFrom(r)
		serror.SetDefault(serror.NewFormatter(&core.CONF, buf.String(), file, nil))
		_, err = Run(buf, file)
		if err != nil {
			log.Fatalln("\n" + err.Error())
		}
	} else {
		fmt.Print(core.ASCII_ART, "\n")
		debug.Log("got nothing, starting repl...")
		repl(Run)
	}
}

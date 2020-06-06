package runner

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/ripta/mwnci/pkg/evaluator"
	"github.com/ripta/mwnci/pkg/lexer"
	"github.com/ripta/mwnci/pkg/object"
	"github.com/ripta/mwnci/pkg/parser"
)

var t0 = time.Now()

func REPL(out io.Writer, in io.Reader) error {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, "» ")
		if !scanner.Scan() {
			break
		}

		p := parser.New(lexer.New(scanner.Text()))

		prog := p.ParseProgram()
		if len(p.Errors()) != 0 {
			dumpErrors(out, p.Errors())
			continue
		}

		if res := evaluator.Eval(prog, env); res != nil {
			log(out, "%s\n", res.Inspect())
		}
	}

	log(out, "👋🏽\n")
	return nil
}

func log(out io.Writer, format string, a ...interface{}) {
	t := time.Since(t0).Truncate(time.Millisecond)
	a = append([]interface{}{t}, a...)
	_, _ = fmt.Fprintf(out, "(%s) "+format, a...)
}

func dumpErrors(out io.Writer, errors []string) {
	log(out, "Error:\n")
	for _, msg := range errors {
		log(out, "\t%s\n", msg)
	}
}

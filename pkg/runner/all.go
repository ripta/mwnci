package runner

import (
	"errors"
	"io"

	"github.com/ripta/mwnci/pkg/evaluator"
	"github.com/ripta/mwnci/pkg/lexer"
	"github.com/ripta/mwnci/pkg/object"
	"github.com/ripta/mwnci/pkg/parser"
)

func All(out io.Writer, in io.Reader) error {
	bs, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	env := object.NewEnvironment()
	p := parser.New(lexer.New(string(bs)))

	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		dumpErrors(out, p.Errors())
		return nil
	}

	res := evaluator.Eval(prog, env)
	if res.Type() == object.ErrorObj {
		return errors.New(res.Inspect())
	}

	return nil
}

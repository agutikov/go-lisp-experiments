package lispy

import (
	"io/ioutil"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
	"github.com/agutikov/go-lisp-experiments/lispy/syntax/lexer"
	"github.com/agutikov/go-lisp-experiments/lispy/syntax/parser"
	"github.com/agutikov/go-lisp-experiments/lispy/syntax/token"
)

func pasrse_bytes(bytes []byte) ast.Sequence {
	p := parser.NewParser()
	lex := lexer.NewLexer(bytes)

	if lex.Scan().Type == token.INVALID {
		// no valid tokens - return empty sequence
		return ast.Sequence{}
	}
	lex.Reset()

	st, err := p.Parse(lex)
	if err != nil {
		panic(err)
	}

	seq, ok := st.(ast.Sequence)
	if !ok {
		panic("Invalid parser output type")
	}

	return seq
}

func ParseStr(s string) ast.Sequence {
	return pasrse_bytes([]byte(s))
}

func ParseFile(filename string) ast.Sequence {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return pasrse_bytes(bytes)
}

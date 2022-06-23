package parser

import (
	"reflect"
	"testing"

	"github.com/agutikov/go-lisp-experiments/lispy/parser/ast"
	"github.com/agutikov/go-lisp-experiments/lispy/parser/lexer"
	"github.com/agutikov/go-lisp-experiments/lispy/parser/parser"
)

func Test_Sexpr(t *testing.T) {
	input := []byte("hello #gocc-lispy_1")
	lex := lexer.NewLexer(input)
	p := parser.NewParser()
	st, err := p.Parse(lex)
	if err != nil {
		panic(err)
	}
	w, ok := st.(*ast.List)
	if !ok {
		t.Fatalf("This is not a Sexpr")
	}
	expected := []ast.Atom{ast.Atom{Value: "hello"}, ast.Atom{Value: "#gocc-lispy_1"}}
	if !reflect.DeepEqual(w.List, expected) {
		t.Fatalf("Wrong Sexpr %v", w)
	}
}

package parser

import (
	"reflect"
	"testing"

	"github.com/agutikov/go-lisp-experiments/lispy/parser/ast"
	"github.com/agutikov/go-lisp-experiments/lispy/parser/lexer"
	"github.com/agutikov/go-lisp-experiments/lispy/parser/parser"
)

type ParserTestItem struct {
	input  string
	output ast.Sexpr
}

func Test_Sexpr(t *testing.T) {
	tests := []ParserTestItem{
		{"str", ast.Sexpr{ast.Atom{"str"}}},
		{"()", ast.Sexpr{ast.List{}}},
		{"( ( ) ( ) )", ast.Sexpr{ast.List{[]ast.Sexpr{ast.Sexpr{ast.List{}}, ast.Sexpr{ast.List{}}}}}},
		{"(a b)", ast.Sexpr{ast.List{[]ast.Sexpr{ast.Sexpr{ast.Atom{"a"}}, ast.Sexpr{ast.Atom{"b"}}}}}},
	}

	p := parser.NewParser()

	for _, test := range tests {
		lex := lexer.NewLexer([]byte(test.input))
		st, err := p.Parse(lex)
		if err != nil {
			panic(err)
		}
		s, ok := st.(*ast.Sexpr)
		t.Logf("%q, %+v -> %+v", test.input, test.output, *s)
		if !ok {
			t.Fatalf("This is not a Sexpr")
		}
		if !reflect.DeepEqual(*s, test.output) {
			t.Fatalf("Wrong Sexpr %v", *s)
		}
	}
}

package parser

import (
	"reflect"
	"testing"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
	"github.com/agutikov/go-lisp-experiments/lispy/syntax/lexer"
	"github.com/agutikov/go-lisp-experiments/lispy/syntax/parser"
)

type ParserTestItem struct {
	input  string
	output ast.Sexpr
}

func Test_Sexpr(t *testing.T) {

	tests := []ParserTestItem{
		{"str", ast.Symbol{"str"}},
		{"()", ast.List{}},
		{"( ( ) ( ) )", ast.List{[]ast.Sexpr{ast.List{}, ast.List{}}}},
		{"(a b)", ast.List{[]ast.Sexpr{ast.Symbol{"a"}, ast.Symbol{"b"}}}},
		{"0", ast.Number{0}},
		{"(+ 99 -1000)", ast.List{[]ast.Sexpr{ast.Symbol{"+"}, ast.Number{99}, ast.Number{-1000}}}},

		{"\"\"", ast.Str{""}},
		{"\"a\"", ast.Str{"a"}},
		{"\" \"", ast.Str{" "}},
		{"\" a \"", ast.Str{" a "}},
		{"\"\\\\\"", ast.Str{"\\"}},
		{"\"\\n\"", ast.Str{"\n"}},
		{"\"\\\"\"", ast.Str{"\""}},
		{"\" \\\" X \\\" \"", ast.Str{" \" X \" "}},

		{"(\"\" \" \" \"string literal\")",
			ast.List{[]ast.Sexpr{ast.Str{""}, ast.Str{" "}, ast.Str{"string literal"}}}},
	}

	p := parser.NewParser()

	for _, test := range tests {
		t.Logf("%q, expected: %#v", test.input, test.output)

		lex := lexer.NewLexer([]byte(test.input))
		st, err := p.Parse(lex)
		if err != nil {
			panic(err)
		}

		s, ok := st.(ast.Sexpr)

		t.Logf("    -> %#v", s)

		if !ok {
			t.Fatalf("This is not a Sexpr")
		}
		if !reflect.DeepEqual(s, test.output) {
			t.Fatalf("Wrong Sexpr %#v", s)
		}
	}
}

package parser

import (
	"reflect"
	"testing"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
	"github.com/agutikov/go-lisp-experiments/lispy/syntax/lexer"
	"github.com/agutikov/go-lisp-experiments/lispy/syntax/parser"
)

type ParserTestSexpr struct {
	input  string
	output ast.Sexpr
}

func Test_Sexpr(t *testing.T) {

	tests := []ParserTestSexpr{
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

		{"nil", ast.Nil{}},
		{"(nil nil)", ast.List{[]ast.Sexpr{ast.Nil{}, ast.Nil{}}}},

		{"'()", ast.Quote{ast.List{}}},
		{",()", ast.Unquote{ast.List{}}},
		{"'(x ,y)", ast.Quote{ast.List{[]ast.Sexpr{
			ast.Symbol{"x"}, ast.Unquote{ast.Symbol{"y"}},
		}}}},
		{"'(,('()))", ast.Quote{ast.List{[]ast.Sexpr{
			ast.Unquote{ast.List{[]ast.Sexpr{ast.Quote{ast.List{}}}}},
		}}}},

		{"(lambda (x) (- x))", ast.Lambda{
			Args: ast.List{[]ast.Sexpr{ast.Symbol{"x"}}},
			Body: ast.List{[]ast.Sexpr{ast.Symbol{"-"}, ast.Symbol{"x"}}},
		}},

		{"(if t f ())", ast.If{
			Test:      ast.Bool(true),
			PosBranch: ast.Bool(false),
			NegBranch: ast.List{},
		}},

		{"(define foo 10)", ast.Define{
			Sym:   ast.Symbol{"foo"},
			Value: ast.Number{10},
		}},

		{"(set! foo 10)", ast.Set{
			Sym:   ast.Symbol{"foo"},
			Value: ast.Number{10},
		}},
	}

	p := parser.NewParser()

	for _, test := range tests {
		expected := ast.Sequence{[]ast.Sexpr{test.output}}
		t.Logf("%q, expected: %#v", test.input, expected)

		lex := lexer.NewLexer([]byte(test.input))
		st, err := p.Parse(lex)
		if err != nil {
			panic(err)
		}

		s, ok := st.(ast.Sequence)

		t.Logf("    -> %#v", s)

		if !ok {
			t.Fatalf("This is not a Sequence")
		}
		if !reflect.DeepEqual(s, expected) {
			t.Fatalf("Wrong Sexpr:\n%#v\nExpected:\n%#v\nInput: %q", s, expected, test.input)
		}
	}
}

type ParserTestSeq struct {
	input  string
	output ast.Sequence
}

func Test_Sequence(t *testing.T) {
	tests := []ParserTestSeq{
		{"a", ast.Sequence{[]ast.Sexpr{ast.Symbol{"a"}}}},
		{"() ; comment\n", ast.Sequence{[]ast.Sexpr{ast.List{}}}},
		{";; line 1\n () ; line 2\n", ast.Sequence{[]ast.Sexpr{ast.List{}}}},
		{"() ; ()\n", ast.Sequence{[]ast.Sexpr{ast.List{}}}},

		{"\";\"", ast.Sequence{[]ast.Sexpr{ast.Str{";"}}}},
	}

	p := parser.NewParser()

	for _, test := range tests {
		expected := test.output
		t.Logf("%q, expected: %#v", test.input, expected)

		lex := lexer.NewLexer([]byte(test.input))
		st, err := p.Parse(lex)
		if err != nil {
			panic(err)
		}

		s, ok := st.(ast.Sequence)

		t.Logf("    -> %#v", s)

		if !ok {
			t.Fatalf("This is not a Sequence")
		}
		if !reflect.DeepEqual(s, expected) {
			t.Fatalf("Wrong Sexpr %#v", s)
		}
	}
}

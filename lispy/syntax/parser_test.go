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
	output ast.Any
}

func Test_Sexpr(t *testing.T) {

	tests := []ParserTestSexpr{
		{"str", ast.Symbol{"str"}},
		{"()", ast.List{}},
		{"( ( ) ( ) )", ast.List{ast.List{}, ast.List{}}},
		{"(a b)", ast.List{ast.Symbol{"a"}, ast.Symbol{"b"}}},
		{"0", ast.IntNum(0)},
		{"(+ 99 -1000)", ast.List{
			ast.Symbol{"+"}, ast.IntNum(99), ast.IntNum(-1000),
		}},

		{"\"\"", ast.Str{""}},
		{"\"a\"", ast.Str{"a"}},
		{"\" \"", ast.Str{" "}},
		{"\" a \"", ast.Str{" a "}},
		{"\"\\\\\"", ast.Str{"\\"}},
		{"\"\\n\"", ast.Str{"\n"}},
		{"\"\\\"\"", ast.Str{"\""}},
		{"\" \\\" X \\\" \"", ast.Str{" \" X \" "}},

		{"(\"\" \" \" \"string literal\")",
			ast.List{ast.Str{""}, ast.Str{" "}, ast.Str{"string literal"}},
		},

		{"nil", ast.Nil{}},
		{"(nil nil)", ast.List{ast.Nil{}, ast.Nil{}}},

		{"'()", ast.Quote{ast.List{}}},
		{",()", ast.Unquote{ast.List{}}},
		{"'(x ,y)", ast.Quote{ast.List{
			ast.Symbol{"x"}, ast.Unquote{ast.Symbol{"y"}},
		}}},
		{"'(,('()))", ast.Quote{ast.List{
			ast.Unquote{ast.List{ast.Quote{ast.List{}}}}},
		}},

		{"(lambda (x) (- x))", ast.Lambda{
			Args: []ast.Symbol{ast.Symbol{"x"}},
			Body: ast.List{ast.Symbol{"-"}, ast.Symbol{"x"}},
		}},

		{"(if t false ())", ast.If{
			Test:      ast.Bool(true),
			PosBranch: ast.Bool(false),
			NegBranch: ast.List{},
		}},

		{"(define foo 10)", ast.Define{
			Sym:   ast.Symbol{"foo"},
			Value: ast.IntNum(10),
		}},

		{"(set! foo 10)", ast.Set{
			Sym:   ast.Symbol{"foo"},
			Value: ast.IntNum(10),
		}},

		{"(defun foo (arg) ())", ast.Defun{
			Sym: ast.Symbol{"foo"},
			L: ast.Lambda{
				Args: []ast.Symbol{ast.Symbol{"arg"}},
				Body: ast.List{},
			},
		}},
	}

	p := parser.NewParser()

	for _, test := range tests {
		expected := ast.Sequence{test.output}
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
		{"a", ast.Sequence{ast.Symbol{"a"}}},
		{"() ; comment\n", ast.Sequence{ast.List{}}},
		{";; line 1\n () ; line 2\n", ast.Sequence{ast.List{}}},
		{"() ; ()\n", ast.Sequence{ast.List{}}},

		{"\";\"", ast.Sequence{ast.Str{";"}}},
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

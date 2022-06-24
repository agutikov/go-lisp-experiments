package lispy

import (
	"reflect"
	"testing"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
)

func Test_ParserDirect(t *testing.T) {
	lines := []string{
		"nil",
		"()",
		"0",
		"t",
		"false",
		"(nil 1 () (1 2) t)",
		"(define foo (+ (* x x) (* y y)))",
	}
	for _, line := range lines {
		t.Logf("%q", line)
		lst := ParseStr(line)
		s := LispyStr(lst)
		if s != line {
			t.Errorf("String representations are not equal: %q and %q", line, s)
		}
	}
}

func Test_ParserReverse(t *testing.T) {
	exprs := []Any{
		ast.Nil{},
		List{},
		ast.IntNum(0),
		Bool(true),
		Bool(false),
		List{
			ast.Nil{}, ast.IntNum(1), List{},
			List{ast.IntNum(1), Bool(false)},
			Bool(true),
		},
		List{
			Symbol{"_if"}, Symbol{"x"}, ast.IntNum(1), ast.IntNum(0),
		},
	}
	for _, expr := range exprs {
		e := ast.Sequence{expr}
		t.Logf("%#v", e)
		s := LispyStr(e)
		lst := ParseStr(s)
		if reflect.TypeOf(lst) != reflect.TypeOf(e) {
			t.Errorf("Types are not equal: %#v and %#v", e, lst)
		}
		if !reflect.DeepEqual(lst, e) {
			t.Errorf("Values are not equal: %#v and %#v", e, lst)
		}
	}
}

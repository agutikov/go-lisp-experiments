package lispy

import (
	"reflect"
	"testing"
)

func Test_ParserDirect(t *testing.T) {
	lines := []string{
		"nil",
		"()",
		"0",
		"t",
		"f",
		"(nil 1 () (1 2) t)",
		"(define foo (x y) (+ (* x x) (* y y)))",
	}
	for _, line := range lines {
		lst := ParseExpr(line)
		s := LispyStr(lst)
		if s != line {
			t.Errorf("String representations are not equal: %q and %q", line, s)
		}
	}
}

func Test_ParserReverse(t *testing.T) {
	exprs := []Any{
		nil,
		List{},
		Int(0),
		Bool(true),
		Bool(false),
		List{nil, Int(1), List{}, List{Int(1), Bool(false)}, Bool(true)},
		List{Builtin("if"), Symbol("x"), Int(1), Int(0)},
	}
	for _, expr := range exprs {
		s := LispyStr(expr)
		lst := ParseExpr(s)
		if reflect.TypeOf(lst) != reflect.TypeOf(expr) {
			t.Errorf("Types are not equal: %v and %v", expr, lst)
		}
		if !reflect.DeepEqual(lst, expr) {
			t.Errorf("Values are not equal: %v and %v", expr, lst)
		}
	}
}

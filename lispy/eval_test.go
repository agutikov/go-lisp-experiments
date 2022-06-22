package lispy

import (
	"reflect"
	"testing"
)

func Test_Eval(t *testing.T) {
	examples := [][]Any{
		{nil, nil},
		{List{}, List{}},
		{List{Builtin("quote"), List{Symbol("xxx"), Int(1)}}, List{Symbol("xxx"), Int(1)}},
	}
	for _, test := range examples {
		expected := test[1]
		expr := test[0]
		e := StdEnv()
		result := e.Eval(expr)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Not expected Eval() result: %v -> %v; expected: %v", expr, result, expected)
		}
	}
}

func Test_define(t *testing.T) {
	expr := "(define foo (lambda (x) (* x x)))"
	lst := NewParser(expr).ParseList()
	e := StdEnv()
	e.Eval(lst)
	v, ok := e.named_objects["foo"]
	if !ok {
		t.Errorf("define fails to update env")
	}
	f := func(args ...Any) Any { return nil }
	if reflect.TypeOf(v) != reflect.TypeOf(f) {
		t.Errorf("Invalid env object type")
	}
}

func Test_set(t *testing.T) {
	e := StdEnv()

	s1 := `(begin
		(define minus (lambda (x) (- 0 x)))
		(define bar (lambda (x) (* (foo x) 2)))
		(define foo (lambda (x) (if (> x 0) x 0)))
		(bar -1)
	)`
	expr1 := NewParser(s1).ParseList()

	r1 := e.Eval(expr1)
	if LispyStr(r1) != "0" {
		t.Errorf("Unexpected r1: %v", r1)
	}

	// TODO: set! replaces the foo in the env and bar behavior changes

	s2 := `(begin
		(set! foo (lambda (x) (minus x)))
		(bar -1)
	)`
	expr2 := NewParser(s2).ParseList()

	r2 := e.Eval(expr2)
	if LispyStr(r2) != "2" {
		t.Errorf("Unexpected r2: %v", r2)
	}
}

func Test_EvalStr(t *testing.T) {
	examples := [][]string{
		{"nil", "nil"},
		{"()", "()"},
		{"(quote (x 2 3))", "(x 2 3)"},
		{"(list 1 t nil ())", "(1 t nil ())"},
		{"(begin (define r 10) (* 3.14159265 (* r r)))", "314.159265"},
		{"(cons 1 ())", "(1)"},
		{"(cons 1 nil)", "(1)"},
		{"(cons 3 (cons 2 (cons 1 nil)))", "(3 2 1)"},
		{"(define x (list 1 2 3 4))", "nil"},
		{"(car x)", "1"},
		{"(cdr x)", "(2 3 4)"},
	}
	e := StdEnv()
	for _, test := range examples {
		expr := test[0]
		expected := test[1]
		lst := NewParser(expr).ParseList()
		result := e.Eval(lst)
		res_str := LispyStr(result)
		if res_str != expected {
			t.Errorf("Not expected Eval() result: %q -> %q, expected: %q", expr, res_str, expected)
		}
	}
}

func Test_lambda(t *testing.T) {
	f1 := Lambda("(lambda (x y) (/ (* x x) (* y y)))")
	r1 := f1(Int(2), Float(4))
	if r1 != Float(0.25) {
		t.Errorf("Unexpected r1: %v", r1)
	}
}

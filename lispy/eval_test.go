package lispy

import (
	"reflect"
	"testing"
)

func Test_Eval(t *testing.T) {
	examples := [][]Any{
		{nil, nil},
		{List{}, List{}},
		{List{Builtin("quote"), List{Symbol("xxx"), FromInt(1)}}, List{Symbol("xxx"), FromInt(1)}},
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
	lst := ParseExpr(expr)
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
		(define bar (lambda (x) (* (foo x) 2)))
		(define foo (lambda (x) (if (> x 0) x 0)))
		(bar -1)
	)`
	expr1 := ParseExpr(s1)

	r1 := e.Eval(expr1)
	if LispyStr(r1) != "0" {
		t.Errorf("Unexpected r1: %v", r1)
	}

	// TODO: set! replaces the foo in the env and bar behavior changes

	s2 := `(begin
		(set! foo (lambda (x) (- x)))
		(bar -1)
	)`
	expr2 := ParseExpr(s2)

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
		{"(define x (list 1 2 3 4))", "(1 2 3 4)"},
		{"(car x)", "1"},
		{"(cdr x)", "(2 3 4)"},
		{"(map - x)", "(-1 -2 -3 -4)"},
		{"(map cons x ())", "((1) (2) (3) (4))"},
		{"(map * (list 1 2) (list 10 20) (list -1 1))", "(-10 40)"},
		{"(apply + x)", "10"},
		{"(apply + 0 1 (list 2 3) 4)", "10"},
		{"(or nil 0 () t)", "t"},
		{"(cons nil nil)", "(nil)"},
		{"(and t 1 (cons nil nil) f)", "f"},
	}
	e := StdEnv()
	for _, test := range examples {
		expr := test[0]
		expected := test[1]
		lst := ParseExpr(expr)
		result := e.Eval(lst)
		res_str := LispyStr(result)
		if res_str != expected {
			t.Errorf("Not expected Eval() result: %q -> %q, expected: %q", expr, res_str, expected)
		}
	}
}

func Test_lambda(t *testing.T) {
	f1 := Lambda("(lambda (x y) (/ (* x x) (* y y)))")
	r1 := f1(FromInt(2), FromFloat(4))
	if !equal(r1, FromFloat(0.25)) {
		t.Errorf("Unexpected r1: %q", LispyStr(r1))
	}

	// from README.md
	fact := Lambda("(define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1))))))")
	r3 := fact(FromInt(100))
	expected3 := "93326215443944152681699238856266700490715968264381621468592963895217599993229915608941463976156518286253697920827223758251185210916864000000000000000000000000"
	if LispyStr(r3) != expected3 {
		t.Errorf("Unexpected r3: %q", LispyStr(r3))
	}

	zip2 := Lambda("(lambda (slice_1 slice_2) (map list slice_1 slice_2))")
	a := List{0, 1, 2}
	b := List{"str", true}
	r2 := zip2(a, b)
	expected2 := List{List{0, "str"}, List{1, true}, List{2, nil}}
	if !reflect.DeepEqual(r2, expected2) {
		t.Errorf("Unexpected r2: %q", LispyStr(r2))
	}
}

func Benchmark_Lambda(b *testing.B) {
	fact := Lambda("(define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1))))))")
	for i := 0; i < b.N; i++ {
		fact(FromInt(100))
	}
}

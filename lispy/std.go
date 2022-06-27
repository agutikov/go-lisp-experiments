package lispy

import (
	"math"
	"math/big"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
)

func car(arg ...Any) Any {
	return to_list(arg[0])[0]
}

func cdr(args ...Any) Any {
	return to_list(args[0])[1:]
}

func cons(args ...Any) Any {
	if len(args) != 2 {
		panic("'cons' function requires exactly 2 arguments (const head exp), while provided: " + LispyStr(args))
	}
	var r List
	tail := List{}
	if !is_nil(args[1]) {
		tail = to_list(args[1])
	}
	r = append(List{args[0]}, tail...)
	return r
}

func list(args ...Any) Any {
	return append(List{}, args...)
}

func fold_bools(f func(bool, bool) bool, init bool, args ...Any) Bool {
	acc := init

	for _, item := range args {
		acc = f(acc, if_test(item))
	}

	return Bool(acc)
}

func bool_to_int(b Bool) Int {
	if bool(b) {
		return ast.IntNum(1)
	} else {
		return ast.IntNum(0)
	}
}

func fold_nums(name string, f_i func(Int, Int) Int, f_f func(Float, Float) Float, init Any, args ...Any) Any {
	var acc Any
	acc = init
	//TODO: type switch for multiple args
	for _, item := range args {
		switch a := acc.(type) {
		case Int:
			switch v := item.(type) {
			case Bool:
				acc = f_i(a, bool_to_int(v))
			case Int:
				acc = f_i(a, v)
			case Float:
				acc = f_f(int_to_float(a), v)
			default:
				panic("Invalid '" + name + "' argument: " + LispyStr(args))
			}
		case Float:
			switch v := item.(type) {
			case Bool:
				acc = f_f(a, int_to_float(bool_to_int(v)))
			case Int:
				acc = f_f(a, int_to_float(v))
			case Float:
				acc = f_f(a, v)
			default:
				panic("Invalid '" + name + "' argument: " + LispyStr(args))
			}
		default:
			panic("'" + name + "' Error")
		}
	}

	return acc
}

func numeric_2_args(name string, f_i func(Int, Int) Int, f_f func(Float, Float) Float, args ...Any) Any {
	if len(args) != 2 {
		panic("'" + name + "' requires exactly 2 arguments, provided: " + LispyStr(args))
	}

	lhs := args[0]
	rhs := args[1]

	switch x := lhs.(type) {
	case Int:
		switch y := rhs.(type) {
		case Int:
			return f_i(x, y)
		case Float:
			return f_f(int_to_float(x), y)
		default:
			panic("Invalid '" + name + "' argument: " + LispyStr(args))
		}
	case Float:
		switch y := rhs.(type) {
		case Int:
			return f_f(x, int_to_float(y))
		case Float:
			return f_f(x, y)
		default:
			panic("Invalid '" + name + "' argument: " + LispyStr(args))
		}
	default:
		panic("'" + name + "' Error")
	}
}

//TODO: Number type that holds the big.Int or big.Rat or anything else and converts it if necessary

func sum(args ...Any) Any {
	return fold_nums("+",
		func(a Int, b Int) Int {
			r := Int{big.NewInt(0)}
			r.Value = r.Value.Add(a.Value, b.Value)
			return r
		},
		func(a Float, b Float) Float {
			r := ast.FloatNum(0)
			r.Value = r.Value.Add(a.Value, b.Value)
			return r
		},
		Int{big.NewInt(0)}, args...,
	)
}

func minus(arg Any) Any {
	switch x := arg.(type) {
	case Int:
		z := big.NewInt(0)
		return Int{z.Neg(x.Value)}
	case Float:
		z := big.NewRat(0, 1)
		return Float{z.Neg(x.Value)}
	default:
		panic("Invalid unary '-' argument: " + LispyStr(arg))
	}
}

func sub(args ...Any) Any {
	if len(args) == 1 {
		return minus(args[0])
	}
	return numeric_2_args("-",
		func(a Int, b Int) Int {
			r := Int{big.NewInt(0)}
			r.Value = r.Value.Sub(a.Value, b.Value)
			return r
		},
		func(a Float, b Float) Float {
			r := ast.FloatNum(0)
			r.Value = r.Value.Sub(a.Value, b.Value)
			return r
		},
		args...,
	)
}

func prod(args ...Any) Any {
	return fold_nums("*",
		func(a Int, b Int) Int {
			r := Int{big.NewInt(0)}
			r.Value = r.Value.Mul(a.Value, b.Value)
			return r
		},
		func(a Float, b Float) Float {
			r := ast.FloatNum(0)
			r.Value = r.Value.Mul(a.Value, b.Value)
			return r
		},
		Int{big.NewInt(1)}, args...,
	)
}

func div(args ...Any) Any {
	return numeric_2_floats("/",
		func(a Float, b Float) Any {
			r := ast.FloatNum(0)
			r.Value = r.Value.Quo(a.Value, b.Value)
			return r
		},
		args...,
	)
}

func numeric_2_floats(name string, f func(Float, Float) Any, args ...Any) Any {
	if len(args) != 2 {
		panic("'" + name + "' requires exactly 2 arguments, provided: " + LispyStr(args))
	}

	lhs := args[0]
	rhs := args[1]

	switch x := lhs.(type) {
	case Int:
		switch y := rhs.(type) {
		case Int:
			return f(int_to_float(x), int_to_float(y))
		case Float:
			return f(int_to_float(x), y)
		default:
			panic("Invalid '" + name + "' argument: " + LispyStr(args))
		}
	case Float:
		switch y := rhs.(type) {
		case Int:
			return f(x, int_to_float(y))
		case Float:
			return f(x, y)
		default:
			panic("Invalid '" + name + "' argument: " + LispyStr(args))
		}
	default:
		panic("'" + name + "' Error")
	}
}

func gt(args ...Any) Any {
	return numeric_2_floats(">", func(a Float, b Float) Any { return Bool(a.Value.Cmp(b.Value) > 0) }, args...)
}

func lt(args ...Any) Any {
	return numeric_2_floats(">", func(a Float, b Float) Any { return Bool(a.Value.Cmp(b.Value) < 0) }, args...)
}

func ge(args ...Any) Any {
	return numeric_2_floats(">", func(a Float, b Float) Any { return Bool(a.Value.Cmp(b.Value) >= 0) }, args...)
}

func le(args ...Any) Any {
	return numeric_2_floats(">", func(a Float, b Float) Any { return Bool(a.Value.Cmp(b.Value) <= 0) }, args...)
}

func list_cmp(a List, b List) Bool {
	if len(a) != len(b) {
		return false
	}

	for i, _ := range a {
		if !equal(a[i], b[i]) {
			return false
		}
	}

	return true
}

func equal(a Any, b Any) Bool {
	switch x := a.(type) {
	case List:
		switch y := b.(type) {
		case List:
			return list_cmp(x, y)
		default:
			return Bool(false)
		}
	case Int:
		switch y := b.(type) {
		case Int:
			return x.Value.Cmp(y.Value) == 0
		default:
			return Bool(false)
		}
	case Float:
		switch y := b.(type) {
		case Float:
			return x.Value.Cmp(y.Value) == 0
		default:
			return Bool(false)
		}
	default:
		return Bool(a == b)
	}
}

func eq(args ...Any) Any {
	if len(args) != 2 {
		panic("equality requires exactly 2 arguments, provided: " + LispyStr(args))
	}

	return equal(args[0], args[1])
}

//TODO: package constraints is not in GOROOT
//func max[T constraints.Ordered](a T, b T) T {
func max(a int, b int) int {
	if b < a {
		return a
	} else {
		return b
	}
}

func get(s []Any, i int, d Any) Any {
	if len(s) > i {
		return s[i]
	} else {
		return d
	}
}

func zip(v []Any) [][]Any {
	tmp := [][]Any{}
	size := 0
	for _, item := range v {
		l := to_list(item)
		tmp = append(tmp, l)
		size = max(size, len(l))
	}
	r := [][]Any{}
	for i := 0; i < size; i++ {
		item := []Any{}
		for _, row := range tmp {
			item = append(item, get(row, i, nil))
		}
		r = append(r, item)
	}
	return r
}

func lispy_map(args ...Any) Any {
	if len(args) < 2 {
		return nil
	}
	f := to_function(args[0])
	a := zip(args[1:])

	r := List{}
	for _, item := range a {
		r = append(r, f(item...))
	}
	return r
}

//TODO: flatten l1 unwraps only upper level of lists
func flatten_l1(args List) List {
	r := List{}

	for _, item := range args {
		switch v := item.(type) {
		case List:
			for _, i := range v {
				r = append(r, i)
			}
		default:
			r = append(r, v)
		}
	}

	return r
}

func apply(args ...Any) Any {
	if len(args) < 2 {
		return nil
	}
	f := to_function(args[0])
	a := flatten_l1(args[1:])
	return f(a...)
}

func to_float(n Any) Float {
	switch v := n.(type) {
	case Float:
		return v
	case Int:
		return int_to_float(v)
	default:
		panic("NAN: " + ast.String(n))
	}
}

func Zero() *big.Rat {
	r := big.NewRat(0, 1)
	return r
}

func Mul(a, b *big.Rat) *big.Rat {
	return Zero().Mul(a, b)
}

func Pow(a *big.Rat, e uint64) *big.Rat {
	result := Zero().Set(a)
	for i := uint64(0); i < e-1; i++ {
		result = Mul(result, a)
	}
	return result
}

func pow(args ...Any) Any {
	if len(args) != 2 {
		panic("'pow' requires exactly 2 arguments, provided: " + LispyStr(args))
	}

	base := to_float(args[0])
	exp := to_float(args[1])

	b, _ := base.Value.Float64()
	e, _ := exp.Value.Float64()

	return ast.FloatNum(math.Pow(b, e))
}

func fact(n *big.Int) *big.Int {
	p := big.NewInt(1)
	a := big.NewInt(1)
	one := big.NewInt(1)

	for a.Cmp(n) <= 0 {
		p.Mul(p, a)
		a.Add(a, one)
	}

	return p
}

func fact_r(n *big.Int) *big.Int {
	var p *big.Int

	if n.Sign() > 0 {
		next := big.NewInt(-1)
		next = next.Add(next, n)
		p = fact_r(next)
		return p.Mul(p, n)
	} else {
		return big.NewInt(1)
	}
}

func __fact(args ...Any) Any {
	arg := args[0]
	switch n := arg.(type) {
	case Int:
		return Int{fact(n.Value)}
	default:
		panic("Invalid '__fact' argument: " + LispyStr(args))
	}
}

func __fact_r(args ...Any) Any {
	arg := args[0]
	switch n := arg.(type) {
	case Int:
		return Int{fact_r(n.Value)}
	default:
		panic("Invalid '__fact' argument: " + LispyStr(args))
	}
}

func StdEnv() *Env {
	env := Env{}

	env.named_objects = map[string]Any{
		"enable-print-elapsed": Bool(false),
		"enable-trace":         Bool(false),

		"car":  car,
		"cdr":  cdr,
		"cons": cons,
		"list": list,
		"+":    sum,
		"-":    sub,
		"*":    prod,
		"/":    div,
		">":    gt,
		"<":    lt,
		">=":   ge,
		"<=":   le,
		"=":    eq,
		//TODO: common way to check the number of args and the types
		"begin":  func(args ...Any) Any { return args[len(args)-1] },
		"pi":     ast.FloatNum(math.Pi),
		"eq?":    func(args ...Any) Any { return Bool(args[0] == args[1]) },
		"equal?": eq,
		"length": func(args ...Any) Any { return Int{big.NewInt(int64(len(to_list(args[0]))))} },
		"not":    func(args ...Any) Any { return Bool(!if_test(args[0])) },
		"and":    func(args ...Any) Any { return fold_bools(func(x bool, y bool) bool { return x && y }, true, args...) },
		"or":     func(args ...Any) Any { return fold_bools(func(x bool, y bool) bool { return x || y }, false, args...) },
		"apply":  apply,
		"map":    lispy_map,

		"pow": pow,

		"__fact": __fact,
		"__fact_r": __fact_r,
	}

	return &env
}

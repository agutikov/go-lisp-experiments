package lispy

import (
	"math"
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
	if args[1] != nil {
		tail = to_list(args[1])
	}
	r = append(List{args[0]}, tail...)
	return r
}

func list(args ...Any) Any {
	return List(args)
}

func fold_nums(name string, f_i func(Int, Int) Int, f_f func(Float, Float) Float, init Any, args ...Any) Any {
	var acc Any
	acc = init
	//TODO: type switch for multiple args
	for _, item := range args {
		switch a := acc.(type) {
		case Int:
			switch v := item.(type) {
			case Int:
				acc = f_i(a, v)
			case Float:
				acc = f_f(Float(a), v)
			default:
				panic("Invalid '" + name + "' argument: " + LispyStr(args))
			}
		case Float:
			switch v := item.(type) {
			case Int:
				acc = f_f(a, Float(v))
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
			return f_f(Float(x), y)
		default:
			panic("Invalid '" + name + "' argument: " + LispyStr(args))
		}
	case Float:
		switch y := rhs.(type) {
		case Int:
			return f_f(x, Float(y))
		case Float:
			return f_f(x, y)
		default:
			panic("Invalid '" + name + "' argument: " + LispyStr(args))
		}
	default:
		panic("'" + name + "' Error")
	}
}

func sum(args ...Any) Any {
	return fold_nums("+",
		func(a Int, b Int) Int { return a + b },
		func(a Float, b Float) Float { return a + b },
		Int(0), args...,
	)
}

func minus(arg Any) Any {
	switch x := arg.(type) {
	case Int:
		return Int(-int64(x))
	case Float:
		return Float(-float64(x))
	default:
		panic("Invalid unary '-' argument: " + LispyStr(arg))
	}
}

func sub(args ...Any) Any {
	if len(args) == 1 {
		return minus(args[0])
	}
	return numeric_2_args("-",
		func(a Int, b Int) Int { return a - b },
		func(a Float, b Float) Float { return a - b },
		args...,
	)
}

func prod(args ...Any) Any {
	return fold_nums("*",
		func(a Int, b Int) Int { return a * b },
		func(a Float, b Float) Float { return a * b },
		Int(1), args...,
	)
}

func div(args ...Any) Any {
	return numeric_2_args("/",
		func(a Int, b Int) Int { return a / b },
		func(a Float, b Float) Float { return a / b },
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
			return f(Float(x), Float(y))
		case Float:
			return f(Float(x), y)
		default:
			panic("Invalid '" + name + "' argument: " + LispyStr(args))
		}
	case Float:
		switch y := rhs.(type) {
		case Int:
			return f(x, Float(y))
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
	return numeric_2_floats(">", func(a Float, b Float) Any { return Bool(a > b) }, args...)
}

func lt(args ...Any) Any {
	return numeric_2_floats(">", func(a Float, b Float) Any { return Bool(a < b) }, args...)
}

func ge(args ...Any) Any {
	return numeric_2_floats(">", func(a Float, b Float) Any { return Bool(a >= b) }, args...)
}

func le(args ...Any) Any {
	return numeric_2_floats(">", func(a Float, b Float) Any { return Bool(a <= b) }, args...)
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

func StdEnv() *Env {
	env := Env{}

	env.named_objects = map[Symbol]Any{
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
		"pi":     Float(math.Pi),
		"eq?":    func(args ...Any) Any { return Bool(args[0] == args[1]) },
		"equal?": eq,
		"length": func(args ...Any) Any { return Int(len(to_list(args[0]))) },
		"not":    func(args ...Any) Any { return Bool(!if_test(args[0])) },
		"apply":  apply,
		"map":    lispy_map,
	}

	return &env
}
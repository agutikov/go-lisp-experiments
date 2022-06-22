package lispy

import (
	"fmt"
	"reflect"
)

type Bool bool
type Int int64
type Str string
type Float float64

type Symbol string
type Builtin string
type Any interface{}
type List []Any

type PureFunction = func(...Any) Any

func to_symbol(s Any) Symbol {
	switch v := s.(type) {
	case Symbol:
		return v
	default:
		panic("Invalid symbol: " + lispstr(s))
	}
}

func to_list(s Any) List {
	switch v := s.(type) {
	case List:
		return v
	default:
		panic("Invalid list: " + lispstr(s))
	}
}

func to_function(s Any) PureFunction {
	switch v := s.(type) {
	case PureFunction:
		return v
	default:
		fmt.Println(reflect.TypeOf(v))
		panic("Invalid function: " + lispstr(s))
	}
}

func lispstr(expr Any) string {
	if expr == nil {
		return "nil"
	}
	switch v := expr.(type) {
	case Builtin:
		return string(v)
	case Symbol:
		return string(v)
	case Int:
		return fmt.Sprintf("%d", v)
	case Float:
		return fmt.Sprintf("%f", v)
	case Str:
		return "\"" + string(v) + "\""
	case Bool:
		if v {
			return "t"
		} else {
			return "f"
		}
	case List:
		s := "("
		for i, item := range v {
			if i != 0 {
				s += " "
			}
			s += lispstr(item)
		}
		s += ")"
		return s
	case PureFunction:
		return fmt.Sprintf("function{%p}", v)
	default:
		return "Unknown"
	}
}

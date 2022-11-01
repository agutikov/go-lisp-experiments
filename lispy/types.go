package lispy

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
)

type Bool = ast.Bool
type Int = ast.Int
type Float = ast.Float
type Str = ast.Str
type Symbol = ast.Symbol
type Any = ast.Any
type List = ast.List
type NIl = ast.Nil

type PureFunction = func([]Any) Any

func int_to_float(v Int) Float {
	r := new(big.Rat)
	r.SetInt(v.Value)
	return Float{r}
}

func to_symbol(s Any) Symbol {
	switch v := s.(type) {
	case Symbol:
		return v
	default:
		panic("Invalid symbol: " + LispyStr(s))
	}
}

func to_list(s Any) List {
	switch v := s.(type) {
	case List:
		return v
	default:
		panic("Invalid list: " + LispyStr(s))
	}
}

func to_int(s Any) Int {
	switch v := s.(type) {
	case Int:
		return v
	default:
		panic("Invalid int: " + LispyStr(s))
	}
}

func to_function(s Any) PureFunction {
	switch v := s.(type) {
	case PureFunction:
		return v
	default:
		panic("Invalid function: " + LispyStr(s) + fmt.Sprintf(" type: %v", reflect.TypeOf(v)) + "; probably evaluated the unquoted list")
	}
}

func if_test(value Any) bool {
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case ast.Nil:
		return false
	case Bool:
		return bool(v)
	case List:
		return len(v) > 0
	case Int:
		return v.Value.Sign() != 0
	case Str:
		return len(v.Value) > 0
	case Float:
		return v.Value.Sign() != 0
	default:
		return true
	}
}

func is_nil(value Any) bool {
	if value == nil {
		return true
	}
	switch value.(type) {
	case ast.Nil:
		return true
	default:
		return false
	}
}

func LispyStr(expr Any) string {
	return ast.String(expr)
}

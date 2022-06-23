package lispy

import (
	"fmt"
	"math/big"
	"reflect"
)

type Bool bool
type Int struct {
	v *big.Int
}
type Str string
type Float struct {
	v *big.Float
}

type Symbol string
type Builtin string
type Any interface{}
type List []Any

type PureFunction = func(...Any) Any

func FromInt(v int) Int {
	return Int{big.NewInt(int64(v))}
}

func FromFloat(v float64) Float {
	return Float{big.NewFloat(v)}
}

func int_to_float(v Int) Float {
	r := new(big.Float)
	r.SetInt(v.v)
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
	case Bool:
		return bool(v)
	case List:
		return len(v) > 0
	case Int:
		return v.v.Sign() != 0
	case Str:
		return len(v) > 0
	case Float:
		return v.v.Sign() != 0
	default:
		return true
	}
}

func LispyStr(expr Any) string {
	if expr == nil {
		return "nil"
	}
	switch v := expr.(type) {
	case Builtin:
		return string(v)
	case Symbol:
		return string(v)
	case Int:
		return v.v.String()
	case Float:
		return v.v.String()
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
			s += LispyStr(item)
		}
		s += ")"
		return s
	case PureFunction:
		return fmt.Sprintf("function{%p}", v)
	default:
		return "Unknown"
	}
}

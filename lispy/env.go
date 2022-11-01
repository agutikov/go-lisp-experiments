package lispy

import (
	"fmt"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
)

type LambdaArg struct {
	index int
}

type Env struct {
	parent        *Env
	named_objects map[string]Any
}

func (env *Env) Print() Any {
	fmt.Printf("named_objects: %#v\n", env.named_objects)
	fmt.Printf("parent: %p\n", env.parent)
	if env.parent != nil {
		env.parent.Print()
	}
	return nil
}

func newEnv(parent *Env) *Env {
	e := Env{parent: parent, named_objects: map[string]Any{}}
	return &e
}

func (env *Env) assign_vars(vars []ast.Symbol, values ...Any) {
	if len(vars) != len(values) {
		panic("Invalid number of values provided")
	}

	for i := range vars {
		env.named_objects[vars[i].Name] = values[i]
	}
}

func (env *Env) define_lambda_args(symbols []ast.Symbol) {
	for i, sym := range symbols {
		env.named_objects[sym.Name] = LambdaArg{index: i}
	}
}

func (env *Env) symbol_lookup(s Symbol) Any {
	if val, ok := env.named_objects[s.Name]; ok {
		return val
	} else if env.parent != nil {
		return env.parent.symbol_lookup(s)
	} else {
		//TODO: panic or return nil or error?
		panic("Undefined symbol: \"" + s.Name + "\"")
	}
}

func (env *Env) lambda_symbol_lookup(s Symbol) (Any, bool) {
	val, ok := env.named_objects[s.Name]
	if ok {
		return val, true
	}
	if env.parent != nil {
		return env.parent.lambda_symbol_lookup(s)
	}
	return nil, false
}

func (env *Env) env_lookup(s string) *Env {
	if _, ok := env.named_objects[s]; ok {
		return env
	} else if env.parent != nil {
		return env.parent.env_lookup(s)
	} else {
		//TODO: panic or return nil or error?
		panic("Undefined symbol: \"" + s + "\"")
	}
}

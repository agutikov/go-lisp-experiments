package lispy

import (
	"fmt"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
)

type LambdaPreEvalContext struct {
	env               *Env
	arg_name_to_index map[string]int
}

type LambdaCallContext struct {
	args []Any
}

func newLambdaPreEvalContext(e *Env, args ast.LambdaArgs) *LambdaPreEvalContext {
	name_to_index := map[string]int{}
	for i, sym := range args {
		name_to_index[sym.Name] = i
	}
	c := LambdaPreEvalContext{env: e, arg_name_to_index: name_to_index}
	return &c
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

package lispy

import "fmt"

type Env struct {
	parent        *Env
	named_objects map[Symbol]Any
}

func (env *Env) print() Any {
	fmt.Println("named_objects: ", env.named_objects)
	fmt.Println("parent: ", env.parent)
	fmt.Println()
	if env.parent != nil {
		env.parent.print()
	}
	return nil
}

func newEnv(parent *Env) *Env {
	e := Env{parent: parent, named_objects: map[Symbol]Any{}}
	return &e
}

func (env *Env) assign_vars(vars []Symbol, values List) {
	if len(vars) != len(values) {
		panic("Invalid number of values provided")
	}

	for i := range vars {
		env.named_objects[vars[i]] = values[i]
	}
}

func (env *Env) symbol_lookup(s Symbol) Any {
	if val, ok := env.named_objects[s]; ok {
		return val
	} else if env.parent != nil {
		return env.parent.symbol_lookup(s)
	} else {
		//TODO: panic or return nil or error?
		panic("Undefined symbol: \"" + s + "\"")
	}
}

func (env *Env) env_lookup(s Symbol) *Env {
	if _, ok := env.named_objects[s]; ok {
		return env
	} else if env.parent != nil {
		return env.parent.env_lookup(s)
	} else {
		//TODO: panic or return nil or error?
		panic("Undefined symbol: \"" + s + "\"")
	}
}

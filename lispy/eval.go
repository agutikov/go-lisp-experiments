package lispy

import (
	"fmt"
	"time"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
)

func (env *Env) eval_if(expr ast.If) Any {
	v := env.eval_expr(expr.Test)

	var r ast.Any

	if if_test(v) {
		r = expr.PosBranch
	} else {
		r = expr.NegBranch
	}

	return env.eval_expr(r)
}

func (env *Env) eval_define(d ast.Define) Any {
	v := env.eval_expr(d.Value)
	env.named_objects[d.Sym.Name] = v
	return v
}

func (env *Env) eval_set(s ast.Set) Any {
	value := env.eval_expr(s.Value)
	env.env_lookup(s.Sym.Name).named_objects[s.Sym.Name] = value
	return value
}

func (env *Env) eval_lambda(l ast.Lambda) Any {
	// Return callable which will
	return func(args ...Any) Any {
		// Eval body in the new nested environment
		e := newEnv(env)
		e.assign_vars(l.Args, args...)
		return e.eval_expr(l.Body)
	}
}

func (env *Env) eval_quote(q ast.Quote) Any {
	//TODO: unquote
	return q.Value
}

func (env *Env) eval_args(args ...ast.Any) []Any {
	r := []Any{}
	for _, elem := range args {
		r = append(r, env.eval_expr(elem))
	}
	return r
}

func (env *Env) eval_list(lst List) Any {
	if len(lst) == 0 {
		return lst
	}
	head := lst[0]
	tail := lst[1:]
	f_value := env.eval_expr(head)
	f := to_function(f_value)
	args := env.eval_args(tail...)
	return f(args...)
}

func quote_if_list(value Any) Any {
	switch v := value.(type) {
	case List:
		return ast.Quote{Value: v}
	default:
		return v
	}
}

func (env *Env) eval_sequence(seq ast.Sequence) Any {
	var r Any
	r = nil
	for _, expr := range seq {
		r = env.eval_expr(expr)
	}
	return r
}

func (env *Env) _eval_expr(expr Any) Any {
	switch v := expr.(type) {
	case List:
		return env.eval_list(v)
	case ast.Sequence:
		return env.eval_sequence(v)
	case ast.Quote:
		return env.eval_quote(v)
	case ast.Define:
		return env.eval_define(v)
	case ast.If:
		return env.eval_if(v)
	case ast.Set:
		return env.eval_set(v)
	case ast.Lambda:
		return env.eval_lambda(v)
	case Symbol:
		// Symbol atom is a name of object in the environment
		return env.symbol_lookup(v)
	default:
		// Other atoms are const literals
		return v
	}
}

func (env *Env) eval_expr(expr Any) Any {
	//fmt.Printf("eval_expr(%#v)\n", expr)

	r := env._eval_expr(expr)

	//env.Print()
	if if_test(env.symbol_lookup(ast.Symbol{"enable-trace"})) {
		fmt.Printf("eval_expr():  %s  ->  %s \n", LispyStr(expr), LispyStr(r))
	}
	//fmt.Printf("eval_expr():  %#v  ->  %#v \n", expr, r)
	return r
}

func (env *Env) Eval(expr Any) Any {
	started := time.Now()

	r := quote_if_list(env.eval_expr(expr))

	elapsed := time.Since(started)
	if if_test(env.symbol_lookup(ast.Symbol{"enable-print-elapsed"})) {
		fmt.Println(" elapsed: ", elapsed)
	}

	return r
}

func Lambda(s string) PureFunction {
	return to_function(StdEnv().eval_expr(ParseStr(s)))
}

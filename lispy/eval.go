package lispy

import (
	"fmt"
	"time"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
)

type LambdaPlaceholder struct {
	f func([]Any) Any
}

func (env *Env) eval_defun(df ast.Defun) func([]Any) Any {
	// Temporary Env that we will use for body pre-eval
	pre_eval_env := newEnv(env)

	// Placeholder for function instance that will appear after function body pre-eval
	fwd := LambdaPlaceholder{}

	// Put the placeholder into the temporary env - to allow recursive function find it's name
	pre_eval_env.named_objects[df.Sym.Name] = func(args []Any) Any {
		return fwd.f(args)
	}

	// Eval the lambda in the temporary env
	lambda := pre_eval_env.eval_lambda(df.L)

	// Put the lambda into the placeholder
	fwd.f = lambda

	// Put the lambda into the env (complete 'define')
	env.named_objects[df.Sym.Name] = lambda

	// Return lambda (as 'define' does)
	return lambda
}

func (env *Env) eval_lambda(l ast.Lambda) func([]Any) Any {
	pre_eval_ctx := newLambdaPreEvalContext(env, l.Args)

	// pre-eval lambda body in the temporary env
	r := pre_eval_ctx.lambda_eval_body(l.Body)

	if r.is_constant {
		return func([]Any) Any {
			return r.value
		}
	} else {
		return func(args []Any) Any {
			ctx := LambdaCallContext{args: args}
			return r.function(&ctx)
		}
	}
}

func (env *Env) eval_simple_lambda(l ast.SimpleLambda) Any {
	// Return callable which will
	return func(args []Any) Any {
		// eval body in the new nested environment
		e := newEnv(env)
		e.assign_vars(l.Args, args)
		return e.eval_expr(l.Body)
	}
}

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

func (env *Env) eval_quoted_expr(expr Any) Any {
	switch v := expr.(type) {
	case List:
		return env.eval_quoted_list(v)
	case ast.Unquote:
		return env.eval_expr(v.Value)
	default:
		return v
	}
}

func (env Env) eval_quoted_list(args List) Any {
	lst := List{}
	for _, item := range args {
		lst = append(lst, env.eval_quoted_expr(item))
	}
	return lst
}

func (env *Env) eval_quote(q ast.Quote) Any {
	return env.eval_quoted_expr(q.Value)
}

func (env *Env) eval_args(args []ast.Any) []Any {
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

	// eval car into callable
	f_value := env.eval_expr(head)
	f := to_function(f_value)

	// eval args
	args := env.eval_args(tail)

	// call function with args
	return f(args)
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
		started := time.Now()

		r = env.eval_expr(expr)

		if if_test(env.symbol_lookup(ast.Symbol{"enable-trace"})) {
			fmt.Printf("%s  ->  %s \n", LispyStr(expr), LispyStr(r))
		}

		elapsed := time.Since(started)
		if if_test(env.symbol_lookup(ast.Symbol{"enable-print-elapsed"})) {
			fmt.Println(" elapsed: ", elapsed)
		}
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
	case ast.SimpleLambda:
		return env.eval_simple_lambda(v)
	case ast.Lambda:
		return env.eval_lambda(v)
	case ast.Defun:
		return env.eval_defun(v)
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
	//fmt.Printf("eval_expr():  %#v  ->  %#v \n", expr, r)
	return r
}

func (env *Env) Eval(seq ast.Sequence) Any {

	r := quote_if_list(env.eval_sequence(seq))

	return r
}

func Function(s string) PureFunction {
	return to_function(StdEnv().eval_expr(ParseStr(s)))
}

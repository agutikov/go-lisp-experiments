package lispy

import (
	"fmt"
	"time"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
)

func (env *Env) lambda_eval_if(expr ast.If) func(*Env) Any {
	// Pre-eval test and all branches
	test := env.lambda_eval_body(expr.Test)
	pos_branch := env.lambda_eval_body(expr.PosBranch)
	neg_branch := env.lambda_eval_body(expr.NegBranch)

	// Return callable
	return func(e *Env) Any {
		if if_test(test(e)) {
			return pos_branch(e)
		} else {
			return neg_branch(e)
		}
	}
}

func (env *Env) lambda_eval_quoted_expr(expr Any) func(*Env) Any {
	switch v := expr.(type) {
	case List:
		return env.lambda_eval_quoted_list(v) //TODO: optimize
	case ast.Unquote:
		return env.lambda_eval_body(v.Value)
	default:
		return func(*Env) Any {
			return v
		}
	}
}

func (env Env) lambda_eval_quoted_list(args List) func(*Env) Any {
	f_lst := []func(*Env) Any{}
	for _, item := range args {
		f_lst = append(f_lst, env.lambda_eval_quoted_expr(item))
	}
	return func(env *Env) Any {
		lst := List{}
		for _, f_item := range f_lst {
			lst = append(lst, f_item(env))
		}
		return lst
	}
}

func (env *Env) lambda_eval_quote(q ast.Quote) func(*Env) Any {
	return env.lambda_eval_quoted_expr(q.Value)
}

func (env *Env) lambda_eval_define(d ast.Define) func(*Env) Any {
	//TODO
	panic("define inside lambda body not implemented")
}

func (env *Env) lambda_eval_set(s ast.Set) func(*Env) Any {
	//TODO
	panic("set! inside lambda body not implemented")
}

func (env *Env) lambda_eval_lambda(lambda ast.Lambda) func(*Env) Any {
	//TODO
	panic("lambda inside lambda body not implemented")
}

func (env *Env) lambda_eval_defun(df ast.Defun) func(*Env) Any {
	//TODO
	panic("defun inside lambda body not implemented")
}

// Call a function inside lambda body
func (env *Env) lambda_eval_list(lst List) func(*Env) Any {
	if len(lst) == 0 {
		return func(*Env) Any {
			return lst
		}
	}
	head := lst[0]
	tail := lst[1:]

	// pre-eval car into callable that will return function
	get_f := env.lambda_eval_body(head)
	//f_value := env.eval_expr(head)
	//f := to_function(f_value)

	// pre-eval args
	args_f := []func(*Env) Any{}
	for _, elem := range tail {
		args_f = append(args_f, env.lambda_eval_body(elem))
	}

	return func(e *Env) Any {
		// get the function
		f_value := get_f(e)
		f := to_function(f_value)

		// eval args values with env
		args := []Any{}
		for _, arg_f := range args_f {
			args = append(args, arg_f(e))
		}

		// Call the function
		return f(args...)
	}
}

func (env *Env) lambda_eval_symbol(sym Symbol) func(*Env) Any {
	// Inside lambda body
	value, ok := env.lambda_symbol_lookup(sym)
	if ok {
		// If symbol already exists
		switch v := value.(type) {
		case LambdaArg:
			// If it is an argument use - return callable
			return func(e *Env) Any {
				// that takes the argument by index from args
				return e.lambda_args[v.index]
			}
		default:
			// If anything else - just cache a value
			return func(*Env) Any {
				return value
			}
		}
	}
	// If symbol not defined yet - will lookup it when called
	return func(e *Env) Any {
		return e.symbol_lookup(sym)
	}
}

type LambdaArg struct {
	index int
}

// Pre-eval lambda body into function with single Env argument
func (env *Env) lambda_eval_body(item Any) func(*Env) Any {
	switch v := item.(type) {
	case List:
		return env.lambda_eval_list(v)
	case ast.Sequence:
		// Sequence is not possible inside lambda body
		panic("Lambda pre-eval ERROR: Sequence appears")
	case ast.Quote:
		return env.lambda_eval_quote(v)
	case ast.Define:
		return env.lambda_eval_define(v)
	case ast.Defun:
		return env.lambda_eval_defun(v)
	case ast.If:
		return env.lambda_eval_if(v)
	case ast.Set:
		return env.lambda_eval_set(v)
	case ast.Lambda:
		return env.lambda_eval_lambda(v)
	case Symbol:
		return env.lambda_eval_symbol(v)
	default:
		// Other atoms are const literals
		return func(*Env) Any {
			return v
		}
	}
}

type LambdaPlaceholder struct {
	f func(...Any) Any
}

func (env *Env) eval_defun(df ast.Defun) func(...Any) Any {
	// Temporary Env that we will use for body pre-eval
	pre_eval_env := newEnv(env)

	// Placeholder for function instance that will appear after function body pre-eval
	fwd := LambdaPlaceholder{}

	// Put the placeholder into the temporary env - to allow recursive function find it's name
	pre_eval_env.named_objects[df.Sym.Name] = func(args ...Any) Any {
		return fwd.f(args...)
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

func (env *Env) eval_lambda(l ast.Lambda) func(...Any) Any {
	// Env that is used during pre-eval of lambda body
	pre_eval_env := newEnv(env)

	// Put LambdaArg objects into the temporary env
	pre_eval_env.define_lambda_args(l.Args)

	// pre-eval lambda body in the temporary env
	body_f := pre_eval_env.lambda_eval_body(l.Body)

	// Return callable that
	return func(args ...Any) Any {
		// Creates env to run the function in - like a call frame
		e := newEnv(env)

		// Put args into the "call frame"
		e.lambda_args = args

		// Call a pre-evaluated lambda
		return body_f(e)
	}
}

func (env *Env) eval_simple_lambda(l ast.SimpleLambda) Any {
	// Return callable which will
	return func(args ...Any) Any {
		// eval body in the new nested environment
		e := newEnv(env)
		e.assign_vars(l.Args, args...)
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

	// eval car into callable
	f_value := env.eval_expr(head)
	f := to_function(f_value)

	// eval args
	args := env.eval_args(tail...)

	// call function with args
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

func Lambda(s string) PureFunction {
	return to_function(StdEnv().eval_expr(ParseStr(s)))
}

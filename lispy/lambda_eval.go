package lispy

import (
	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
)

func (env *Env) lambda_eval_if(expr ast.If) func(...Any) Any {
	// Pre-eval test and all branches
	test := env.lambda_eval_body(expr.Test)
	pos_branch := env.lambda_eval_body(expr.PosBranch)
	neg_branch := env.lambda_eval_body(expr.NegBranch)

	// Return callable
	return func(args ...Any) Any {
		if if_test(test(args...)) {
			return pos_branch(args...)
		} else {
			return neg_branch(args...)
		}
	}
}

func (env *Env) lambda_eval_quoted_expr(expr Any) func(...Any) Any {
	switch v := expr.(type) {
	case List:
		return env.lambda_eval_quoted_list(v)
	case ast.Unquote:
		return env.lambda_eval_body(v.Value)
	default:
		return func(...Any) Any {
			return v
		}
	}
}

func (env Env) lambda_eval_quoted_list(args List) func(...Any) Any {
	f_lst := []func(...Any) Any{}
	for _, item := range args {
		f_lst = append(f_lst, env.lambda_eval_quoted_expr(item))
	}
	return func(args ...Any) Any {
		lst := List{}
		for _, f_item := range f_lst {
			lst = append(lst, f_item(args...))
		}
		return lst
	}
}

func (env *Env) lambda_eval_quote(q ast.Quote) func(...Any) Any {
	return env.lambda_eval_quoted_expr(q.Value)
}

func (env *Env) lambda_eval_lambda(lambda ast.Lambda) func(...Any) Any {
	//TODO
	panic("lambda inside lambda body not implemented")
}

func (env *Env) lambda_eval_set(s ast.Set) func(...Any) Any {
	//TODO
	panic("set! inside lambda body not implemented")
}

func (env *Env) lambda_eval_define(d ast.Define) func(...Any) Any {
	//TODO
	panic("define inside lambda body not implemented")
}

func (env *Env) lambda_eval_defun(df ast.Defun) func(...Any) Any {
	//TODO
	panic("defun inside lambda body not implemented")
}

func (env *Env) lambda_eval_args(args ...ast.Any) []func(...Any) Any {
	r := []func(...Any) Any{}
	for _, elem := range args {
		r = append(r, env.lambda_eval_body(elem))
	}
	return r
}

// Call a function inside lambda body
func (env *Env) lambda_eval_list(lst List) func(...Any) Any {
	if len(lst) == 0 {
		return func(...Any) Any {
			return List{}
		}
	}
	head := lst[0]
	tail := lst[1:]

	// pre-eval car into callable that will return function
	get_f := env.lambda_eval_body(head)

	// pre-eval args
	args_f := env.lambda_eval_args(tail...)

	return func(args ...Any) Any {
		// get the function
		f_value := get_f(args...)
		f := to_function(f_value)

		// eval args values with env
		values := []Any{}
		for _, arg_f := range args_f {
			values = append(values, arg_f(args...))
		}

		// Call the function
		return f(values...)
	}
}

func (env *Env) lambda_eval_symbol(sym Symbol) func(...Any) Any {
	// Inside lambda body
	value := env.symbol_lookup(sym)

	// If symbol already exists
	switch v := value.(type) {
	case LambdaArg:
		// If it is an argument use - return callable
		return func(args ...Any) Any {
			// that takes the argument by index from args
			return args[v.index]
		}
	default:
		// If anything else - just cache a value
		return func(...Any) Any {
			return value
		}
	}
}

// Pre-eval lambda body into function with single Env argument
func (env *Env) lambda_eval_body(item Any) func(...Any) Any {
	switch v := item.(type) {
	case List:
		return env.lambda_eval_list(v)
	case ast.Sequence:
		// Sequence is not possible inside lambda body
		panic("Lambda pre-eval ERROR: Sequence")
	case ast.SimpleLambda:
		panic("slambda inside lambda not allowed")
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
		return func(...Any) Any {
			return v
		}
	}
}

package lispy

import (
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
		return env.lambda_eval_quoted_list(v)
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

func (env *Env) lambda_eval_lambda(lambda ast.Lambda) func(*Env) Any {
	l_func := env.eval_lambda(lambda)

	return func(*Env) Any {
		return l_func
	}
}

func (env *Env) lambda_eval_set(s ast.Set) func(*Env) Any {
	// Inside lambda body
	body_f := env.lambda_eval_body(s.Value)

	value, ok := env.lambda_symbol_lookup(s.Sym)
	if ok {
		// If symbol already exists
		switch v := value.(type) {
		case LambdaArg:
			// If it is an argument
			return func(e *Env) Any {
				value := body_f(e)
				e.lambda_args[v.index] = value
				return value
			}
		default:
			// If anything else
			return func(*Env) Any {
				//TODO: cache env where to set value
			}
		}
	}
	// If symbol not defined yet - will lookup it when called
	return func(e *Env) Any {
	}
}

func (env *Env) lambda_eval_define(d ast.Define) func(*Env) Any {
	//TODO
	panic("define inside lambda body not implemented")
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

// Pre-eval lambda body into function with single Env argument
func (env *Env) lambda_eval_body(item Any) func(*Env) Any {
	switch v := item.(type) {
	case List:
		return env.lambda_eval_list(v)
	case ast.Sequence:
		// Sequence is not possible inside lambda body
		panic("Lambda pre-eval ERROR: Sequence appears")
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
		return func(*Env) Any {
			return v
		}
	}
}

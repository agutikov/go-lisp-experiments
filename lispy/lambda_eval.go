package lispy

import (
	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
)

func (env *LambdaPreEvalContext) lambda_eval_if(expr ast.If) func(*LambdaCallContext) Any {
	// Pre-eval test and all branches
	test := env.lambda_eval_body(expr.Test)
	pos_branch := env.lambda_eval_body(expr.PosBranch)
	neg_branch := env.lambda_eval_body(expr.NegBranch)

	// Return callable
	return func(ctx *LambdaCallContext) Any {
		if if_test(test(ctx)) {
			return pos_branch(ctx)
		} else {
			return neg_branch(ctx)
		}
	}
}

func (env *LambdaPreEvalContext) lambda_eval_quoted_expr(expr Any) func(*LambdaCallContext) Any {
	switch v := expr.(type) {
	case List:
		return env.lambda_eval_quoted_list(v)
	case ast.Unquote:
		return env.lambda_eval_body(v.Value)
	default:
		return func(*LambdaCallContext) Any {
			return v
		}
	}
}

func (env *LambdaPreEvalContext) lambda_eval_quoted_list(args List) func(*LambdaCallContext) Any {
	f_lst := []func(*LambdaCallContext) Any{}
	for _, item := range args {
		f_lst = append(f_lst, env.lambda_eval_quoted_expr(item))
	}
	return func(ctx *LambdaCallContext) Any {
		lst := List{}
		for _, f_item := range f_lst {
			lst = append(lst, f_item(ctx))
		}
		return lst
	}
}

func (env *LambdaPreEvalContext) lambda_eval_quote(q ast.Quote) func(*LambdaCallContext) Any {
	return env.lambda_eval_quoted_expr(q.Value)
}

func (env *LambdaPreEvalContext) lambda_eval_lambda(lambda ast.Lambda) func(*LambdaCallContext) Any {
	//TODO
	panic("lambda inside lambda body not implemented")
}

func (env *LambdaPreEvalContext) lambda_eval_set(s ast.Set) func(*LambdaCallContext) Any {
	//TODO
	panic("set! inside lambda body not implemented")
}

func (env *LambdaPreEvalContext) lambda_eval_define(d ast.Define) func(*LambdaCallContext) Any {
	//TODO
	panic("define inside lambda body not implemented")
}

func (env *LambdaPreEvalContext) lambda_eval_defun(df ast.Defun) func(*LambdaCallContext) Any {
	//TODO
	panic("defun inside lambda body not implemented")
}

func (env *LambdaPreEvalContext) lambda_eval_args(args []ast.Any) []func(*LambdaCallContext) Any {
	r := []func(*LambdaCallContext) Any{}
	for _, elem := range args {
		r = append(r, env.lambda_eval_body(elem))
	}
	return r
}

// Call a function inside lambda body
func (env *LambdaPreEvalContext) lambda_eval_list(lst List) func(*LambdaCallContext) Any {
	if len(lst) == 0 {
		return func(*LambdaCallContext) Any {
			return List{}
		}
	}
	head := lst[0]
	tail := lst[1:]

	// pre-eval car into callable that will return function
	get_f := env.lambda_eval_body(head)

	// pre-eval args
	args_f := env.lambda_eval_args(tail)

	return func(ctx *LambdaCallContext) Any {
		// get the function
		f_value := get_f(ctx)
		f := to_function(f_value)

		// eval args values with env
		values := []Any{}
		for _, arg_f := range args_f {
			values = append(values, arg_f(ctx))
		}

		// Call the function
		return f(values)
	}
}

func (env *LambdaPreEvalContext) lambda_eval_symbol(sym Symbol) func(*LambdaCallContext) Any {
	// Inside lambda body
	if arg_index, ok := env.arg_name_to_index[sym.Name]; ok {
		return func(ctx *LambdaCallContext) Any {
			return ctx.args[arg_index]
		}
	}

	value := env.env.symbol_lookup(sym)

	return func(*LambdaCallContext) Any {
		return value
	}
}

// Pre-eval lambda body into function with single Env argument
func (env *LambdaPreEvalContext) lambda_eval_body(item Any) func(*LambdaCallContext) Any {
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
		return func(*LambdaCallContext) Any {
			return v
		}
	}
}

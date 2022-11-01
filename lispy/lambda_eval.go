package lispy

import (
	"github.com/agutikov/go-lisp-experiments/lispy/syntax/ast"
)

type LambdaEvalBodyResult struct {
	is_constant bool
	value       Any
	function    func(*LambdaCallContext) Any
}

func constResult(value Any) LambdaEvalBodyResult {
	return LambdaEvalBodyResult{
		is_constant: true,
		value:       value,
	}
}

func funcResult(f func(*LambdaCallContext) Any) LambdaEvalBodyResult {
	return LambdaEvalBodyResult{
		is_constant: false,
		function:    f,
	}
}

func (env *LambdaPreEvalContext) lambda_eval_if(expr ast.If) LambdaEvalBodyResult {
	// Pre-eval test and all branches
	test := env.lambda_eval_body(expr.Test)
	pos_branch := env.lambda_eval_body(expr.PosBranch)
	neg_branch := env.lambda_eval_body(expr.NegBranch)

	if test.is_constant {
		if if_test(test.value) {
			return pos_branch
		} else {
			return neg_branch
		}
	}

	if pos_branch.is_constant {
		if neg_branch.is_constant {
			return funcResult(func(ctx *LambdaCallContext) Any {
				if if_test(test.function(ctx)) {
					return pos_branch.value
				} else {
					return neg_branch.value
				}
			})
		} else {
			return funcResult(func(ctx *LambdaCallContext) Any {
				if if_test(test.function(ctx)) {
					return pos_branch.value
				} else {
					return neg_branch.function(ctx)
				}
			})
		}
	} else {
		if neg_branch.is_constant {
			return funcResult(func(ctx *LambdaCallContext) Any {
				if if_test(test.function(ctx)) {
					return pos_branch.function(ctx)
				} else {
					return neg_branch.value
				}
			})
		} else {
			return funcResult(func(ctx *LambdaCallContext) Any {
				if if_test(test.function(ctx)) {
					return pos_branch.function(ctx)
				} else {
					return neg_branch.function(ctx)
				}
			})
		}
	}
}

func (env *LambdaPreEvalContext) lambda_eval_quoted_expr(expr Any) LambdaEvalBodyResult {
	panic("not implemented")

	/*
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
	*/
}

func (env *LambdaPreEvalContext) lambda_eval_quoted_list(args List) LambdaEvalBodyResult {
	panic("not implemented")

	/*
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
	*/
}

func (env *LambdaPreEvalContext) lambda_eval_quote(q ast.Quote) LambdaEvalBodyResult {
	return env.lambda_eval_quoted_expr(q.Value)
}

func (env *LambdaPreEvalContext) lambda_eval_lambda(lambda ast.Lambda) LambdaEvalBodyResult {
	//TODO
	panic("lambda inside lambda body not implemented")
}

func (env *LambdaPreEvalContext) lambda_eval_set(s ast.Set) LambdaEvalBodyResult {
	//TODO
	panic("set! inside lambda body not implemented")
}

func (env *LambdaPreEvalContext) lambda_eval_define(d ast.Define) LambdaEvalBodyResult {
	//TODO
	panic("define inside lambda body not implemented")
}

func (env *LambdaPreEvalContext) lambda_eval_defun(df ast.Defun) LambdaEvalBodyResult {
	//TODO
	panic("defun inside lambda body not implemented")
}

func (env *LambdaPreEvalContext) lambda_eval_args(args []ast.Any) []LambdaEvalBodyResult {
	r := []LambdaEvalBodyResult{}
	for _, elem := range args {
		r = append(r, env.lambda_eval_body(elem))
	}
	return r
}

func (ctx *LambdaCallContext) lambda_eval_args(args []LambdaEvalBodyResult) []Any {
	values := []Any{}
	for _, arg := range args {
		if arg.is_constant {
			values = append(values, arg.value)
		} else {
			values = append(values, arg.function(ctx))
		}
	}
	return values
}

// Call a function inside lambda body
func (env *LambdaPreEvalContext) lambda_eval_list(lst List) LambdaEvalBodyResult {
	if len(lst) == 0 {
		return constResult(List{})
	}
	head := lst[0]
	tail := lst[1:]

	// pre-eval car into callable that will return function
	f_r := env.lambda_eval_body(head)

	// pre-eval args
	args := env.lambda_eval_args(tail)

	if f_r.is_constant {
		f := to_function(f_r.value)

		return funcResult(func(ctx *LambdaCallContext) Any {
			values := ctx.lambda_eval_args(args)

			return f(values)
		})
	} else {
		return funcResult(func(ctx *LambdaCallContext) Any {
			// get the function
			f_value := f_r.function(ctx)
			f := to_function(f_value)

			// eval args values with env
			values := ctx.lambda_eval_args(args)

			// Call the function
			return f(values)
		})
	}
}

func (env *LambdaPreEvalContext) lambda_eval_symbol(sym Symbol) LambdaEvalBodyResult {
	// Inside lambda body
	if arg_index, ok := env.arg_name_to_index[sym.Name]; ok {
		return funcResult(func(ctx *LambdaCallContext) Any {
			return ctx.args[arg_index]
		})
	}

	value := env.env.symbol_lookup(sym)

	return constResult(value)
}

// Pre-eval lambda body into function with single Env argument
func (env *LambdaPreEvalContext) lambda_eval_body(item Any) LambdaEvalBodyResult {
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
		return constResult(v)
	}
}

package lispy

func (env *Env) eval_if(expr ...Any) Any {
	//TODO: verify number of args on parsing stage
	if len(expr) != 4 {
		panic("'if' statement requires exactly 3 arguments (if test conseq alt), while provided: " + LispyStr(expr))
	}
	test := expr[1]
	conseq := expr[2]
	alt := expr[3]

	v := env.Eval(test)

	var r Any

	if if_test(v) {
		r = conseq
	} else {
		r = alt
	}

	return env.Eval(r)
}

func (env *Env) eval_define(expr ...Any) Any {
	//TODO: verify number of args on parsing stage
	if len(expr) != 3 {
		panic("'define' statement requires exactly 2 arguments (define name exp), while provided: " + LispyStr(expr))
	}
	name := expr[1]
	exp := expr[2]

	value := env.Eval(exp)

	switch s := name.(type) {
	case Symbol:
		env.named_objects[s] = value
	default:
		panic("Invalid define name argument")
	}

	return value
}

func (env *Env) eval_set(expr ...Any) Any {
	//TODO: verify number of args on parsing stage
	if len(expr) != 3 {
		panic("'set!' statement requires exactly 2 arguments (set! name exp), while provided: " + LispyStr(expr))
	}
	name := expr[1]
	exp := expr[2]

	value := env.Eval(exp)

	switch s := name.(type) {
	case Symbol:
		env.env_lookup(s).named_objects[s] = value
	default:
		panic("Invalid set! name argument")
	}

	return value
}

func lambda_args(vars Any) []Symbol {
	args := to_list(vars)

	a := []Symbol{}

	for _, elem := range args {
		a = append(a, to_symbol(elem))
	}

	return a
}

func (env *Env) eval_lambda(expr ...Any) Any {
	//TODO: verify number of args on parsing stage
	if len(expr) != 3 {
		panic("'lambda' statement requires exactly 2 arguments (lambda (vars...) body), while provided: " + LispyStr(expr))
	}
	p := expr[1]
	body := expr[2]

	params := lambda_args(p)

	// Return callable which will
	return func(args ...Any) Any {
		// Eval body in the new nested environment
		e := newEnv(env)
		e.assign_vars(params, args)
		return e.Eval(body)
	}
}

func eval_quote(expr ...Any) Any {
	if len(expr) != 2 {
		panic("'quote' statement requires exactly 1 argument (quote exp), while provided: " + LispyStr(expr))
	}
	return expr[1]
}

func (env *Env) eval_builtin(s Builtin, expr List) Any {
	switch s { //TODO: is there any benefit of using map[Symbol]func(*Env, ...Any)Any ?
	case "quote":
		//TODO: unquote????
		return eval_quote(expr...)
	case "if":
		return env.eval_if(expr...)
	case "define":
		//TODO: what define should return ?
		return env.eval_define(expr...)
	case "set!":
		//TODO: what set! should return ?
		return env.eval_set(expr...)
	case "lambda":
		return env.eval_lambda(expr...)
	default:
		panic("Unknown builtin")
	}
}

func (env *Env) eval_args(args ...Any) List {
	r := make(List, 0)
	for _, elem := range args {
		r = append(r, env.Eval(elem))
	}
	return r
}

func (env *Env) eval_expr(expr List) Any {
	head := expr[0]
	tail := expr[1:]
	f_value := env.Eval(head)
	f := to_function(f_value)
	args := env.eval_args(tail...)
	return f(args...)
}

func (env *Env) eval_list(expr List) Any {
	if len(expr) == 0 {
		return expr
	}
	head := expr[0]

	switch s := head.(type) {
	case Builtin:
		return env.eval_builtin(s, expr)
	default:
		return env.eval_expr(expr)
	}
}

func (env *Env) Eval(expr Any) Any {
	switch v := expr.(type) {
	case List:
		return env.eval_list(v)
	case Symbol:
		// Symbol atom is a name of object in the environment
		return env.symbol_lookup(v)
	default:
		// Other atoms are const literals
		return v
	}
}

func Lambda(s string) PureFunction {
	return to_function(StdEnv().Eval(ParseExpr(s)))
}

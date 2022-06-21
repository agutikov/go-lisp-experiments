package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Bool bool
type Int int64
type Str string
type Float float64

type Symbol string
type Builtin string
type Any interface{}
type List []Any

type PureFunction = func(List) Any
type EnvFunction = func(*Env, List) Any

func to_symbol(s Any) Symbol {
	switch v := s.(type) {
	case Symbol:
		return v
	default:
		panic("Invalid symbol")
	}
}

func to_list(s Any) List {
	switch v := s.(type) {
	case List:
		return v
	default:
		panic("Invalid list")
	}
}

func to_function(s Any) PureFunction {
	switch v := s.(type) {
	case PureFunction:
		return v
	default:
		fmt.Println(reflect.TypeOf(v))
		panic("Invalid function, ")
	}
}

func lispstr(expr Any) string {
	if expr == nil {
		return "nil"
	}
	switch v := expr.(type) {
	case Builtin:
		return string(v)
	case Symbol:
		return string(v)
	case Int:
		return fmt.Sprintf("%d", v)
	case Float:
		return fmt.Sprintf("%f", v)
	case Str:
		return "\"" + string(v) + "\""
	case Bool:
		if v {
			return "t"
		} else {
			return "f"
		}
	case List:
		s := "("
		for i, item := range v {
			if i != 0 {
				s += " "
			}
			s += lispstr(item)
		}
		s += ")"
		return s
	case PureFunction:
		return fmt.Sprintf("function{%p}", v)
	default:
		return "Unknown"
	}
}

type Env struct {
	parent        *Env
	named_objects map[Symbol]Any
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

func if_test(value Any) bool {
	//TODO: not sure about conversion of other types to bool
	// maybe everything except nil and false is true
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case Bool:
		return bool(v)
	case List:
		return len(v) > 0
	case Int:
		return int(v) != 0
	case Str:
		return len(v) > 0
	case Float:
		return float64(v) != 0.0
	default:
		panic("Non-data if test argument")
	}
}

func (env *Env) eval_if(expr List) Any {
	//TODO: verify number of args on parsing stage
	if len(expr) != 4 {
		panic("'if' statement requires exactly 3 arguments (if test conseq alt), while provided: " + lispstr(expr))
	}
	test := expr[1]
	conseq := expr[2]
	alt := expr[3]

	v := env.eval(test)

	if if_test(v) {
		return conseq
	} else {
		return alt
	}
}

func (env *Env) eval_define(expr List) Any {
	//TODO: verify number of args on parsing stage
	if len(expr) != 3 {
		panic("'define' statement requires exactly 2 arguments (define name exp), while provided: " + lispstr(expr))
	}
	name := expr[1]
	exp := expr[2]

	value := env.eval(exp)

	switch s := name.(type) {
	case Symbol:
		env.named_objects[s] = value
	default:
		panic("Invalid define name argument")
	}

	return nil
}

func (env *Env) eval_set(expr List) Any {
	//TODO: verify number of args on parsing stage
	if len(expr) != 3 {
		panic("'set!' statement requires exactly 2 arguments (set! name exp), while provided: " + lispstr(expr))
	}
	name := expr[1]
	exp := expr[2]

	value := env.eval(exp)

	switch s := name.(type) {
	case Symbol:
		env.env_lookup(s).named_objects[s] = value
	default:
		panic("Invalid set! name argument")
	}

	return nil
}

func lambda_args(vars Any) []Symbol {
	args := to_list(vars)

	a := []Symbol{}

	for _, elem := range args {
		a = append(a, to_symbol(elem))
	}

	return a
}

func (env *Env) eval_lambda(expr List) Any {
	//TODO: verify number of args on parsing stage
	if len(expr) != 3 {
		panic("'lambda' statement requires exactly 2 arguments (lambda (vars...) body), while provided: " + lispstr(expr))
	}
	p := expr[1]
	body := expr[2]

	params := lambda_args(p)

	// Return callable which will
	return func(args List) Any {
		// Eval body in the new nested environment
		e := newEnv(env)
		e.assign_vars(params, args)
		return e.eval(body)
	}
}

func eval_quote(expr List) Any {
	if len(expr) != 2 {
		panic("'quote' statement requires exactly 1 argument (quote exp), while provided: " + lispstr(expr))
	}
	return expr[1:]
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

func (env *Env) eval_builtin(s Builtin, expr List) Any {
	switch s { //TODO: is there any benefit of using map[Symbol]EnvFunction ?
	case "quote":
		//TODO: unquote????
		return eval_quote(expr)
	case "if":
		return env.eval_if(expr)
	case "define":
		//TODO: what define should return ?
		return env.eval_define(expr)
	case "set!":
		//TODO: what set! should return ?
		return env.eval_set(expr)
	case "lambda":
		return env.eval_lambda(expr)
	case "print-env":
		return env.print()
	default:
		panic("Unknown builtin")
	}
}

func (env *Env) eval_args(args List) List {
	r := make(List, 0)
	for _, elem := range args {
		r = append(r, env.eval(elem))
	}
	return r
}

func (env *Env) eval_expr(expr List) Any {
	head := expr[0]
	tail := expr[1:]
	f_value := env.eval(head)
	f := to_function(f_value)
	args := env.eval_args(tail)
	return f(args)
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

func (env *Env) eval(expr Any) Any {
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

func car(args List) Any {
	l := to_list(args[0])
	if len(l) == 0 {
		return nil
	}
	return l[0]
}

func cdr(args List) Any {
	return to_list(args[0])[1:]
}

func cons(args List) Any {
	if len(args) != 2 {
		panic("'cons' function requires exactly 2 arguments (const head exp), while provided: ")
	}
	var r List
	tail := List{}
	if args[1] != nil {
		tail = to_list(args[1])
	}
	r = append(List{args[0]}, tail...)
	return r
}

func list(args List) Any {
	return args
}

func standard_env() *Env {
	env := Env{}

	env.named_objects = map[Symbol]Any{
		"car":  car,
		"cdr":  cdr,
		"cons": cons,
		"list": list,
	}

	return &env
}

func tokenize(s string) []string {
	s = strings.Replace(s, "(", " ( ", -1)
	s = strings.Replace(s, ")", " ) ", -1)
	return strings.Fields(s)
}

type Parser struct {
	tokens []string
	pos    int
}

func newParser(s string) *Parser {
	p := Parser{pos: 0}
	p.tokens = tokenize(s)
	return &p
}

var builtins = map[string]Any{
	"if":        Builtin("if"),
	"quote":     Builtin("quote"),
	"define":    Builtin("define"),
	"set!":      Builtin("set!"),
	"lambda":    Builtin("lambda"),
	"t":         Bool(true),
	"f":         Bool(false),
	"nil":       nil,
	"print-env": Builtin("print-env"),
}

func (p *Parser) parse_atom(token string) Any {
	if v, ok := builtins[token]; ok {
		return v
	}

	int_val, err := strconv.Atoi(token)
	if err == nil {
		return Int(int_val)
	}

	float_val, err := strconv.ParseFloat(token, 64)
	if err == nil {
		return Float(float_val)
	}

	//TODO: strings - with parser generator
	return Symbol(token)
}

func (p *Parser) parse_list() Any {
	token := p.tokens[p.pos]
	p.pos++
	if token == "(" {
		l := List{}
		for p.tokens[p.pos] != ")" {
			l = append(l, p.parse_list())
		}
		p.pos++
		return l
	} else if token == ")" {
		panic("Unexpected ')'")
	} else {
		return p.parse_atom(token)
	}
}

func exec(env *Env, line string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	expr := newParser(line).parse_list()
	r := env.eval(expr)
	fmt.Println(lispstr(r))
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	env := standard_env()
	for {
		fmt.Print("lisp> ")
		line, _ := reader.ReadString('\n')
		exec(env, line)
	}
}

func test() {
	env := standard_env()

	l2 := List{Symbol("cons"), Symbol("list"), nil}
	fmt.Println(l2)
	r2 := env.eval(l2)
	fmt.Println(r2)

	l3 := List{Symbol("cons"), Int(42), l2}
	fmt.Println(l3)
	r3 := env.eval(l3)
	fmt.Println(r3)

	l4 := newParser("(car (cdr (list 1 2.01 3)))").parse_list()
	fmt.Println(l4)
	r4 := env.eval(l4)
	fmt.Println(r4)
}

func main() {
	repl()
}

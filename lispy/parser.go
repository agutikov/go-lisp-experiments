package lispy

import (
	"strconv"
	"strings"
)

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
	"if":     Builtin("if"),
	"quote":  Builtin("quote"),
	"define": Builtin("define"),
	"set!":   Builtin("set!"),
	"lambda": Builtin("lambda"),
	"t":      Bool(true),
	"f":      Bool(false),
	"nil":    nil,
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

func ParseExpr(s string) Any {
	return newParser(s).parse_list()
}

package lispy

import (
	"math/big"
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

	n := new(big.Int)
	n, ok := n.SetString(token, 10)

	if ok {
		return Int{n}
	}

	f := new(big.Rat)
	f, ok = f.SetString(token)

	if ok {
		return Float{f}
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

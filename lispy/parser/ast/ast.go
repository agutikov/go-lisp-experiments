package ast

import (
	"strings"

	"github.com/agutikov/go-lisp-experiments/lispy/parser/token"
)

type Attrib interface{}

type Atom struct {
	Value string
}

type Sexpr struct {
	List []Atom
}

func NewAtom(id Attrib) (*Atom, error) {
	return &Atom{string(id.(*token.Token).Lit)}, nil
}

func (this *Atom) String() string {
	return this.Value
}

func NewSexpr(head Attrib, tail Attrib) (*Sexpr, error) {
	atom := tail.(*Atom)
	if head == nil {
		return &Sexpr{List: []Atom{*atom}}, nil
	}
	s := head.(*Sexpr)
	s.List = append(s.List, *atom)
	return s, nil
}

func Map[From any, To any](f func(From) To, args []From) []To {
	r := []To{}
	for _, arg := range args {
		r = append(r, f(arg))
	}
	return r
}

func (this *Sexpr) String() string {
	return strings.Join(Map(func(a Atom) string { return a.String() }, this.List), " ")
}

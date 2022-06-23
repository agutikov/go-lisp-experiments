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
	Value interface{}
}

type List struct {
	List []Sexpr
}

type Sequence struct {
	List []Sexpr
}

func NewAtom(id Attrib) (*Atom, error) {
	return &Atom{string(id.(*token.Token).Lit)}, nil
}

func NewSequence(s Attrib) (*Sequence, error) {
	sexpr := *s.(*Sexpr)
	return &Sequence{[]Sexpr{sexpr}}, nil
}

func Cons(car Attrib, cdr Attrib) (*Sequence, error) {
	sexpr := *car.(*Sexpr)
	seq := *cdr.(*Sequence)

	seq.List = append([]Sexpr{sexpr}, seq.List...)

	return &seq, nil
}

func NewList(seq Attrib) (*List, error) {
	lst := List{}
	if seq != nil {
		lst.List = seq.(*Sequence).List
	}
	return &lst, nil
}

func NewSexpr(s Attrib) (*Sexpr, error) {
	switch s := s.(type) {
	case *Atom:
		return &Sexpr{*s}, nil
	case *List:
		return &Sexpr{*s}, nil
	default:
		panic("Invalid NewSexpr() argument")
	}
}

func Map[From any, To any](f func(From) To, args []From) []To {
	r := []To{}
	for _, arg := range args {
		r = append(r, f(arg))
	}
	return r
}

func (this *Atom) String() string {
	return this.Value
}

func (this *List) String() string {
	return "( " + strings.Join(Map(func(a Sexpr) string { return a.String() }, this.List), " ") + " )"
}

func (this *Sequence) String() string {
	return strings.Join(Map(func(a Sexpr) string { return a.String() }, this.List), " ")
}

func (this *Sexpr) String() string {
	switch s := this.Value.(type) {
	case Atom:
		return s.String()
	case List:
		return s.String()
	default:
		panic("Invalid Sexpr")
	}
}

package ast

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/token"
)

type Attrib interface{}

type Nil struct{}

type Bool bool

type Symbol struct {
	Name string
}

type Number struct {
	Value int
}

type Str struct {
	Value string
}

type Sexpr interface{}

type Quote struct {
	Value Sexpr
}

type Unquote struct {
	Value Sexpr
}

type List struct {
	List []Sexpr
}

type Sequence struct {
	Seq []Sexpr
}

type If struct {
	Test      Sexpr
	PosBranch Sexpr
	NegBranch Sexpr
}

type Define struct {
	Sym   Symbol
	Value Sexpr
}

type Set struct {
	Sym   Symbol
	Value Sexpr
}

type Lambda struct {
	Args List
	Body Sexpr
}

func NewSymbol(t Attrib) (Symbol, error) {
	name := string(t.(*token.Token).Lit)
	return Symbol{name}, nil
}

func NewNumber(t Attrib) (Number, error) {
	s := string(t.(*token.Token).Lit)
	val, err := strconv.Atoi(s)
	if err != nil {
		return Number{}, errors.New("invalid number literal: \"" + s + "\"")
	}
	return Number{val}, nil
}

func str_replace_escaped(s string) string {
	r := strings.NewReplacer(
		"\\\\", "\\",
		"\\\"", "\"",
		"\\n", "\n",
		"\\t", "\t",
	) //TODO: Should be already a way to "compile" a string
	return r.Replace(s)
}

func NewStr(t Attrib) (Str, error) {
	s := string(t.(*token.Token).Lit)
	unquoted := s[1 : len(s)-1]
	compiled := str_replace_escaped(unquoted)
	return Str{compiled}, nil
}

func NewSequence(s Attrib) (Sequence, error) {
	sexpr := s.(Sexpr)
	return Sequence{[]Sexpr{sexpr}}, nil
}

func Cons(car Attrib, cdr Attrib) (Sequence, error) {
	sexpr := car.(Sexpr)
	seq := cdr.(Sequence)

	seq.Seq = append([]Sexpr{sexpr}, seq.Seq...)

	return seq, nil
}

func NewList(seq Attrib) (List, error) {
	lst := List{}
	if seq != nil {
		lst.List = seq.(Sequence).Seq
	}
	return lst, nil
}

func NewQuote(sexpr Attrib) (Quote, error) {
	return Quote{sexpr.(Sexpr)}, nil
}

func NewUnquote(sexpr Attrib) (Unquote, error) {
	return Unquote{sexpr.(Sexpr)}, nil
}

func NewDefine(symbol Attrib, value Attrib) (Define, error) {
	return Define{Sym: symbol.(Symbol), Value: value.(Sexpr)}, nil
}

func NewSet(symbol Attrib, value Attrib) (Set, error) {
	return Set{Sym: symbol.(Symbol), Value: value.(Sexpr)}, nil
}

func NewIf(test Attrib, pos_branch Attrib, neg_branch Attrib) (If, error) {
	return If{
		Test:      test.(Sexpr),
		PosBranch: pos_branch.(Sexpr),
		NegBranch: neg_branch.(Sexpr),
	}, nil
}

func NewLambda(args Attrib, body Attrib) (Lambda, error) {
	return Lambda{Args: args.(List), Body: body.(Sexpr)}, nil
}

func Map[From any, To any](f func(From) To, args []From) []To {
	r := []To{}
	for _, arg := range args {
		r = append(r, f(arg))
	}
	return r
}

func (this Str) String() string {
	return fmt.Sprintf("%q", this.Value)
}

func (this Number) String() string {
	return fmt.Sprintf("%d", this.Value)
}

func (this Symbol) String() string {
	return this.Name
}

func (this List) String() string {
	return "(" + strings.Join(Map(func(a Sexpr) string { return String(a) }, this.List), " ") + ")"
}

func (this Sequence) String() string {
	return strings.Join(Map(func(a Sexpr) string { return String(a) }, this.Seq), "\n")
}

func String(this Sexpr) string {
	return fmt.Sprintf("%+v", this)
}

func (this Quote) String() string {
	return "'" + String(this.Value)
}

func (this Unquote) String() string {
	return "," + String(this.Value)
}

func (this Nil) String() string {
	return "nil"
}

func (this Bool) String() string {
	if bool(this) {
		return "t"
	} else {
		return "f"
	}
}

func (this If) String() string {
	return fmt.Sprintf("(if %+v %+v %+v)", this.Test, this.PosBranch, this.NegBranch)
}

func (this Define) String() string {
	return fmt.Sprintf("(define %+v %+v)", this.Sym, this.Value)
}

func (this Set) String() string {
	return fmt.Sprintf("(set! %+v %+v)", this.Sym, this.Value)
}

func (this Lambda) String() string {
	return fmt.Sprintf("(lambda %+v %+v)", this.Args, this.Body)
}

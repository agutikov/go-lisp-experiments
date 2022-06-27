package ast

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/agutikov/go-lisp-experiments/lispy/syntax/token"
)

type Attrib interface{}

type Nil struct{} //TODO: replace with nil

type Bool bool

type Symbol struct {
	Name string
}

type Number struct {
	//TODO: number with dynamic conversion
}

type Int struct {
	Value *big.Int
}

type Float struct {
	Value *big.Rat
}

type Str struct {
	Value string
}

type Any interface{}

type Quote struct {
	Value Any
}

type Unquote struct {
	Value Any
}

type List []Any

type Sequence []Any

type If struct {
	Test      Any
	PosBranch Any
	NegBranch Any
}

type Define struct {
	Sym   Symbol
	Value Any
}

type Set struct {
	Sym   Symbol
	Value Any
}

type Lambda struct {
	Args []Symbol
	Body Any
}

type Defun struct {
	Sym Symbol
	L Lambda
}

func NewSymbol(t Attrib) (Symbol, error) {
	name := string(t.(*token.Token).Lit)
	return Symbol{name}, nil
}

func NewInt(t Attrib) (Int, error) {
	s := string(t.(*token.Token).Lit)

	n := new(big.Int)
	n, ok := n.SetString(s, 10)
	if !ok {
		return Int{nil}, errors.New("invalid Int literal: \"" + s + "\"")
	}

	return Int{n}, nil
}

func IntNum(i int64) Int {
	return Int{big.NewInt(i)}
}

func FloatNum(f float64) Float {
	r := new(big.Rat)
	r.SetFloat64(f)
	return Float{r}
}

func NewFloat(t Attrib) (Float, error) {
	s := string(t.(*token.Token).Lit)

	f := new(big.Rat)
	n, ok := f.SetString(s)
	if !ok {
		return Float{nil}, errors.New("invalid Float literal: \"" + s + "\"")
	}

	return Float{n}, nil
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
	sexpr := s.(Any)
	return Sequence{sexpr}, nil
}

func Cons(car Attrib, cdr Attrib) (Sequence, error) {
	sexpr := car.(Any)
	seq := cdr.(Sequence)

	seq = append([]Any{sexpr}, seq...)

	return seq, nil
}

func NewList(seq Attrib) (List, error) {
	lst := List{}
	if seq != nil {
		lst = List(seq.(Sequence))
	}
	return lst, nil
}

func NewQuote(sexpr Attrib) (Quote, error) {
	return Quote{sexpr.(Any)}, nil
}

func NewUnquote(sexpr Attrib) (Unquote, error) {
	return Unquote{sexpr.(Any)}, nil
}

func NewDefine(symbol Attrib, value Attrib) (Define, error) {
	return Define{Sym: symbol.(Symbol), Value: value.(Any)}, nil
}

func NewSet(symbol Attrib, value Attrib) (Set, error) {
	return Set{Sym: symbol.(Symbol), Value: value.(Any)}, nil
}

func NewIf(test Attrib, pos_branch Attrib, neg_branch Attrib) (If, error) {
	return If{
		Test:      test.(Any),
		PosBranch: pos_branch.(Any),
		NegBranch: neg_branch.(Any),
	}, nil
}

func NewLambda(args Attrib, body Attrib) (Lambda, error) {
	a := []Symbol{}
	for _, item := range args.(Sequence) {
		a = append(a, item.(Symbol))
	}
	return Lambda{Args: a, Body: body.(Any)}, nil
}

func NewDefun(sym Attrib, args Attrib, body Attrib) (Defun, error) {
	a := []Symbol{}
	for _, item := range args.(Sequence) {
		a = append(a, item.(Symbol))
	}
	return Defun{Sym: sym.(Symbol), L: Lambda{Args: a, Body: body.(Any)}}, nil
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

func (this Int) String() string {
	return this.Value.String()
}

func (this Float) String() string {
	f := big.NewFloat(0)
	f.SetPrec(0)
	f.SetRat(this.Value)
	return f.Text('f', int(f.MinPrec())) //TODO: get precision from env
}

func (this Symbol) String() string {
	return this.Name
}

func (this List) String() string {
	return "(" + strings.Join(Map(func(a Any) string { return String(a) }, this), " ") + ")"
}

func (this Sequence) String() string {
	return strings.Join(Map(func(a Any) string { return String(a) }, this), "\n")
}

func String(this Any) string {
	if this == nil {
		return "nil"
	}
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
		return "false"
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

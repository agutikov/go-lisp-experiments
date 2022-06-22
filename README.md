# Experimenting with Lisp implementation in Go

Basically this is [Peter Norvig's Lis.py](http://norvig.com/lispy.html) translated to Go.
But with some limitations induced by the fact that Python is already a language with dynamic types,
while in Go we need to implement dynamic typing manually.

Go-lispy is a subset of Sheme, with following implemented:
* atoms, booleans, integer and float numbers
* special forms (keywords): cons, if, define, set!, lambda
* functions:
  * list functions: **car, cdr, cons, list, length**
  * arithmetic: +, -, *, /
  * comparison:
    * arithmetic: **>, <, >=, <=**
    * equality for all types: **=**
  * boolean functions: **not, and, or**
  * functional: **apply, map**
  * other: **begin**

Go-lispy implements lists with slices - so there is no dotted pairs like in classic Lisp.


## HOWTO

```sh
# compile
$ cd go-lisp-experiments/
$ make

# run tests and benchmark
$ make test

# run repl, Ctrl+D to exit
$ ./go-lispy
go-lis.py>

# eval command line arguments
$ ./go-lispy '(begin (define r 10)'     '(* pi (* r r)))'
314.159265

```

## How to use it as a library

```Go
package main

import (
    "fmt"
    "github.com/agutikov/go-lisp-experiments/lispy"
)

func main() {
    // create env
    env := lispy.StdEnv()

    // parse the expression string
    expr := ParseExpr("(car (cdr (cons 1 (list 2 3 4))))")

    // eval the expr
    r := env.Eval(expr)

    // print result
    fmt.Println(LispyStr(r))
}
```

### Embed the Lispy lambda into the Go code

```Go
// Get an executable from lambda expression
// NOTE: that go-lispy doesn't have long integers - so result would be float
fact := lispy.Lambda("(define fact (lambda (n) (if (<= n 1) 1.0 (* n (fact (- n 1))))))")

// Call the function
// NOTE: if go-lispy interpreter interacts with values - then lispy types should be used
v := fact(lispy.Int(100))
fmt.Println(v)

// If go-lispy interpreter will not interact with values - then any types could be used
zip2 := lispy.Lambda("(lambda (slice_1 slice_2) (map list slice_1 slice_2))")
a := lispy.List{0, 1, 2}
b := lispy.List{"str", true}
r := zip2(a, b)
fmt.Println(r)

```



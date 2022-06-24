# Experimenting with Lisp implementation in Go

Basically this is [Peter Norvig's Lis.py](http://norvig.com/lispy.html) translated into Go.
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

It is not so compact as original *lis.py*, but a little bit faster: go-lispy computes (fact 100) in 0.6ms, while original *lis.py* takes 3ms.

Parsing is implemented with [Gocc](https://github.com/goccmack/gocc)

Syntax definition: [lispy/syntax/lispy.bnf](lispy/syntax/lispy.bnf)


## HOWTO

```sh
# dependencies
$ go get github.com/goccmack/gocc
$ go install github.com/goccmack/gocc

# compile
$ cd go-lisp-experiments/
$ make

# run tests and benchmark
$ make test

# run repl, Ctrl+D to exit
$ ./go-lispy
go-lis.py>

# eval command line arguments
$ ./go-lispy -e '(begin (define r 10) (* pi (* r r)))'
314.159265

# eval file
$ ./go-lispy -e '(set! enable-print-elapsed t) (set! enable-trace t)' ./fact-bench.lsp
$ ./go-lispy -e '(set! enable-print-elapsed t) (set! enable-trace t)' ./lispy-test.lsp

```

## Extra features (in addition to original lis.py)

#### Quote and unquote

```
go-lis.py> '(1 ,(- 0 1) 2)
'(1 -1 2)

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

#### Embed the Lispy lambda into the Go code

```Go
// Get an executable from lambda expression
fact := lispy.Lambda("(define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1))))))")

// Call the function
// NOTE: if go-lispy interpreter interacts with values - then lispy types should be used
v := fact(lispy.FromInt(100))
fmt.Println(LispyStr(v))

// If go-lispy interpreter will not interact with values - then any types could be used
zip2 := lispy.Lambda("(lambda (slice_1 slice_2) (map list slice_1 slice_2))")
a := lispy.List{0, 1, 2}
b := lispy.List{"str", true}
r := zip2(a, b)
fmt.Println(r)

```

## Go-Lispy features

#### Big numbers

```
go-lis.py> (define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1))))))
function{0x49b740}
go-lis.py> (fact 1000)
402387260077093773543702433923003985719374864210714632543799910429938512398629020592044208486969404800479988610197196058631666872994808558901323829669944590997424504087073759918823627727188732519779505950995276120874975462497043601418278094646496291056393887437886487337119181045825783647849977012476632889835955735432513185323958463075557409114262417474349347553428646576611667797396668820291207379143853719588249808126867838374559731746136085379534524221586593201928090878297308431392844403281231558611036976801357304216168747609675871348312025478589320767169132448426236131412508780208000261683151027341827977704784635868170164365024153691398281264810213092761244896359928705114964975419909342221566832572080821333186116811553615836546984046708975602900950537616475847728421889679646244945160765353408198901385442487984959953319101723355556602139450399736280750137837615307127761926849034352625200015888535147331611702103968175921510907788019393178114194545257223865541461062892187960223838971476088506276862967146674697562911234082439208160153780889893964518263243671616762179168909779911903754031274622289988005195444414282012187361745992642956581746628302955570299024324153181617210465832036786906117260158783520751516284225540265170483304226143974286933061690897968482590125458327168226458066526769958652682272807075781391858178889652208164348344825993266043367660176999612831860788386150279465955131156552036093988180612138558600301435694527224206344631797460594682573103790084024432438465657245014402821885252470935190620929023136493273497565513958720559654228749774011413346962715422845862377387538230483865688976461927383814900140767310446640259899490222221765904339901886018566526485061799702356193897017860040811889729918311021171229845901641921068884387121855646124960798722908519296819372388642614839657382291123125024186649353143970137428531926649875337218940694281434118520158014123344828015051399694290153483077644569099073152433278288269864602789864321139083506217095002597389863554277196742822248757586765752344220207573630569498825087968928162753848863396909959826280956121450994871701244516461260379029309120889086942028510640182154399457156805941872748998094254742173582401063677404595741785160829230135358081840096996372524230560855903700624271243416909004153690105933983835777939410970027753472000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
```


# What I've learned about Lisp

1. Syntactic form differs from procedure call in a way that different syntactic forms can eval or not eval some of it's arguments, while procedure call (as a syntactic form itself) eval all arguments before calling the procedure.
For example:
  - procedure call - eval all args before call;
  - quote - does not eval any arg;
  - define and set! - does not eval first, eval second;
  - if - eval first, and then eval one of second or third;
  - lambda - does not eval any of two, but will eval body when been called, or can eval body partially, or behave any different way;

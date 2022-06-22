# Experimenting with Lisp implementation in Go

Basically this is [Peter Norvig's Lis.py](http://norvig.com/lispy.html) implemented in Go.


## HOWTO

```sh
# compile
$ cd go-lisp-experiments/
$ make

# run repl, Ctrl+D to exit
$ ./go-lispy
go-lis.py>

# eval in command line
$ ./go-lispy '(begin (define r 10) (* pi (* r r)))'
314.159265

```







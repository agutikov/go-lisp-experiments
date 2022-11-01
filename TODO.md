


- lambda body pre-eval and Go lambda composition





- MACRO ?????? let ????

- casting SPELs in Go-lispy (https://www.lisperati.com/casting.html)



- std:
  - math.go, Number, etc...
  - list.go
  - functional.go : apply, map
  - bool.go
  - common.go : equality, begin, ...
  - macro.go : ??? builtin macros

- debug tools
  - verbose
    - top level expressions
    - whole tree with indentation and levels coloring
  - step-by-step eval
  - history
  - inspect env, export complete env or user's symbols


- eval string :))))))
- sprintf, split, join




- call Go functions and methods by reflect.ValueOf and MethodByName

- goroutines, parallellism

- limit the available symbol set for evaluation of expr only in terms of DSL


- type inference
  - procedures have description: num of args, types, etc... so apply could verify it
  - no needs to do checks inside the functions
  - use apply() in eval_funcation_call() after car been evaluated


- S-expressions, cons, dot pairs, cons cells
  - parse into S-expressions
  - compare eval performance for slices and dot pairs
  - then "compile" into slices with pre-evaluation of const expressions
  - then "compile" with lambda composition
  - then call final lambda

- monads ??? in lisp ???




===================================



## Go sum type

https://medium.com/@haya14busa/sum-union-variant-type-in-go-and-static-check-tool-of-switch-case-handling-3bfc61618b1e

https://making.pusher.com/alternatives-to-sum-types-in-go/


## Go project setup

https://blog.boot.dev/golang/golang-project-structure/

https://www.wolfe.id.au/2020/03/10/how-do-i-structure-my-go-project/

https://dev.to/jinxankit/go-project-structure-and-guidelines-4ccm

https://github.com/golang-standards/project-layout


## S-expressions and Cons

https://en.wikipedia.org/wiki/S-expression

https://en.wikipedia.org/wiki/Cons


## Lisp implementaions in C

https://habr.com/ru/post/150805/

https://github.com/rui314/minilisp

https://carld.github.io/2017/06/20/lisp-in-less-than-200-lines-of-c.html

https://www.buildyourownlisp.com/contents


## Lisp implementations in Go

https://github.com/janne/go-lisp

https://github.com/nukata/lisp-in-go

https://github.com/chenzhuoyu/simple-lisp

https://github.com/amirgamil/lispy


================================================================================


Command Line Lexer - the simplest possible tool for command line args.

Gets all flags in a map and rest positionals in a slice.
-a -b : bool flags
-f xxx : with args - ints or strings (all positionals are strings)
-- x -file --y : only positional args after double dash
- : single minus also a flag, not positional
-x 1 -x 2 : multiple entries will became a slice
  if you want do --files x y z - just, ... go to hell
-g=3 : equal sign is allowed

With int levels=n option set to n > 0 will treat positionals as commands and subcommands
and arrange them in a linked-list starting from globals.

It doesn't have any assumptions about required or optional arguments
because it is not a config management tool.
If someone wants - can make a validator on top of this structure
and get a full-featured CLI library.
But most frequently it's not necessary.

Next - lvl 2 - would be from config structure auto-generated set of options and arguments.

And next - last level - would be a validator of the config with all possible ways of validation:
- options format:
  - data type: string, number, bool, ....
  - available values
  - string regex
  - validation callback
- external environment dependencies:
  - files: must/may exist or must/may not exist
    - input, output, rewrite protection, rotation, etc...
- dependencies between options:
  - presence
  - values
  - format
  - predicates on multiple deps
  - ... and anything else you may want



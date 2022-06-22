package main

import (
	"bufio"
	"example/user/golisp/lispy"
	"fmt"
	"os"
	"strings"
)

func exec(env *lispy.Env, line string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	expr := lispy.newParser(line).parse_list()
	r := env.eval(expr)
	fmt.Println(lispstr(r))
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	env := standard_env()
	for {
		fmt.Print("lisp> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		exec(env, line)
	}
}

func main() {
	if len(os.Args) > 1 {
		env := standard_env()
		exec(env, strings.Join(os.Args[1:], " "))
	} else {
		repl()
	}
}

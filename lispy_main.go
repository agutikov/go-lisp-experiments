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

	expr := lispy.NewParser(line).ParseList()
	r := env.Eval(expr)
	fmt.Println(lispy.LispyStr(r))
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	env := lispy.StdEnv()
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
		env := lispy.StdEnv()
		exec(env, strings.Join(os.Args[1:], " "))
	} else {
		repl()
	}
}

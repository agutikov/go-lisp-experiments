package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/agutikov/go-lisp-experiments/lispy"
)

func exec(env *lispy.Env, line string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	expr := lispy.ParseExpr(line)
	r := env.Eval(expr)
	fmt.Println(lispy.LispyStr(r))
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	env := lispy.StdEnv()
	for {
		fmt.Print("go-lis.py> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			break
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			exec(env, line)
		} else {
			fmt.Println()
		}
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

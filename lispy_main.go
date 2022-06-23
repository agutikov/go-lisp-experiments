package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/agutikov/go-lisp-experiments/lispy"
)

func exec(env *lispy.Env, line string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	started := time.Now()
	expr := lispy.ParseExpr(line)
	r := env.Eval(expr)
	elapsed := time.Since(started)
	fmt.Println(lispy.LispyStr(r))
	fmt.Println(" elapsed:", elapsed)
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

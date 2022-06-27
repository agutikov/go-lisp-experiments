package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/agutikov/go-lisp-experiments/lispy"
	"github.com/agutikov/go-lisp-experiments/cmdlex"
)

func exec(env *lispy.Env, line string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	expr := lispy.ParseStr(line)
	r := env.Eval(expr)
	fmt.Println(lispy.LispyStr(r))
}

func exec_file(env *lispy.Env, filename string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	expr := lispy.ParseFile(filename)
	r := env.Eval(expr)
	fmt.Println(lispy.LispyStr(r))
}

func repl(env *lispy.Env) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("go-lis.py> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			exec(env, line)
		}
	}
}



func main() {
	args := cmdlex.ParseCmdLineArgs(os.Args, 0)

	env := lispy.StdEnv()

	if exprs, ok := args.Options["e"]; ok {
		// Eval cmdline args
		exec(env, strings.Join(exprs, "\n"))
	}

	for _, filename := range args.Positional[1:] {
		if filename == "-" {
			repl(env)
		} else {
			exec_file(env, filename)
		}
	}

	if len(args.Positional) == 1 && len(args.Options) == 0 {
		// No other options - run repl
		repl(env)
	}
}

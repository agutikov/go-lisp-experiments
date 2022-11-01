package main

import (
	"fmt"
	"math/big"

	"github.com/agutikov/go-lisp-experiments/lispy"
)

func example() {
	// create env
	env := lispy.StdEnv()

	// parse the expression string
	expr := lispy.ParseStr("(car (cdr (cons 1 (list 2 3 4))))")

	// eval the expr
	r := env.Eval(expr)

	// print result
	fmt.Println(lispy.LispyStr(r))
}

func fact_bench() {
	fact := lispy.Function("(defun fact (n) (if (<= n 1) 1 (* n (fact (- n 1)))))")

	// Call the function
	// NOTE: if go-lispy interpreter interacts with values - then lispy types should be used
	args := []lispy.Any{lispy.Int{Value: big.NewInt(1000)}}
	v := fact(args)

	fmt.Println(lispy.LispyStr(v))
}

func native_fact_r_bench() {
	fact := lispy.Function("__fact_r")

	// Call the function
	// NOTE: if go-lispy interpreter interacts with values - then lispy types should be used
	args := []lispy.Any{lispy.Int{Value: big.NewInt(1000)}}
	v := fact(args)

	fmt.Println(lispy.LispyStr(v))
}

func main() {
	native_fact_r_bench()
}

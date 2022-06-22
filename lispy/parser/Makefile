

# the order is important, because of dependencies
all:
	gocc ./lispy.bnf
	go build ./util
	go build ./token
	go build ./ast
	go build ./errors
	go build ./lexer
	go build ./parser

test:
	go test -v
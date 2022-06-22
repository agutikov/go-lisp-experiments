

all: go-lispy

GO_LISPY_SRCS := lispy_main.go

BINARIES := go-lispy go-lispy-debug


go-lispy: ${GO_LISPY_SRCS}
	go build -o $@ $^

go-lispy-debug: $(GO_LISPY_SRCS)
	go build -gcflags "-N" -o $@ $^


.PHONY: test
test:
	go test ./...


.PHONY: clean
clean:
	go clean
	rm ${BINARIES}





all: go-lispy

GO_LISPY_SRCS := lispy_main.go

BINARIES := go-lispy go-lispy-debug


go-lispy: $(GO_LISPY_SRCS) | lispy
	go build -o $@ $^

go-lispy-debug: $(GO_LISPY_SRCS) | lispy
	go build -gcflags "-N" -o $@ $^


.PHONY: lispy
lispy:
	$(MAKE) -C $@

.PHONY: test
test:
	go test ./...


.PHONY: clean
clean:
	go clean
	rm $(BINARIES)





all: go-lispy

GO_LISPY_SRCS := lispy_main.go

BINARIES := go-lispy go-lispy-debug

# TODO: dependency from lispy package
go-lispy: $(GO_LISPY_SRCS) | lispy
	go build -o $@ $^

go-lispy-debug: $(GO_LISPY_SRCS) | lispy
	go build -gcflags "-N" -o $@ $^

.PHONY: run
run:
	go run $(GO_LISPY_SRCS)

.PHONY: lispy
lispy:
	$(MAKE) -C $@

.PHONY: test
test:
	go test -v ./...


.PHONY: clean
clean:
	go clean
	rm -f $(BINARIES)



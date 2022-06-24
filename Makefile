

all: go-lispy

GO_LISPY_SRCS := lispy_main.go

BINARIES := go-lispy go-lispy-debug

PKGS := lispy cmdlex

# TODO: dependency from lispy package
go-lispy: $(GO_LISPY_SRCS) | ${PKGS}
	go build -o $@ $^

go-lispy-debug: $(GO_LISPY_SRCS) | ${PKGS}
	go build -gcflags "-N" -o $@ $^

.PHONY: run
run:
	go run $(GO_LISPY_SRCS)

.PHONY: lispy cmdlex
${PKGS}:
	$(MAKE) -C $@

.PHONY: test
test:
	go test ./...


.PHONY: clean
clean:
	go clean
	rm -f $(BINARIES)



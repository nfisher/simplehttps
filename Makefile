# ex : shiftwidth=2 tabstop=2 softtabstop=2 :                                      

SHELL := /bin/sh
SRC := $(wildcard *.go) cmd/simplehttps/main.go
GIT_REV := $(shell git rev-parse --short HEAD)

.PHONY: all
all: lint.out vet.out coverage.out bench.out

.PHONY: bench
bench: bench.out

bench.out: $(SRC)
	go test -bench ./... | tee bench.out

cover.out: $(SRC)
	go test -v -cover -covermode atomic -coverprofile cover.out

coverage.html: cover.out
	go tool cover -html=cover.out -o coverage.html

coverage.out: cover.out
	go tool cover -func=cover.out | tee coverage.out

.PHONY: clean
clean:
	rm *.out
	go clean -i ./...

.PHONY: fast
fast: vet cov

lint.out: $(SRC)
	golint -set_exit_status | tee lint.out

.PHONY: test
test: coverage.out

vet.out: install
	go vet -composites=false -v ./... | tee vet.out

install: $(SRC)
	go install -v -ldflags "-X main.Version=$(GIT_REV)" github.com/nfisher/simplehttps/cmd/simplehttps

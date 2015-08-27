SHELL := /bin/sh
SRC := $(wildcard *.go)
EXE := simplehttps

.PHONY: all
all: test vet $(EXE)

$(EXE): $(SRC)
	go build

.PHONY: run
run: $(EXE)
	./$(EXE)

.PHONY: install
install:
	go install

.PHONY: test
test:
	go test -v

.PHONY: vet
vet:
	go vet -x

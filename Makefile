SHELL := /bin/sh
SRC := $(wildcard *.go)
EXE := simplehttp

$(EXE): $(SRC)
	go build

.PHONY: run
run: $(EXE)
	./$(EXE)

.PHONY: install
install:
	go install

SHELL := /bin/sh
SRC := $(wildcard *.go)
EXE := simplehttps
COV := coverage.out

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

.PHONY: cov
cov: $(COV)
	go tool cover -func=coverage.out

.PHONY: htmlcov
htmlcov: $(COV)
	go tool cover -html=coverage.out

$(COV): $(SRC)
	go test -v -covermode=count -coverprofile=coverage.out

.PHONY: vet
vet:
	go vet -x

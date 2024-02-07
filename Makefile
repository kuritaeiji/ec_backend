.DEFAULT_GOAL := build

tidy:
	go mod tidy
.PHONY: tidy

fmt: tidy
	goimports -l -w config share util enduser systemuser
.PHONY: fmt

lint: fmt
	golangci-lint run ./...
.PHONY: lint

gen: lint
	go generate ./...
.PHONY: gen

run: gen
	export ENV=dev && go run -race enduser/api/main.go
.PHONY: run

build: gen
	go build -o main enduser/api/main.go
.PHONY: build

test: gen
	export ENV=test && go test -v -cover -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html
.PHONY: test

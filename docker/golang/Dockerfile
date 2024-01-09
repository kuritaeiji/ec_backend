FROM golang:1.21.4-bullseye as dev

WORKDIR /go/app

RUN apt update
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/vektra/mockery/v2@v2.37.1

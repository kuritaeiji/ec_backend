# 開発環境
FROM golang:1.21.4-bullseye as dev

WORKDIR /go/app
ENV WD=/go/app

RUN apt update
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/vektra/mockery/v2@v2.37.1

# 以下本番環境用
# docker image build --target pro -t ec_backend:v1 -f ./docker/golang/Dockerfile .
FROM golang:1.21.4-bullseye as builder

WORKDIR /go/app

COPY . .
RUN make tidy
RUN make build

FROM ubuntu:22.04 as pro

RUN apt update && apt install -y tzdata ca-certificates

USER 1001

ENV ENV=pro
ENV TZ=Asia/Tokyo
ENV WD=/go/app

WORKDIR /go/app

COPY --from=builder /go/app/main .

RUN mkdir env
COPY --from=builder /go/app/env/pro.env ./env

CMD ["./main"]

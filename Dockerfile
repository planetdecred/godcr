FROM golang 

RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.19.1

ENV GO111MODULE on

RUN go version

ENTRYPOINT [ "./run_tests.sh", "lint"]
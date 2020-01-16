FROM golang 

RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.19.1

ENV GO111MODULE on

RUN go version

RUN apt-get update

RUN apt-get install libwayland-dev libx11-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev --yes

ENTRYPOINT [ "./run_tests.sh", "lint"]
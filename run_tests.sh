#!/bin/bash

# Usage: ./run_tests.sh [all | build | lint]

set -e

function build_and_test() {
    go build
    go test
    go vet
}

function lint() {
    golangci-lint run --deadline=10m \
      --disable-all \
      --enable govet \
      --enable staticcheck \
      --enable gosimple \
      --enable unconvert \
      --enable ineffassign \
      --enable structcheck \
      --enable goimports \
      --enable misspell \
      --enable unparam
}

option=$1
if [[ "$option" = "build" ]]; then
    build_and_test
elif [[ "$option" = "lint" ]]; then
    lint
elif [[ "$option" = "all" ]]; then
    build_and_test
    lint
else
    echo "Usage: ./run_tests.sh [all | build | lint]"
fi

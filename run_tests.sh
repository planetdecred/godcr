#!/bin/bash

# Usage: ./run_tests.sh [all | build | lint]

set -e

function build_and_test() {
    go build
    go test
    go vet
}

function lint() {
    golangci-lint run --deadline=10m
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

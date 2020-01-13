#!/bin/bash

# Usage: ./run_docker.sh [setup | build]

set -e

function setup() {
    if ! docker ps -a | grep -q buildenv ; then
        docker build -t golangci . 
        docker run --name buildenv -v "$(pwd)"/.:/src -w /src golangci go build
    else
        echo "buildenv already up"
    fi
}

function build() {
    docker start -ia buildenv
}

option=$1
if [[ "$option" = "setup" ]]; then
    setup
elif [[ "$option" = "build" ]]; then
    build
else
    echo "Usage: ./run_docker.sh [setup | build]"
fi

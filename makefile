# This how we want to name the binary output
BINARY=godcr

VERSION="0.0.1"
BUILD=`date -u +"%Y-%m-%dT%H:%M:%SZ"`
# dev or prod
BuildEnv="prod" 

LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.BuildDate=${BUILD} -X main.BuildEnv=${BuildEnv}"

export GOARCH=amd64

all: clean darwin windows

freebsd:
	GOOS=freebsd go build -trimpath ${LDFLAGS} -o ${BINARY}-freebsd-${GOARCH}

darwin:
	GOOS=darwin go build -trimpath ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH}

windows:
	GOOS=windows go build -trimpath ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe
 
# Cleans our project: deletes old binaries
clean:
	-rm -f ${BINARY}-*

.PHONY: clean darwin windows

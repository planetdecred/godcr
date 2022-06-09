# This how we want to name the binary output
BINARY=godcr

VERSION="1.7.0"
BUILD=`date -u +"%Y-%m-%dT%H:%M:%SZ"`
# dev or prod
BuildEnv="prod" 

LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.BuildDate=${BUILD} -X main.BuildEnv=${BuildEnv}"
# LDFLAGSWIN adds the -H=windowsgui flag to windows build to prevent cli from starting alongside godcr
LDFLAGSWIN= -ldflags "-H=windowsgui -w -s -X main.Version=${VERSION} -X main.BuildDate=${BUILD} -X main.BuildEnv=${BuildEnv}"

all: clean macos windows linux freebsd

freebsd:
	GOOS=freebsd GOARCH=amd64 go build -trimpath ${LDFLAGS} -o ${BINARY}-freebsd-${GOARCH}
	GOOS=freebsd GOARCH=arm64 go build -trimpath ${LDFLAGS} -o ${BINARY}-freebsd-${GOARCH}

linux:
	GOOS=linux GOARCH=amd64 go build -trimpath ${LDFLAGS} -o ${BINARY}-freebsd-${GOARCH}
	GOOS=linux GOARCH=arm64 go build -trimpath ${LDFLAGS} -o ${BINARY}-freebsd-${GOARCH}

macos:
	GOOS=darwin GOARCH=amd64 go build -trimpath ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH}

windows:
	GOOS=windows GOARCH=amd64 go build -trimpath ${LDFLAGSWIN} -o ${BINARY}-windows-${GOARCH}.exe
 
# Cleans our project: deletes old binaries
clean:
	-rm -f ${BINARY}-*


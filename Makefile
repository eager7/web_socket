#/bin/bash
# This is how we want to name the binary output
TARGET=example
SRC=example.go
# These are the values we want to pass for Version and BuildTime
GITTAG=1.0.0
BUILD_TIME=`date +%Y%m%d%H%M%S`
# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X main.Version=${GITTAG} -X main.Build_Time=${BUILD_TIME} -s -w"

default: mod

mod:
	export GOPROXY=https://goproxy.io && GO111MODULE=on go build -v ${LDFLAGS} -o ${TARGET} ${SRC}

clean:
	-rm ${TARGET}

check:
	GO111MODULE=on golangci-lint run ./...
#!/usr/bin/env bash

#
# Helper scripts for Local
#

SCRIPT_DIR=$(cd $(dirname "$0"); pwd)

function usage() {
cat <<_EOT_
Usage:
  $0 Command

Example.
  $0 build

Command:
  build     Build app binary.
  fmt       Format sources.
  run       Run on local.
  test      Run test on local.
  install   Install dependency modules
_EOT_
exit 1
}

# Usage
usage() {
    echo "usage: scripts/local.sh build|run"
}

build() {
    cd "${SCRIPT_DIR}/../cmd" || exit
    go build -p 2 -v -x -mod vendor -tags=local local.go
}

fmt() {
    go fmt ./...
}

run() {
    cd "${SCRIPT_DIR}/../cmd" || exit
#    go run local.go
    ${SCRIPT_DIR}/../cmd/local
    #open http://localhost:8080
}

cmd_test() {
    cd ${SCRIPT_DIR}/..

    if [[ "$#" -ge 1 ]]; then
        if [[ "$1" == "nocache"  ]]; then
            ARGS="-count 1"
        fi
    fi

    # @see https://stackoverflow.com/questions/16353016/how-to-go-test-all-tests-in-my-project/35852900#35852900
    # NG
    #go test -v -cover "${ARGS}" ./...
    # OK
    go test -v -cover -tags="test local" ${ARGS} ./...
}

install() {
    go env -w GO111MODULE=on
    go mod vendor -v
}

if [[ $# -lt 1 ]]; then
    usage
fi

if [[ $1 = "build" ]]; then
    build
elif [[ $1 = "fmt" ]]; then
    fmt
elif [[ $1 = "run" ]]; then
    run
elif [[ $1 = "test" ]]; then
    cmd_test
elif [[ $1 = "install" ]]; then
    install
else
    usage
fi

#!/usr/bin/env bash

#
# Helper scripts for Local
#

SCRIPT_DIR=$(cd $(dirname "$0"); pwd)
ARGV=("$@")
ARGC=$#

function usage() {
cat <<_EOT_
Usage:
  $0 Command

Example.
  $0 build

Command:
  build         Build app binary.
  fmt           Format sources.
  run           Run on local.
  test          Run test on local.
  install       Install dependency modules
  docker:run    Run docker on local.
  docker:stop   Stop docker on local.
_EOT_
exit 1
}

build() {
    cd "${SCRIPT_DIR}/../cmd" || exit
    go build -p 2 -v -x -mod vendor main.go

    cd "${SCRIPT_DIR}/../cmd/migrate" || exit
    go build -p 2 -v -x -mod vendor migrate.go
}

fmt() {
    go fmt ./...
}

run() {
#    docker_run

#    # Call if it's entered Ctrl+C
#    trap docker_cleanup SIGINT

    echo Run go cmd.
    cd "${SCRIPT_DIR}/../cmd" || exit
    ${SCRIPT_DIR}/../cmd/main

#    docker_cleanup
}
docker_run() {
    docker-compose -f ${SCRIPT_DIR}/../deployments/local/docker-compose.yml up -d
}
docker_cleanup() {
    docker-compose -f ${SCRIPT_DIR}/../deployments/local/docker-compose.yml down
}

cmd_test() {
    docker_run

    cd ${SCRIPT_DIR}/..

    ARGS=""
    if [[ "$ARGC" -ge 2 ]]; then
        if [[ "${ARGV[1]}" == "nocache"  ]]; then
            ARGS="-count 1"
        fi
    fi
    if [[ "$ARGC" -ge 3 ]]; then
        ARGS="$ARGS -run ${ARGV[2]}"
    fi

    # @see https://stackoverflow.com/questions/16353016/how-to-go-test-all-tests-in-my-project/35852900#35852900
    # NG
    #go test -v -cover "${ARGS}" ./...
    # OK
    go test -v -cover -tags="test local" ${ARGS} ./...

    docker_cleanup
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
elif [[ $1 = "docker:run" ]]; then
    docker_run
elif [[ $1 = "docker:stop" ]]; then
    docker_cleanup
else
    usage
fi

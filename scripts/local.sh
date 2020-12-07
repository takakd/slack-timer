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
    #go build -p 2 -v -x -mod vendor main.go
    go build -p 2 -v -x main.go

    cd "${SCRIPT_DIR}/../cmd/migrate" || exit
    #go build -p 2 -v -x -mod vendor migrate.go
    go build -p 2 -v -x migrate.go
}

fmt() {
    go fmt ./...
    # Ref: https://gist.github.com/bgentry/fd1ffef7dbde01857f66#gistcomment-1618537
    goimports -w $(find . -type f -name "*.go" -not -path "./vendor/*")
    golint ./cmd/... ./internal/...
    go vet ./cmd/... ./internal/...
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
    if [[ "$DATABASE_URL" != "" ]]; then
        docker-compose -f ${SCRIPT_DIR}/../deployments/local/docker-compose.yml up -d
        # wait for DB is up
        sleep 2
    fi
}
docker_cleanup() {
    if [[ "$DATABASE_URL" != "" ]]; then
        docker-compose -f ${SCRIPT_DIR}/../deployments/local/docker-compose.yml down
    fi
}

cmd_test() {
    docker_run

    cd ${SCRIPT_DIR}/..

    # @see https://stackoverflow.com/questions/16353016/how-to-go-test-all-tests-in-my-project/35852900#35852900
    # NG
    #go test -v -cover "${ARGS}" ./...
    # OK
    sh -c "go $(echo ${ARGV[@]})"

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

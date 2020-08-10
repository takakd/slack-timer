#!/usr/bin/env bash

#
# Helper scripts for Heroku.
#

SCRIPT_DIR=$(cd $(dirname "$0"); pwd)

function usage() {
cat <<_EOT_
Usage:
  $0 Command

Example.
  $0 build

Command:
  build     Build heroku binary.
  run       Run on local.
_EOT_
exit 1
}

build() {
    cd "${SCRIPT_DIR}/../cmd" || exit
    go mod vendor
    go build -p 2 -v -x -mod vendor -tags=heroku -o ../bin/cmd ./heroku.go
}

run() {
    cd "${SCRIPT_DIR}/.." || exit
    echo "RUN: http://localhost:5000"
    heroku local
}

if [[ $# -lt 1 ]]; then
  usage
fi

if [[ $1 = "build" ]]; then
    build
elif [[ $1 = "run" ]]; then
  run
else
  usage
fi

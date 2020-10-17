#!/usr/bin/env bash

# golang-migrate helper script
#   https://github.com/golang-migrate/migrate

usage() {
    echo "Usage: migrate.sh COMMAND

COMMAND:
    Input golang-migrate command and options.
    'create' is limited to SQL.
    e.g.
        migrate.sh create NAME
        migrate.sh up
        migrate.sh down 2
"
}

if [[ "$1" = "help" ]]; then
    usage
    exit 0
fi

SCRIPT_DIR=$(cd $(dirname "$0"); pwd)
ENV_FILE=${SCRIPT_DIR}/../configs/.env

if [[ -e "$ENV_FILE" ]]; then
    source "$ENV_FILE"
fi

# Limit create to SQL
if [[ "$1" = "create" ]]; then
    # Exclude options
    param=$(echo "$@" | sed -e 's/create //')
    # Remove schema from env value
    dir=$(echo "$MIGRATION_DIR" | sed -e 's/file:\/\///')
    migrate -source "$MIGRATION_DIR" -database "$DATABASE_URL" create -dir $dir -ext sql $param
else
    migrate -source "$MIGRATION_DIR" -database "$DATABASE_URL" $@
fi

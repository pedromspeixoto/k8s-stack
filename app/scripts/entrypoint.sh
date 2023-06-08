#!/usr/bin/env bash
set -e

# load env file if exists
ENV_FILE=/opt/todo-api/service.env
if test -f "$ENV_FILE"; then
    echo "loading env variables"
    set -o allexport
    source $ENV_FILE
    set +o allexport
fi

# run binary
./todo-api
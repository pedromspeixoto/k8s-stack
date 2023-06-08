#!/usr/bin/env bash

go get -u github.com/swaggo/swag/cmd/swag@v1.8.3

export GO111MODULE=on
export GOOS=linux
export GOARCH=amd64

go mod tidy
go build -a -o bin/todo-api main.go
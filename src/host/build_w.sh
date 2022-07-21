#!/bin/sh

go env -w GOOS=windows
go build
go env -w GOOS=linux

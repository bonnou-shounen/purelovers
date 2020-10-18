#!/bin/sh

cd `dirname $0`

go build -o dist/purelovers cmd/purelovers/main.go

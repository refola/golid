#!/bin/sh
echo "Formatting code."
find . -type f -name \*.go -print0 | xargs -0 -r -n1 gofmt -w

echo "Running tests."
go test ./parse/

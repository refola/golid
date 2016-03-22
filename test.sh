#!/bin/sh

# Make sure we're at the script's location before doing stuff.
cd "$(dirname "$(readlink -f "$0")")" || exit 1

echo "Processing code."
# gofmt everything
find . -type f -name \*.go -print0 | xargs -0 -r -n1 gofmt -w

echo "Running tests."
go test ./parse/

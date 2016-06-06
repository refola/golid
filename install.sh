#!/bin/sh

# Make sure we're at the script's location before doing stuff.
cd "$(dirname "$(readlink -f "$0")")" || exit 1

echo -n "Formatting code... "
find . -type f -name \*.go -print0 | xargs -0 -r -n1 gofmt -w
echo -n "Installing... "
go install github.com/refola/golid/cmd/golid
echo "Done!"

#!/bin/sh

# Make sure we're at the script's location before doing stuff.
cd "$(dirname "$(readlink -f "$0")")" || exit 1

# Get to the latest version.
./install.sh

echo "Running tests."
go test ./parse/

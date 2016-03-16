#!/bin/bash
# Run all the .plgo tests to make sure they do what they're meant to.

# Make sure we're at the right place before doing stuff.
cd "$(dirname "$(readlink -f "$0")")" || exit 1
cd ./tests/classic || exit 1
echo "Running tests in $PWD."

for f in ./*; do
    f="${f/.\//}" # strip leading ./
    echo -e "\nRunning $f:"
    piklisp "$f" # convert to Go
    out="plgo_${f/.plgo/}.go"
    go run "$out" # run and remove generated file
    rm "$out"
done

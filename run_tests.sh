#!/bin/bash
# Run all the .gol tests to make sure they do what they're meant to.

## Usage: test-dir directory
# Tests every .gol file in given directory
test-dir() {
	echo "Running tests in $1"
	cd "$1" || exit 1
	for f in ./*; do
		if [ "$f" = "${f/.gol/}.gol" ]; then
			echo "Testing $f."
			f="${f/.\//}" # strip leading ./
			echo -e "\nRunning $f:"
			golid "$f" # convert to Go
			out="gol_${f/.gol/}.go"
			go run "$out" # run and remove generated file
			rm "$out"
		fi
	done
}

# Make sure we're at the right place before doing stuff.
test_dir="$(dirname "$(readlink -f "$0")")/tests" || exit 1
for d in "$test_dir"/*; do
	test-dir "$d"
done

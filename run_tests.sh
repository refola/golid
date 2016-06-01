#!/bin/bash
# Run all the .gol tests to make sure they do what they're meant to.

## Usage: clr color_code text
# Changes text color according to given code, from 0 thru 7.
clr() {
    echo -e "\e[0;3${1}m$2\e[0m"
}

## Usage: test-dir directory
# Tests every .gol file in given directory
test-dir() {
    clr 2 "Running tests in $1"
    cd "$1" || exit 1
    local f
    local out_file
    local output
    local expected
    local show="failed"
    if [ -f "./output.go" ]; then
        expected="$(go run ./output.go)"
    else
        clr 2 "No output.go in $1. Showing raw output instead."
        show="all"
    fi
    for f in ./*; do
        if [ "$f" = "${f/.gol/}.gol" ]; then
            f="${f/.\//}" # strip leading ./
            out_file="gol_${f/.gol/}.go" # Golid → Go name conversion
            if golid "$f" 2>/dev/null; then # convert to Go
                # run generated file if Golid worked
                if ! output="$(go run "$out_file" 2>/dev/null)"; then
                    clr 1 "Error in $f: Could not run generated '$out_file'."
                fi
                rm "$out_file" # and then clean up
                if [ "$show" = "all" ]; then
                    clr 7 "$f:"
                    echo "$output"
                elif [ "$output" != "$expected" ]; then
                    clr 1 "$f doesn't produce expected output."
                    clr 1 "$f's results:"
                    echo "$output"
                    clr 1 "\nExpected output:"
                    echo "$expected"
                    echo
                fi
            else
                clr 1 "Error! Golid failed on '$f'!"
            fi
        fi
    done
    echo
}

# Make sure we're at the right place before doing stuff.
test_dir="$(dirname "$(readlink -f "$0")")/tests" || exit 1
for d in "$test_dir"/*; do
    test-dir "$d"
done

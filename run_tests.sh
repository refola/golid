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
    echo -n " ..."
    local messages
    messages=("$(clr 2 "Testing $1")")
    cd "$1" || exit 1
    local f
    local out_file
    local output
    local expected
    local show="failed"
    local failed
    if [ -f "./output.go" ]; then
        expected="$(go run ./output.go)"
    else
        messages+=("$(clr 2 "No output.go in $1. Showing raw output instead.")")
        show="all"
    fi
    for f in ./*; do
        if [ "$f" = "${f/.gol/}.gol" ]; then
            f="${f/.\//}" # strip leading ./
            out_file="gol_${f/.gol/}.go" # Golid → Go name conversion
            if golid "$f" 2>/dev/null; then # convert to Go
                # run generated file if Golid worked
                if ! output="$(go run "$out_file" 2>/dev/null)"; then
                    failed="yes"
                    messages+=("$(clr 1 "Error in $f: Could not run generated '$out_file'.")")
                    messages+=("$(clr 2 "Run test.sh for debugging info.")")
                else
                    rm "$out_file" # and then clean up
                    if [ "$show" = "all" ]; then
                        messages+=("$(clr 7 "$f:")")
                        messages+=("$output")
                    elif [ "$output" != "$expected" ]; then
                        failed="yes"
                        messages+=("$(clr 1 "$f doesn't produce expected output.")")
                        messages+=("$(clr 1 "$f's results:")")
                        messages+=("$output" "")
                    fi
                fi
            else
                failed="yes"
                messages+=("$(clr 1 "Error! Golid failed on '$f'!")")
            fi
        fi
    done
    if [ -z "$failed" ]; then
        messages+=("$(clr 3 "Success!")")
    elif [ -z "$show" ]; then
        messages+=("$(clr 1 "At least one test failed. Here's the expected output.")")
        messages+=("$expected")
    fi
    if [ -n "$failed" ] || [ "$show" = "all" ] || [ -n "$DEBUG" ]; then
        echo
        for message in "${messages[@]}"; do
            echo -e "$message"
        done
    fi
}

# Make sure we're at the right place before doing stuff.
test_dir="$(dirname "$(readlink -f "$0")")/tests" || exit 1

# Make sure the latest Golid version is installed.
"$test_dir"/../install.sh

echo -n "Running tests. Please wait."

# Enable debug if told to
if [ "$1" = "debug" ]; then
    DEBUG="true"
fi

# Run the tests.
for d in "$test_dir"/*; do
    test-dir "$d"
done
echo

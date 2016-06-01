#!/bin/bash
# Run all the .gol tests to make sure they do what they're meant to.

## Usage: test-dir directory
# Tests every .gol file in given directory
test-dir() {
    echo -e "\e[0;32mRunning tests in $1\e[0m"
    cd "$1" || exit 1
    for f in ./*; do
        if [ "$f" = "${f/.gol/}.gol" ]; then
            f="${f/.\//}" # strip leading ./
            echo -e "\e[0;37mTesting $f.\e[0m"
            out="gol_${f/.gol/}.go" # Golid → Go name conversion
            if golid "$f" 2>/dev/null; then # convert to Go
                go run "$out" # run generated file if Golid worked
                rm "$out" # and then clean up
            else
                echo -e "\e[0;31mError! Golid failed on '$f'!\e[0m"
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

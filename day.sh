#!/bin/sh -e

if [ $# != 2 ] || \
  ! echo "$1" | grep -E -q '^[0-9]{4}$' || \
  ! echo "$2" | grep -E -q '^[0-9]{1,2}$'; then
  echo "Usage: $0 <YEAR> <DAY>" >&2
  exit 2
fi

# Remove zero-padding.
year=$(printf "%d" $1)
day=$(printf "%d" $2)

script_dir="$(dirname "$(realpath -s "$0")")"
dir="${script_dir}/$(printf "%04d/%02d" "$year" "$day")"

if [ ! -e "$dir" ]; then
  mkdir -p "$dir"
  cat <<EOF >"${dir}/main.go"
package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	for _, ln := range lib.InputLines("${year}/${day}") {
		fmt.Println(ln)
	}
}
EOF
fi

echo "$dir"

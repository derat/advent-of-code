#!/bin/sh -e

die() {
  echo "$1" >&2
  exit 2
}

usage() {
  prog=$(basename "$0")
  cat <<EOF >&2
Usage:
  $prog <YEAR> <DAY>   Prints (and inits) dir for specified year and day.
  $prog <DAY>          Prints dir for specified day in current dir's year.
  $prog today          Prints dir for today.
  $prog next           Prints dir for day after current dir.
  $prog prev           Prints dir for day before current dir.
  $prog input          Prints input for current dir.
  $prog web            Opens webpage for current dir.
  $prog lib            Prints library directory.
EOF
  exit 2
}

script_dir="$(dirname "$(realpath -s "$0")")"

# Figure out if we're already in a year/day or year directory.
cur_dir=$(pwd)
cur_year=
cur_day=
case "$cur_dir" in
  ${script_dir}/[0-9][0-9][0-9][0-9]/[0-9][0-9] | ${script_dir}/[0-9][0-9][0-9][0-9]/[0-9] )
    cur_year="$(basename "$(dirname "$cur_dir")")"
    cur_day="$(basename "$cur_dir")"
    break
    ;;
  ${script_dir}/[0-9][0-9][0-9][0-9] )
    cur_year="$(basename "$cur_dir")"
    break
    ;;
esac

# Dies with an error if not already in a year/day directory.
check_in_day_dir() {
  if [ -z "$cur_year" ] || [ -z "$cur_day" ]; then
    die "Must be in year/day directory"
  fi
}

year=
day=

if [ $# -eq 1 ]; then
  if [ "$1" = today ]; then
    [ $(date +%m) -eq 12 ] || die "Not in December"
    year=$(date +%Y)
    day=$(date +%d)
  elif [ "$1" = next ]; then
    check_in_day_dir
    year=$cur_year
    day=$(($cur_day + 1))
  elif [ "$1" = prev ]; then
    check_in_day_dir
    year=$cur_year
    day=$(($cur_day - 1))
  elif [ "$1" = input ]; then
    check_in_day_dir
    cat "$HOME/.cache/advent-of-code/$(printf "%d/%d" $cur_year $cur_day)"
    exit 0
  elif [ "$1" = web ]; then
    check_in_day_dir
    xdg-open "$(printf "https://adventofcode.com/%d/day/%d" $cur_year $cur_day)"
    exit 0
  elif [ "$1" = lib ]; then
    echo "${script_dir}/lib"
    exit 0
  else
    if [ -z "$cur_year" ]; then die "Must be in year or year/day directory"; fi
    if ! echo "$1" | grep -E -q '^[0-9]{1,2}$'; then usage; fi
    year="$cur_year"
    day="$1"
  fi
elif [ $# -eq 2 ]; then
  if ! echo "$1" | grep -E -q '^[0-9]{4}$' || \
     ! echo "$2" | grep -E -q '^[0-9]{1,2}$'; then
    usage
  fi
  year="$1"
  day="$2"
else
  usage
fi

if [ "$day" -lt 1 ] || [ "$day" -gt 25 ]; then
  die "Day $day not in range [1, 25]"
fi

# Remove zero-padding.
year=$(printf "%d" "$year")
day=$(printf "%d" "$day")

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

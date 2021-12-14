# advent-of-code

![Build Status](https://storage.googleapis.com/derat-build-badges/213a5a5e-f1c2-4738-abdf-0fb0a4a3dab4.svg)

This repository contains my Go solutions for [Advent of Code] programming
challenges, along with [related library code](./lib) and
[advent.sh](./advent.sh), a shell script that makes various common tasks easier.

I have a function similar to the following declared in my shell ([zsh]):

```sh
advent() {
  case "$1" in
    -h|--help|check|checkall|help|input|run|save|stdin|web)
      $HOME/advent-of-code/advent.sh "$@"
      ;;
    *)
      cd "$($HOME/advent-of-code/advent.sh "$@")"
      ;;
  esac
}
```

This lets me run commands like `advent today` to move to the directory for
today's puzzle.

[Advent of Code]: https://adventofcode.com/
[zsh]: https://en.wikipedia.org/wiki/Z_shell

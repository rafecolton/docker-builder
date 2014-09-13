#!/bin/bash

usage() {
  cat <<USAGE >&2

Usage: ./make.sh <command>

Commands:
-h/--help: show this message
ls: list
add: add current
rm: remove current
USAGE
}

main() {
  #touch $WIP_PATH
  local command="$1"
  shift

  if [[ -z $command ]] || [[ "$command" =~ -h|--help ]] ; then
    usage
    exit 1
  fi

  if ! type "make_${command}" >/dev/null 2>&1 ; then
    usage
    exit 2
  fi

  eval "cmd_${command} \"$@\""
}

cmd_release() {
  # binclean, gox-linux/darwin
  echo "cmd'ing"
}

main "$@"

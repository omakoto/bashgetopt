#!/bin/bash

set -e

cd "${0%/*}/.."

out=bin
mkdir -p "$out"

go build -o "$out/bashgetopt" ./src/cmd/bashgetopt

if [[ "$1" == "-r" ]] ; then
    shift
    "$out/bashgetopt" "$@"
fi

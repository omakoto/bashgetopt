#!/bin/bash

#set -e

cd "${0%/*}/.."

./scripts/build.sh || exit 1

my_usage() {
    echo "MY USAGE"
}

prog="$(./bin/bashgetopt -d 'COMMAND DESCRIPTION' '
s           echo short          # Short flag only
long        echo long           # Long flag only
t|test      echo both           # Short and long
str:          echo %              # String option
int=i         echo %              # Int option
' "$@")"

echo "=== SOURCE === "
echo "$prog"

echo
echo "=== RESULT === "
eval "$prog"
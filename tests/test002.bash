#!/bin/bash

#set -e

cd "${0%/*}/.."

./scripts/build.sh || exit 1

rc=0

export COMPROMISE_BASH_SKIP_BINDS=1
export COMPROMISE_QUIET=1

. <(./sample/pg --bash-completion)

complete -p pg | grep '^complete .* pg$'
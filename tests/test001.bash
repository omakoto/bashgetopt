#!/bin/bash

#set -e

cd "${0%/*}/.."

./scripts/build.sh || exit 1

rc=0

test() {
    echo "Test: $*"
    if ! diff --ignore-trailing-space -u <(cat) <(./sample/pg "$@"); then
        rc=1
        echo "FAIL: $*"
    fi
    return 0
}

test <<'EOF'
ignorecase: 0
context: 0
pattern:
files:
EOF

test -i <<'EOF'
ignorecase: 1
context: 0
pattern:
files:
EOF

test -C 3 <<'EOF'
ignorecase: 0
context: 3
pattern:
files:
EOF

test -i abc def xyz <<'EOF'
ignorecase: 1
context: 0
pattern: abc
files: def xyz
EOF

test -e PAT "x y" <<'EOF'
ignorecase: 0
context: 0
pattern: PAT
files: x y
EOF

exit $rc
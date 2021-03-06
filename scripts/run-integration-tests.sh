#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

# Use ginkgo to run integration tests. If arguments are provided to the
# script, they are treated as the directories containing the tests to run.
# When no arguments are provided, all integration tests are executed.

set -e -u

mchain_dir="$(cd "$(dirname "$0")/.." && pwd)"

# find packages that contain "integration" in the import path
integration_dirs() {
    local packages="$1"

    go list -f {{.Dir}} "$packages" | grep -E '/integration($|/)' | sed "s,${mchain_dir},.,g"
}

main() {
    cd "$mchain_dir"

    local -a dirs=("$@")
    if [ "${#dirs[@]}" -eq 0 ]; then
        dirs=($(integration_dirs "./..."))
    fi

    echo "Running integration tests..."
    ginkgo -keepGoing --slowSpecThreshold 60 -r "${dirs[@]}"
}

main "$@"

#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly dir
cd "${dir}/.."

# Stage 1 sync - intermediate to the ./vendir folder
set -x
vendir sync
{ set +x; } 2>/dev/null

# Patches
./sync/patches/templates-helper/patch.sh

# Store diffs
rm -f ./diffs/*
for f in $(git --no-pager diff --no-exit-code --no-color --no-index vendor/gateway-api helm/gateway-api-crds/templates --name-only) ; do
        [[ "$f" == "helm/gateway-api-crds/templates/_helpers.yaml" ]] && continue
        set +e
        set -x
        git --no-pager diff --no-exit-code --no-color --no-index "vendor/gateway-api/${f#"helm/gateway-api-crds/templates/"}" "${f}" \
                > "./diffs/${f//\//__}.patch" # ${f//\//__} replaces all "/" with "__"
        ret=$?
        { set +x; } 2>/dev/null
        set -e
        if [ $ret -ne 0 ] && [ $ret -ne 1 ] ; then
                exit $ret
        fi
done

#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

helpers_path="./helm/gateway-api-crds/templates/_helpers.yaml"

set -x
if ! [ -f "$helpers_path" ] ; then
    touch "$helpers_path"
fi
{ set +x; } 2>/dev/null

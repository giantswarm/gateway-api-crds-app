#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

raw_crds_path="../../../vendor/gateway-api-inference-extension"
templates_path="./helm/gateway-api-crds/templates/"

set -x

cd "${templates_path}"
yq --split-exp '.spec.group + "-" + .metadata.name + ".yaml"' "${raw_crds_path}"/*.yaml

cd "${repo_dir}"
{ set +x; } 2>/dev/null

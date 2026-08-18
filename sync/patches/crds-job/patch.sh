#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

readonly script_dir_rel=".${script_dir#"${repo_dir}"}"

templates_path="./helm/gateway-api-crds/templates"

set -x

# The CRDs themselves are baked into the installer image, so all the chart contributes is the
# hook Job that applies them plus the RBAC it needs.
mkdir -p "${templates_path}/crds"
cp "${script_dir_rel}/serviceaccount.yaml" "${templates_path}/crds/serviceaccount.yaml"
cp "${script_dir_rel}/rbac.yaml" "${templates_path}/crds/rbac.yaml"
cp "${script_dir_rel}/job.yaml" "${templates_path}/crds/job.yaml"
cp "${script_dir_rel}/networkpolicy.yaml" "${templates_path}/crds/networkpolicy.yaml"

{ set +x; } 2>/dev/null

#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

raw_crds_path="../../../vendor/gateway-api"
templates_path="./helm/gateway-api-crds/templates/"

set -x

cd "${templates_path}"
yq --split-exp '.metadata.annotations."gateway.networking.k8s.io/channel" + "-" + .metadata.name + ".yaml"' "${raw_crds_path}"/*.yaml

# The ValidatingAdmissionPolicy and its Binding share both the channel annotation and
# metadata.name, so the split above collapses them into a single file and the Binding
# wins. Drop that file and write one file per kind instead.
rm -f standard-safe-upgrades.gateway.networking.k8s.io.yaml
yq -N 'select(.kind == "ValidatingAdmissionPolicy")' \
  "${raw_crds_path}/standard-install.yaml" > standard-safe-upgrades-policy.gateway.networking.k8s.io.yaml
yq -N 'select(.kind == "ValidatingAdmissionPolicyBinding")' \
  "${raw_crds_path}/standard-install.yaml" > standard-safe-upgrades-binding.gateway.networking.k8s.io.yaml

cd "${repo_dir}"
{ set +x; } 2>/dev/null

#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

raw_crds_path="${repo_dir}/vendor/gateway-api-inference-extension"
crds_path="./crds"

set -x

# Inference Extension CRDs carry no channel annotation, so the channel is derived from the
# API group. yq's --split-exp only accepts a concatenation expression, hence one call per
# channel instead of a conditional.
mkdir -p "${crds_path}/standard" "${crds_path}/experimental"

cd "${crds_path}/standard"
yq -N --split-exp '.spec.names.plural + ".yaml"' \
  'select(.kind == "CustomResourceDefinition" and .spec.group == "inference.networking.k8s.io")' \
  "${raw_crds_path}"/*.yaml

cd "${repo_dir}/${crds_path}/experimental"
yq -N --split-exp '.spec.names.plural + ".yaml"' \
  'select(.kind == "CustomResourceDefinition" and .spec.group == "inference.networking.x-k8s.io")' \
  "${raw_crds_path}"/*.yaml

cd "${repo_dir}"
{ set +x; } 2>/dev/null

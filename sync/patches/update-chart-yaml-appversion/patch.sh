#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

chart_yaml_path="./helm/gateway-api-crds/Chart.yaml"

set -x

# { set +x; } 2>/dev/null
gateway_api_version=$(yq e '.directories[0].contents[] | select(.path == "gateway-api") | .githubRelease.tag' vendir.yml)
inference_extension_version=$(yq e '.directories[0].contents[] | select(.path == "gateway-api-inference-extension") | .githubRelease.tag' vendir.yml)

yq -i ".appVersion = \"gateway-api: $gateway_api_version, gateway-api-inference-extension: $inference_extension_version\"" $chart_yaml_path

{ set +x; } 2>/dev/null

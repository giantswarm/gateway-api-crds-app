#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

crds_path="./crds"

set -x

# The CRDs are applied by the installer Job, so Helm never owns them and this annotation is
# inert under server-side apply. It is kept as a safety net: clusters that installed an
# earlier chart version have Helm-managed CRDs carrying it, and Helm skips deleting
# annotated resources that disappear from the release manifest on upgrade.
# One yq call per file: passing several files to a single "yq -i" merges all of their
# documents into one output file.
find "${crds_path}" -type f -name '*.yaml' \
  -exec yq -i '.metadata.annotations += {"helm.sh/resource-policy":"keep"}' {} \;

{ set +x; } 2>/dev/null

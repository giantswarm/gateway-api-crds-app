#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

raw_crds_path="${repo_dir}/vendor/gateway-api"
crds_path="./crds"
templates_path="./helm/gateway-api-crds/templates"

set -x

# The CRDs are baked into the installer image, not rendered by Helm, so they live outside
# the chart directory. Start from scratch so upstream renames cannot leave orphans behind.
rm -rf "${crds_path}"
mkdir -p "${crds_path}"

# Drop the previously generated CRD templates. Only the admission policy templates are
# regenerated below; _helpers.tpl and templates/crds/ are left alone.
find "${templates_path}" -maxdepth 1 -type f -name '*.yaml' -delete

cd "${crds_path}"
yq -N --split-exp '.metadata.annotations."gateway.networking.k8s.io/channel" + "/" + .spec.names.plural + ".yaml"' \
  'select(.kind == "CustomResourceDefinition")' "${raw_crds_path}"/*.yaml

cd "${repo_dir}"

# The ValidatingAdmissionPolicy and its Binding are not CRDs. They are small enough to stay
# Helm-managed, so they are emitted as chart templates instead of into the image.
yq -N 'select(.kind == "ValidatingAdmissionPolicy")' \
  "vendor/gateway-api/standard-install.yaml" > "${templates_path}/safe-upgrades-policy.yaml"
yq -N 'select(.kind == "ValidatingAdmissionPolicyBinding")' \
  "vendor/gateway-api/standard-install.yaml" > "${templates_path}/safe-upgrades-binding.yaml"

# Gate them on .Values.install.admissionPolicies, as the CRD templates used to be gated.
for f in "${templates_path}/safe-upgrades-policy.yaml" "${templates_path}/safe-upgrades-binding.yaml" ; do
  { echo '{{- if .Values.install.admissionPolicies }}' ; cat "${f}" ; echo '{{- end }}' ; } > "${f}.tmp"
  mv "${f}.tmp" "${f}"
done

{ set +x; } 2>/dev/null

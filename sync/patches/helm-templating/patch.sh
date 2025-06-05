#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

templates_path="./helm/gateway-api-crds/templates/"

set -x

cd "${templates_path}"

{ set +x; } 2>/dev/null
for f in *.yaml ; do
  [[ "$f" == "_helpers.yaml" ]] && continue

  name=$(yq '.spec.names.plural' $f)
  channel=$(yq '.metadata.annotations."gateway.networking.k8s.io/channel"' $f)
  set -x
  echo "{{ if eq .Values.install.$name \"$channel\" }}" | cat - $f > _temp #
  echo "{{ end }}" >> _temp
  mv _temp $f

  { set +x; } 2>/dev/null
done

set -x
cd "${repo_dir}"
{ set +x; } 2>/dev/null

{{- define "gateway-api-crds.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "gateway-api-crds.installerName" -}}
{{- printf "%s-installer" (include "gateway-api-crds.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "gateway-api-crds.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "gateway-api-crds.labels" -}}
app.kubernetes.io/name: {{ include "gateway-api-crds.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ include "gateway-api-crds.chart" . }}
application.giantswarm.io/team: {{ index .Chart.Annotations "io.giantswarm.application.team" | default "cabbage" | quote }}
{{- end -}}

{{/*
Installer image reference. The tag defaults to the chart version, which the release pipeline
keeps in lockstep with the pushed image. Set crds.image.tag to pin a branch build.
*/}}
{{- define "gateway-api-crds.image" -}}
{{- $registry := .Values.crds.image.registry -}}
{{/* Nested rather than a single "and": older Helm evaluates all "and" arguments, so a nil
.Values.global would be dereferenced. */}}
{{- if .Values.global -}}
{{- if .Values.global.image -}}
{{- if .Values.global.image.registry -}}
{{- $registry = .Values.global.image.registry -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- $tag := .Values.crds.image.tag | default (.Chart.Version | splitList "+" | first) -}}
{{- printf "%s/%s:%s" $registry .Values.crds.image.repository $tag -}}
{{- end -}}

{{/*
Space separated paths of the CRD files selected through .Values.install, as they are laid out
in the installer image. The key is the CRD's plural name and the value is its channel, so a
values entry addresses a file directly. Empty when nothing is selected, which the templates
use to render no installer at all.
*/}}
{{- define "gateway-api-crds.selectedFiles" -}}
{{- $files := list -}}
{{- range $name, $channel := .Values.install -}}
{{/* Nested rather than a single "and": older Helm evaluates all "and" arguments, so
comparing the boolean admissionPolicies against a string would error out. */}}
{{- if kindIs "string" $channel -}}
{{- if or (eq $channel "standard") (eq $channel "experimental") -}}
{{- $files = append $files (printf "/crds/%s/%s.yaml" $channel $name) -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- join " " $files -}}
{{- end -}}

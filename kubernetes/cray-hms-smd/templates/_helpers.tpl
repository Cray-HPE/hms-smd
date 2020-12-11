{{/*
Add helper methods here for your chart
*/}}

{{- define "cray-hms-smd.image-prefix" -}}
{{ $base := index . "cray-service" }}
{{- if $base.imagesHost -}}
{{- printf "%s/" $base.imagesHost -}}
{{- else -}}
{{- printf "" -}}
{{- end -}}
{{- end -}}

{{/*
Helper function to get the proper image tag
*/}}
{{- define "cray-hms-smd.imageTag" -}}
{{- default "latest" .Chart.AppVersion -}}
{{- end -}}
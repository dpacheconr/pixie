{{- /* Defines tolerations */ -}}
{{- define "custom.tolerations" -}}
    {{- if .Values.tolerations -}}
        {{- toYaml .Values.tolerations -}}
    {{- end -}}
{{- end -}}

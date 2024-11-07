{{/* The line below tells Intellij/GoLand to enable the autocompletion based on the *gen.Graph type. */}}
{{/* gotype: entgo.io/ent/entc/gen.Graph */}}

{{ define "crud/helper/update" }}

{{ $builder := .UpdateName }}
{{ $receiver := .UpdateReceiver }}
{{ $fields := .Fields }}
{{- if or (hasSuffix $builder "Update") (hasSuffix $builder "UpdateOne") }}
{{ $fields = .MutableFields }}
{{- end }}

{{ print "// Set" .Name "Full set the " .Name }}
func ({{ $receiver }} *{{ $builder }}) Set{{ .Name }}Full(input *{{ .Name }}) *{{ $builder }} {
{{- range $f := $fields }}
{{- if $f.Nillable}}
{{"if"}} input.{{ $f.StructField }} {{"!= nil {"}}
{{- $setter := print "Set" $f.StructField }}
{{ $receiver }}.{{ $setter }}(*input.{{ $f.StructField }})
{{"}"}}
{{- else}}
{{- if $f.IsTime }}
{{- $setter := print "Set" $f.StructField }}
    {{"if !"}}input.{{ $f.StructField }}.IsZero() {{"{"}}
        {{ $receiver }}.{{ $setter }}(input.{{ $f.StructField }})
    {{"}"}}
{{- else }}
{{- $setter := print "Set" $f.StructField }}
{{ $receiver }}.{{ $setter }}(input.{{ $f.StructField }})
{{- end}}
{{- end}}
{{- end }}
return {{ $receiver }}
}

{{ print "// Set" .Name " set the " .Name }}
func ({{ $receiver }} *{{ $builder }}) Set{{ .Name }}(input *{{ .Name }}) *{{ $builder }} {
{{- range $f := $fields }}
{{- if not $f.Optional }}
{{- if $f.Nillable}}
{{"if"}} input.{{ $f.StructField }} {{"!= nil {"}}
{{- $setter := print "Set" $f.StructField }}
{{ $receiver }}.{{ $setter }}(*input.{{ $f.StructField }})
{{"}"}}
{{- else}}
{{- if $f.IsTime }}
{{- $setter := print "Set" $f.StructField }}
{{"if !"}}input.{{ $f.StructField }}.IsZero() {{"{"}}
{{ $receiver }}.{{ $setter }}(input.{{ $f.StructField }})
{{"}"}}
{{- else }}
{{- $setter := print "Set" $f.StructField }}
{{ $receiver }}.{{ $setter }}(input.{{ $f.StructField }})
{{- end}}

{{- end}}
{{- end}}
{{- end}}
return {{ $receiver }}
}


{{- end -}}

{{/* The line below tells Intellij/GoLand to enable the autocompletion based on the *gen.Graph type. */}}
{{/* gotype: entgo.io/ent/entc/gen.Graph */}}

{{ define "crud/helper/updateone" }}

{{ $builder := .UpdateOneName }}
{{ $receiver := .UpdateOneReceiver }}
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


// Omit allows the unselect one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func ({{ $receiver }} *{{ $builder }}) Omit(fields ...string) *{{ $builder }} {
omits := make(map[string]struct{}, len(fields))
for i := range fields {
omits[fields[i]] = struct{}{}
}
{{ $receiver }}.fields = []string(nil)
for _, col := range {{ .Package }}.Columns {
if _, ok := omits[col]; !ok {
{{ $receiver }}.fields = append({{ $receiver }}.fields, col)
}
}
return {{ $receiver }}
}

{{- end -}}

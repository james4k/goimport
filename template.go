package goimport

import (
	"html/template"
	"io"
)

var tmpl = template.Must(template.New("tmpl").Parse(`<html>
<head>
{{$host := .Host}}
{{range .Imports}}
<meta name="go-import" content="{{$host}} {{.VCS}} {{.Path}}">
{{end}}
</head>
<body>
</body>
</html>`))

func render(w io.Writer, data interface{}) {
	_ = tmpl.Execute(w, data)
	// We could catch the error, but we would either panic, or allow the
	// user to specify a logger which will just complicate the
	// package. The only errors that will happen are likely to be
	// network related.
}

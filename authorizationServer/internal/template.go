package internal

import (
	"html/template"
	"path/filepath"
)

func mustLoadTemplates(pattern string) *template.Template {
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return tmpl
}

var Tmpl = mustLoadTemplates("web/html/*")

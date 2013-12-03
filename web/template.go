package web

import (
	"bytes"
	"html/template"
)

var pageTemplates *template.Template

func TemplateLoad(path string) {
	pageTemplates = template.Must(template.ParseGlob(path))
}

func ProcessTemplate(tname string, ob interface{}) (string, error) {
	bb := new(bytes.Buffer)
	err := pageTemplates.ExecuteTemplate(bb, tname, ob)
	if err != nil {
		return "", err
	}
	return bb.String(), nil
}

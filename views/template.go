package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "There was error executing the template", http.StatusInternalServerError)
		return
	}
}

func Must(tpl Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return tpl
}

func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %v", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	tpl, err := template.ParseFS(fs, pattern...)
	if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %v", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

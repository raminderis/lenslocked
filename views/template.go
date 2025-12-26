package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	t, err := template.ParseFS(fs, patterns...)
	if err != nil {
		log.Printf("parsing template error : %v", err)
		return Template{}, fmt.Errorf("parsing tempalte: %w", err)
	}
	return Template{
		htmlTpl: t,
	}, nil
}

func Parse(filepath string) (Template, error) {
	t, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("parsing template error : %v", err)
		return Template{}, fmt.Errorf("parsing tempalte: %w", err)
	}
	return Template{
		htmlTpl: t,
	}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template error : %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

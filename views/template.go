package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/raminderis/lenslocked/context"
	"github.com/raminderis/lenslocked/models"
)

type public interface {
	Public() string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	t := template.New(patterns[0])
	t = t.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("csrfField not implemented")
		},
		"currentUser": func() (template.HTML, error) {
			return "", fmt.Errorf("currentUser not implemented")
		},
		"errors": func() (template.HTML, error) {
			return "", fmt.Errorf("errors not implemented")
		},
	})
	t, err := t.ParseFS(fs, patterns...)
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

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template error: %v", err)
		http.Error(w, "Error cloning template", http.StatusInternalServerError)
		return
	}
	errMsgs := errMessages(errs)
	tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errMsgs
			},
		},
	)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer //may cause performance issues. as it loads the entire page in mem before rendring to response writer.
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template error : %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func errMessages(errs []error) []string {
	var errMsgs []string
	for _, err := range errs {
		var pubError public
		if errors.As(err, &pubError) {
			errMsgs = append(errMsgs, pubError.Public())
		} else {
			fmt.Println(err)
			errMsgs = append(errMsgs, "Something went wrong. Please try again.")
		}
	}
	return errMsgs
}

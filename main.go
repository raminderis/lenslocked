package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type User struct {
	Email string
	Phone string
}

func executeTemplate(w http.ResponseWriter, filepath string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, "Error parsing template ", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath, nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	data := User{
		Email: "contactis@live.com",
		Phone: "4253222555",
	}
	executeTemplate(w, tplPath, data)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")
	data := User{
		Email: "faqis@live.com",
		Phone: "4253111555",
	}
	executeTemplate(w, tplPath, data)
}

func rawPathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Path: "+r.URL.Path+"</h1>")
	fmt.Fprint(w, "<h1>Raw Path: "+r.URL.RawPath+"</h1>")
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.Get("/dog/*", rawPathHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
	fmt.Println("With a branch starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}

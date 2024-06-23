package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("Parsing template error: %v", err)
		http.Error(w, "There was error parsing the template.", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("Execution template error: %v", err)
		http.Error(w, "There was error executing the template.", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tpath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tpath)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	w.Write([]byte(fmt.Sprintf("hi %v", userID)))
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tpath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tpath)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>FAQ Page</h1>
	<ul>
		<li>
			<b>Is there a free version</b>
			Yes
		</li>
		<li>
			<b>Is there a free version</b>
			Yes
		</li>
		<li>
			<b>Is there a free version</b>
			EMail us at <a href=\"mailto:rami@liv.com\">rami@liv.com</a>
		</li>
	</ul>
	`)
}

func main() {
	r := chi.NewRouter()
	r.With(middleware.Logger).Get("/users/{userID}", usersHandler)
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
	//var router Router
	fmt.Println("listening now on 3000")
	err := http.ListenAndServe("127.0.0.1:3000", r)
	if err != nil {
		panic(err)
	}
}

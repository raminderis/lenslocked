package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact us at</h1><p> Email : <a href=\"mailto:rami@live.com\">raminder@live.com</a>.</p>")
}

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
	}
}

func main() {
	var router Router
	fmt.Println("Starting the server on :3000...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		panic(err)
	}
}

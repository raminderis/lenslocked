package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//set status code as 200 ok
	//w.WriteHeader(http.StatusOK)
	//set content-type header
	bio := `&lt;script&gt;alert(&#34;hello you are hacked&#34;)&lt;/script&gt;`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := fmt.Fprint(w, "<h1>Welcome to my awesome site</h1><p>Bio: "+bio+"</p>")
	if err != nil {
		return
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact us</h1><p> Get in touch with me at <a href=\"mailto:raminder@live.com\">raminder@live.com</a></p>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>FAQ Page:</h1>"+
		"<ul>"+
		"<li><b>Do you have free version</b><p>A. gibberish</p>"+
		"<li><b>Do you have paid version</b><p>A. yes</p>"+
		"</ul>")
}

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}

//7% 1:27pm

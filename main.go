package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
	<h2>
		My Contact:
	</h2>
	<p> Email me at <a href=\"mailto:jon@live.com\">jon@live.com</a.</p>>
	`)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	w.Write([]byte(fmt.Sprintf("hi %v", userID)))
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
	<h2>
		My Contact:
	</h2>
	<p> Email me at <a href=\"mailto:jon@live.com\">jon@live.com</a.</p>>
	`)
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

// type Router struct{}

//	func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("goning  now on 3000")
//		switch r.URL.Path {
//		case "/":
//			homeHandler(w, r)
//		case "/faq":
//			faqHandler(w, r)
//		case "/contact":
//			contactHandler(w, r)
//		default:
//			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		}
//	}
func main() {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome"))
	// })
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/users/{userID}", usersHandler)
	})
	r.Get("/", homeHandler)
	r.Get("/contact", homeHandler)
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

package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/raminderis/lenslocked/controllers"
	"github.com/raminderis/lenslocked/templates"
	"github.com/raminderis/lenslocked/views"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	data := controllers.User{
		Email: "raminder@love.com",
		Phone: "4253000114",
		QA: []controllers.Questions{
			{
				Question: "What is your name?",
				Answer:   "Something",
			},
			{
				Question: "What is your ding?",
				Answer:   "sulu",
			},
			{
				Question: "What is your dong?",
				Answer:   "hini",
			},
		},
	}

	t := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(t, data))

	t = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(t, data))

	t = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.StaticHandler(t, data))

	t = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	r.Get("/signup", controllers.StaticHandler(t, data))

	t = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Post("/users", controllers.StaticHandler(t, data))

	t = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	r.Get("/signin", controllers.StaticHandler(t, data))

	t = views.Must(views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml"))
	r.Get("/reset-pw", controllers.StaticHandler(t, data))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
	fmt.Println("With a branch starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}

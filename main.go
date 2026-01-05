package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/raminderis/lenslocked/controllers"
	"github.com/raminderis/lenslocked/migrations"
	"github.com/raminderis/lenslocked/models"
	"github.com/raminderis/lenslocked/templates"
	"github.com/raminderis/lenslocked/views"
)

func timeHandlerProcessing(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h(w, r)
		fmt.Println("Request Time: ", time.Since(start))
	}
}

func main() {
	// Setup the database
	dbCfg := models.DefaultPostgresConfig()
	pgxConn, err := models.Open(dbCfg)
	if err != nil {
		panic(err)
	}
	defer pgxConn.Close(context.Background())
	fmt.Println("Connected to DB")

	err = dbCfg.MigrateFS(migrations.FS)
	if err != nil {
		panic(err)
	}

	// Setup Services aka Initialize Controller
	usersC := controllers.Users{}
	usersC.Templates.General = views.Must(views.ParseFS(templates.FS, "general-page.gohtml", "tailwind.gohtml"))
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.Signin = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))

	usersC.UserService = &models.UserService{
		DB_CONN: pgxConn,
	}
	usersC.SessionService = &models.SessionService{
		DB_CONN:       pgxConn,
		BytesPerToken: 32,
	}

	// Setup Middleware
	umw := controllers.UserMiddleware{
		SessionService: usersC.SessionService,
	}
	csrfKey := "Q7f9K2pL8xR3mV1tC6zH4bN0wP5sJ8dF"
	csrfMiddleware := csrf.Protect([]byte(csrfKey), csrf.Secure(true), csrf.TrustedOrigins([]string{"localhost:3000"}))

	// Set up router and routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(csrfMiddleware)
	r.Use(umw.SetUser)
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
	r.Get("/faq", timeHandlerProcessing(controllers.StaticHandler(t, data)))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.Signin)
	r.Post("/signin", usersC.SigninProcess)
	r.Post("/signout", usersC.ProcessSignout)
	// r.Get("/users/me", usersC.CurrentUser)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello from user me")
		})

	})
	t = views.Must(views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml"))
	r.Get("/reset-pw", controllers.StaticHandler(t, data))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	fmt.Println("With a branch starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}

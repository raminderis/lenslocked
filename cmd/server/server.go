package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
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

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
	}
	if cfg.PSQL.Host == "" || cfg.PSQL.Port == "" || cfg.PSQL.User == "" || cfg.PSQL.Password == "" || cfg.PSQL.Database == "" {
		return cfg, fmt.Errorf("incomplete PSQL configuration in environment variables")
	}
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	smtp_ssl := os.Getenv("SMTP_SSL")
	cfg.SMTP.DialerSSL, err = strconv.ParseBool(smtp_ssl)
	if err != nil {
		return cfg, err
	}

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")
	return cfg, nil
}

func main() {
	// Load config
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Setup the database
	// dbCfg := models.DefaultPostgresConfig()
	pgxConn, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer pgxConn.Close(context.Background())
	fmt.Println("Connected to DB")

	err = cfg.PSQL.MigrateFS(migrations.FS)
	if err != nil {
		panic(err)
	}

	// Setup Services aka Initialize Controller
	usersC := controllers.Users{}
	usersC.Templates.General = views.Must(views.ParseFS(templates.FS, "general-page.gohtml", "tailwind.gohtml"))
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.Signin = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "forgot-pw.gohtml", "tailwind.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "check-your-email.gohtml", "tailwind.gohtml"))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml"))

	usersC.UserService = &models.UserService{
		DB_CONN: pgxConn,
	}
	usersC.SessionService = &models.SessionService{
		DB_CONN:       pgxConn,
		BytesPerToken: 32,
	}
	usersC.PasswordResetService = &models.PasswordResetService{
		DB_CONN:       pgxConn,
		BytesPerToken: 32,
	}
	usersC.EmailService = models.NewEmailService(cfg.SMTP)

	// Setup Middleware
	umw := controllers.UserMiddleware{
		SessionService: usersC.SessionService,
	}
	csrfKey := cfg.CSRF.Key
	csrfMiddleware := csrf.Protect([]byte(csrfKey), csrf.Secure(cfg.CSRF.Secure), csrf.TrustedOrigins([]string{"localhost:3000"}))

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
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)
	// r.Get("/users/me", usersC.CurrentUser)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello from user me")
		})

	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	fmt.Printf("With a branch starting the server on %s...", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}

package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/raminderis/lenslocked/context"
	"github.com/raminderis/lenslocked/models"
)

type Users struct {
	Templates struct {
		New            Template
		Signin         Template
		General        Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
	}
	Message              string
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
}

type Account struct {
	Email    string
	Password string
	Message  string
	Token    string
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	recievedAcct := Account{
		Email: r.FormValue("email"),
	}
	u.Templates.New.Execute(w, r, recievedAcct)
}

func (u Users) Signin(w http.ResponseWriter, r *http.Request) {
	recievedAcct := Account{
		Email: r.FormValue("email"),
	}
	u.Templates.Signin.Execute(w, r, recievedAcct)
}

func (u Users) General(w http.ResponseWriter, r *http.Request) {
	recievedAcct := Account{
		Message: u.Message,
	}
	u.Templates.General.Execute(w, r, recievedAcct)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := r.PostForm
	newAcct := Account{
		Email:    s.Get("email"),
		Password: s.Get("password"),
	}
	// fmt.Fprintf(w, "Attempting to create user with email %v and password notprintedHere\n", newAcct.Email)
	user, err := u.UserService.Create(newAcct.Email, newAcct.Password)
	if err != nil {
		fmt.Println(err)
		u.Templates.New.Execute(w, r, newAcct, err)
		// http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
	// fmt.Fprintf(w, "User created: %v", user.Email)
	// recievedAcct := Account{
	// 	Message: "User " + user.Email + " is created",
	// }
	// u.Templates.General.Execute(w, r, recievedAcct)
}

func (u Users) SigninProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := r.PostForm
	newAcct := Account{
		Email:    s.Get("email"),
		Password: s.Get("password"),
	}
	// fmt.Fprintf(w, "Attempting to Login User with email %v and password notprintedHere\n", newAcct.Email)
	user, err := u.UserService.Authenticate(newAcct.Email, newAcct.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	// fmt.Fprintf(w, "User authenticated: %+v", user.Email)
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)

	// recievedAcct := Account{
	// 	Message: "User " + user.Email + " is authenticated",
	// }
	// u.Templates.General.Execute(w, r, recievedAcct)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := context.User(ctx)
	if user == nil {
		http.Redirect(w, r, "signin", http.StatusFound)
		return
	}

	// token, err := readCookie(r, CookieSession)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Redirect(w, r, "signin", http.StatusFound)
	// 	return
	// }
	// user, err := u.SessionService.User(token)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Redirect(w, r, "signin", http.StatusFound)
	// 	return
	// }
	recievedAcct := Account{
		Message: "Current User is: " + user.Email,
	}
	u.Templates.General.Execute(w, r, recievedAcct)
}

func (u Users) ProcessSignout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	token, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		recievedAcct := Account{
			Message: "Signout Process something went wrong. " + err.Error(),
		}
		u.Templates.General.Execute(w, r, recievedAcct)
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "signin", http.StatusFound)
}

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	recievedAcct := Account{
		Email: r.FormValue("email"),
	}
	u.Templates.ForgotPassword.Execute(w, r, recievedAcct)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	recievedAcct := Account{
		Email: r.FormValue("email"),
	}
	pwReset, err := u.PasswordResetService.Create(recievedAcct.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetURL := "http://localhost:3000/reset-pw?" + vals.Encode()
	err = u.EmailService.ForgotPassword(recievedAcct.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	u.Templates.CheckYourEmail.Execute(w, r, recievedAcct)
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	recievedAcct := Account{
		Token: r.FormValue("token"),
	}
	u.Templates.ResetPassword.Execute(w, r, recievedAcct)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	recievedAcct := Account{
		Token:    r.FormValue("token"),
		Password: r.FormValue("password"),
	}

	user, err := u.PasswordResetService.Consume(recievedAcct.Token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	err = u.UserService.UpdatePassword(user.ID, recievedAcct.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//Sign in the user now that the pw has be reset.
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.User(token)
		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

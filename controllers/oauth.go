package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
)

type OAuth struct {
	ProviderConfigs map[string]*oauth2.Config
}

func (oa OAuth) Connect(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)
	config, ok := oa.ProviderConfigs[provider]
	if !ok {
		http.Error(w, "Unknown provider or Invaliud OAuth2 Service", http.StatusBadRequest)
		return
	}
	verifier := oauth2.GenerateVerifier()
	setCookie(w, "pkce_verifier", verifier)
	state := csrf.Token(r)
	setCookie(w, "oauth_state", state)
	url := config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
		oauth2.SetAuthURLParam("token_access_type", "offline"),
		oauth2.SetAuthURLParam("redirect_uri", redirectURI(r, provider)),
	)
	http.Redirect(w, r, url, http.StatusFound)
}

func (oa OAuth) Callback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)
	config, ok := oa.ProviderConfigs[provider]
	if !ok {
		http.Error(w, "Unknown provider or Invaliud OAuth2 Service", http.StatusBadRequest)
		return
	}
	state := r.FormValue("state")
	cookieState, err := readCookie(r, "oauth_state")
	if err != nil || cookieState != state {
		if err != nil {
			fmt.Println("Error reading cookie:", err)
		}
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}
	cookieVerifier, err := readCookie(r, "pkce_verifier")
	if err != nil {
		if err != nil {
			fmt.Println("Error reading cookie:", err)
		}
		http.Error(w, "Invalid Verifier state", http.StatusBadRequest)
		return
	}
	deleteCookie(w, "oauth_state")
	code := r.FormValue("code")
	tok, err := config.Exchange(
		r.Context(),
		code,
		oauth2.SetAuthURLParam("redirect_uri", redirectURI(r, provider)),
		oauth2.VerifierOption(cookieVerifier),
	)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	//Presist the token or use it to make API calls
	//Redirect to previous page or dashboard
	// w.Header().Set("Content-Type", "application/json")
	// fmt.Printf("Token: %+v\n", tok)
	// enc := json.NewEncoder(w)
	// enc.SetIndent("", "  ")
	// enc.Encode(tok)

	client := config.Client(r.Context(), tok)
	resp, err := client.Post(
		"https://api.dropboxapi.com/2/files/list_folder",
		"application/json",
		strings.NewReader(`{"path": ""}`))
	if err != nil {
		http.Error(w, "API request failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	unprettyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read API response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var pretty bytes.Buffer
	err = json.Indent(&pretty, unprettyBytes, "", "  ")
	if err != nil {
		http.Error(w, "Failed to format API response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	pretty.WriteTo(w)
	// io.Copy(w, &pretty)
	fmt.Println("Response status:", resp.Status)
}

func redirectURI(r *http.Request, provider string) string {
	if r.Host == "localhost:3000" {
		return "http://" + r.Host + "/oauth/" + provider + "/callback"
	}
	return "http://" + r.Host + "/oauth/" + provider + "/callback"
	// return "https://mywebsite.com/oauth/" + provider + "/callback"
}

package main

import (
	"fmt"
	"net/http"
	"strings"
)

// ResponseWriter is an interface which implements many methods one of which is the Write method.
// The Write method is used to write the response back to the client.
// Request is the pointer to the struct which contains all the information about the incoming request.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>Contact Us</h1>
	<p> email: <a href="mailto:raminderis@live.com">ramidneris@live.com</a></p> 
	<p> phone: <a href="phoneto:4253000115">4253000115</a></p>`)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>Frequently Asked Questions</h1>
	<ul>
		<li><b>Question 1: What is your return policy?</b>
			<div>Answer:</div>
			<ul>
				<li>We accept returns within 30 days of purchase.</li>
				<li>Items must be in original condition.</li>
				<li>Refunds will be processed within 5-7 business days.</li>
			</ul>
		</li>
		<li><b>Question 2: Do you offer international shipping?</b>
			<div>Answer:</div>
			<ul><li>Yes, we ship to most countries worldwide.</li></ul>
		</li>
		<li><b>Question 3: How can I track my order?</b>
			<div>Answer:</div>
			<ul><li>You will receive a tracking number via email once your order has shipped.</li></ul>
		</li>
		<li><b>Question 4: if you have any questions please contact us at?</b>
			<div>Answer:</div>
			<ul><li>Email: <a href="mailto:raminderis@live.com">raminderis@live.com</a></li></ul>
		</li>
	</ul>`)
}

func rawPathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Path: "+r.URL.Path+"</h1>")
	fmt.Fprint(w, "<h1>Raw Path: "+r.URL.RawPath+"</h1>")
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		homeHandler(w, r)
	case r.URL.Path == "/contact":
		contactHandler(w, r)
	case strings.HasPrefix(r.URL.Path, "/dog/"):
		rawPathHandler(w, r)
	case r.URL.Path == "/faq":
		faqHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		homeHandler(w, r)
	case r.URL.Path == "/contact":
		contactHandler(w, r)
	case r.URL.Path == "/faq":
		faqHandler(w, r)
	case strings.HasPrefix(r.URL.Path, "/dog/"):
		rawPathHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func main() {
	var somePathHandler Router
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handlerFunc)
	// http.HandleFunc("/", pathHandler)
	// http.HandleFunc("/contact", contactHandler)
	// http.HandleFunc("/path/", pathHandler)
	// var router http.HandlerFunc = pathHandler
	http.HandleFunc("/", http.HandlerFunc(pathHandler).ServeHTTP)
	// http.Handle("/", http.HandlerFunc(pathHandler))
	fmt.Println("In main now Starting the server on :3000...")
	// http.ListenAndServe(":3000", nil)
	http.ListenAndServe(":3000", somePathHandler)
	// http.ListenAndServe(":3000", router)
}

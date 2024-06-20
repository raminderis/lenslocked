package main

import (
	"fmt"
	"net/http"
	"time"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to Rami's awesome site!</h1>")
}

func main() {
	//mux := http.NewServeMux()
	//mux.HandleFunc("/", handlerFunc)
	//the above 2 lines OR below 1 line.
	http.HandleFunc("/", handlerFunc)
	fmt.Println("logic simulated by sleeping for 3 seconds")
	time.Sleep(3 * time.Second)
	fmt.Println("Starting the server on :3000...")
	err := http.ListenAndServe(":3000", nil) //nil means it passes to defaultservemux which is what is used by HandleFunc,
	// instead of you creating a mux in which case you put mux in place of nil
	if err != nil {
		panic(err)
	}
}

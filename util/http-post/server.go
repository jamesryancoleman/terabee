package main

import (
	"fmt"
	"log"
	"net/http"
)

func HandleLXL(w http.ResponseWriter, req *http.Request) {
	// for testing just print the body.
	log.Printf("\"/terabee/lxl\" endpoint called with method %s:\n%s", req.Method, req.Body)
}

// a default endpoint to confirm receipt of a http-post
func HandleDefaultEndpoint(w http.ResponseWriter, req *http.Request) {
	// for testing just print the body.
	log.Printf("\"/\" root endpoint called with method %s:\n%s", req.Method, req.Body)
}

func main() {
	// set loggin flags
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// define handlers
	http.HandleFunc("/terabee/lxl", HandleLXL)
	http.HandleFunc("/", HandleDefaultEndpoint)

	// start the server
	default_port := 8080
	log.Println("== Starting HTTP Server ==")
	log.Printf("== Listening on port %d ==", default_port)
	http.ListenAndServe(fmt.Sprintf(":%d", default_port), nil)
}

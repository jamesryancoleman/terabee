package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func ReadBody(req *http.Request) []byte {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic("error parsing body of request")
	}
	return body
}

<<<<<<<< HEAD:util/http/server.go
func HandleLXL(w http.ResponseWriter, req *http.Request) {
	// for testing just print the body.
	body := ReadBody(req)
	log.Printf("\"/terabee/lxl\" endpoint called with method %s:\n%s", req.Method, string(body))
========
func HandleFlow(w http.ResponseWriter, req *http.Request) {
	// for testing just print the body.
	body := ReadBody(req)
	log.Printf("\"/terabee/flow\" endpoint called with method %s:\n%s", req.Method, string(body))
>>>>>>>> docker:util/http-echo/server.go
}

// a default endpoint to confirm receipt of a http-post
func HandleDefaultEndpoint(w http.ResponseWriter, req *http.Request) {
	// for testing just print the body.
	body := ReadBody(req)
	log.Printf("\"/\" root endpoint called with method %s:\n%s", req.Method, string(body))
}

func main() {
	// set logging flags
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// define handlers
	http.HandleFunc("/terabee/flow", HandleFlow)
	http.HandleFunc("/", HandleDefaultEndpoint)

	// start the server
	default_port := 8080
	log.Println("== Starting HTTP Server ==")
	log.Printf("== Listening on port %d ==", default_port)
	http.ListenAndServe(fmt.Sprintf(":%d", default_port), nil)
}

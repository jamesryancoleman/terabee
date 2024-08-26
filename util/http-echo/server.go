package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type FlowPayload struct {
	UnixTS int64  `json:"at"`
	Serial string `json:"serial_number"`
	Flow   Flow   `json:"value"`
}

type Flow struct {
	In  int `json:"in"`
	Out int `json:"out"`
}

func (f *FlowPayload) GetTime() time.Time {
	return time.Unix(f.UnixTS, 0)
}

func ReadBody(req *http.Request) []byte {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic("error parsing body of request")
	}
	return body
}

// func HandleLXL(w http.ResponseWriter, req *http.Request) {
// 	// for testing just print the body.
// 	body := ReadBody(req)
// 	log.Printf("\"/terabee/lxl\" endpoint called with method %s:\n%s", req.Method, string(body))
// }

// func HandleFlow(w http.ResponseWriter, req *http.Request) {
// 	// for testing just print the body.
// 	body := ReadBody(req)
// 	log.Printf("\"/terabee/flow\" endpoint called with method %s:\n%s", req.Method, string(body))

// }

// a default endpoint to confirm receipt of a http-post
func HandleDefaultEndpoint(w http.ResponseWriter, req *http.Request) {
	// for testing just print the body.
	body := ReadBody(req)

	var flow FlowPayload
	err := json.Unmarshal(body, &flow)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Printf("%s %s %+v\n", flow.GetTime().Format(time.RFC3339), flow.Serial, flow.Flow)
	// log.Printf("\"/\" root endpoint called with method %s:\n%s", req.Method, string(body))

}

func main() {
	// set logging flags
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// define handlers
	// http.HandleFunc("/terabee/flow", HandleFlow)
	http.HandleFunc("/", HandleDefaultEndpoint)

	// start the server
	default_port := 8080
	log.Println("== Starting HTTP Server ==")
	log.Printf("== Listening on port %d ==", default_port)
	http.ListenAndServe(fmt.Sprintf(":%d", default_port), nil)
}

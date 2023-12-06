package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type FlowIn struct {
	TS     int64       `json:"at"`
	Serial string      `json:"serial_number"`
	Flow   FlowInValue `json:"value"`
}

type FlowInValue struct {
	In    int    `json:"in"`
	Out   int    `json:"out"`
	Reset string `json:"reset_period"`
}

func (f FlowInValue) GetOccupancy() int {
	return f.In - f.Out
}

func Flow2Frost(in FlowIn) FlowFrost {
	out := FlowFrost{}

	out.TS = time.Unix(in.TS, 0).Format(time.RFC3339)
	out.RTS = out.TS
	out.Result = in.Flow.GetOccupancy()
	out.Stream.Id, _ = strconv.Atoi(os.Args[2])

	return out
}

type FlowFrost struct {
	TS     string     `json:"phenomenonTime"`
	RTS    string     `json:"resultTime"`
	Result int        `json:"result"`
	Stream Datastream `json:"Datastream"`
}

type Datastream struct {
	Id int `json:"@iot.id"`
}

func main() {
	// start an http server on port
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

func HandleFlow(w http.ResponseWriter, req *http.Request) {
	// for testing just print the body.
	body := ReadBody(req)
	log.Printf("\"/terabee/flow\" endpoint called with method %s:\n%s", req.Method, string(body))

	// send to frost server
	client := &http.Client{}
	msg := readStdin()
	fmt.Printf("Posting to FROST Server:\n%s\n", msg)
	req, _ = http.NewRequest("POST", os.Args[1], bytes.NewBuffer([]byte(msg)))
	req.Header.Add("Authorization", "Basic "+basicAuth(os.Args[2], os.Args[3]))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response Status:", resp.Status)
	// done frost server

	// send to local mortar server
}

// a default endpoint to confirm receipt of a http-post
func HandleDefaultEndpoint(w http.ResponseWriter, req *http.Request) {
	// for testing just print the body.
	body := ReadBody(req)
	log.Printf("\"/\" root endpoint called with method %s:\n%s", req.Method, string(body))
}

func ReadBody(req *http.Request) []byte {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic("error parsing body of request")
	}
	return body
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func readStdin() string {
	reader := bufio.NewReader(os.Stdin)
	// TODO: this should be updated to be EOF
	text, _ := reader.ReadString('\n')
	return text
}

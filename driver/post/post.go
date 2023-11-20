package main

/*
usage:
		post URL USR PASSWRD

	STDIN – post reads from stdin to know what message to pass to the server
*/

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"os"
)

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

func main() {
	client := &http.Client{}
	msg := readStdin()
	log.Printf("Posting:\n%s\n", msg)
	req, _ := http.NewRequest("POST", os.Args[1], bytes.NewBuffer([]byte(msg)))
	req.Header.Add("Authorization", "Basic "+basicAuth(os.Args[2], os.Args[3]))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	log.Println("Response Status:", resp.Status)
}

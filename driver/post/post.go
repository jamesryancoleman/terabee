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
	"fmt"
	"net/http"
	"os"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func readStdin() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

func main() {
	fmt.Println("STARTED post...")
	client := &http.Client{}

	req, _ := http.NewRequest("POST", os.Args[1], bytes.NewBuffer([]byte(readStdin())))
	req.Header.Add("Authorization", "Basic "+basicAuth(os.Args[2], os.Args[3]))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response Status:", resp.Status)
}

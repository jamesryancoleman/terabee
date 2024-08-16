package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestReadCredentials(t *testing.T) {

	var c Credentials
	test_path := "credentials.json"
	err := ReadCredentials(test_path, &c)
	if err != nil {
		t.FailNow()
	}

	fmt.Printf("%+v\n", c)

}

func TestWriteToPG(t *testing.T) {
	var c Credentials
	test_path := "credentials.json"
	err := ReadCredentials(test_path, &c)
	if err != nil {
		t.FailNow()
	}

	now := time.Now()
	r := TSRow{
		Time:  now.Format(time.RFC3339),
		Value: 84698384, // ASCII: TEST
		Id:    "demo",
	}

	err = PGInsert(r, &c)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}

}

func TestDecodePayload(t *testing.T) {
	samplePayload := []byte{123, 34, 97, 116, 34, 58, 32, 49, 55, 50, 51, 56, 52, 48, 50, 57, 53, 44, 32, 34, 116, 121, 112, 101, 34, 58, 32, 34, 99, 111, 117, 110, 116, 101, 114, 115, 34, 44, 32, 34, 115, 101, 114, 105, 97, 108, 95, 110, 117, 109, 98, 101, 114, 34, 58, 32, 34, 98, 56, 50, 55, 101, 98, 98, 56, 54, 52, 100, 56, 34, 44, 32, 34, 118, 97, 108, 117, 101, 34, 58, 32, 123, 34, 105, 110, 34, 58, 32, 50, 44, 32, 34, 111, 117, 116, 34, 58, 32, 49, 44, 32, 34, 114, 101, 115, 101, 116, 95, 112, 101, 114, 105, 111, 100, 34, 58, 32, 34, 109, 97, 110, 117, 97, 108, 45, 101, 101, 53, 98, 49, 52, 55, 57, 34, 125, 32, 125}

	var fp FlowPayload
	err := json.Unmarshal(samplePayload, &fp)
	if err != nil {
		fmt.Println(err.Error())
	}

	// confirm Unix to RFC3339 works and print payload
	fmt.Printf("%s %+v\n", fp.GetTime().Format(time.RFC3339), fp)
}

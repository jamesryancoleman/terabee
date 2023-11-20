package main

/*
usage:
		convert_flow RAW_MSG ID
*/

import (
	"encoding/json"
	"fmt"
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

type FlowFrost struct {
	TS     string     `json:"phenomenonTime"`
	RTS    string     `json:"resultTime"`
	Result int        `json:"result"`
	Stream Datastream `json:"Datastream"`
}

type Datastream struct {
	Id int `json:"@iot.id"`
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

func main() {
	// read the command line arg
	msg := os.Args[1]

	// unmarshall the raw msg
	flowIn := FlowIn{}
	json.Unmarshal([]byte(msg), &flowIn)

	// convert to frost format
	flowOut := Flow2Frost(flowIn)
	fmt.Println("convert_flow is sending:")
	fmt.Println(flowOut)
	b, _ := json.Marshal(flowOut)

	fmt.Println(string(b))
}

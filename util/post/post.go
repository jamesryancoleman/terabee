package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// This struct is used to parse credentials from a JSON file
type Credentials struct {
	PgAddr     string `json:"pg_addr"`
	PgUser     string `json:"pg_user"`
	PgPassword string `json:"pg_password"`
	PgPort     int    `json:"pg_port"`
	PgDbName   string `json:"pg_dbname"`
	PgTable    string `json:"pg_table"`
	PgColumns  string `json:"column_names"`

	FrostAddr     string `json:"frost_addr"`
	FrostUser     string `json:"frost_user"`
	FrostPassword string `json:"frost_password"`
}

func ReadCredentials(path string, credPointer *Credentials) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(raw, credPointer)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

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
	out.Stream.Id = 1

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

type MortarRow struct {
	TS    string
	value float32
	id    int64
}

type TSRow struct {
	TS    string
	Value float32
	Id    string
}

func (m *TSRow) Prepare() []interface{} {
	args := make([]interface{}, 3)
	args[0] = m.TS
	args[1] = m.Value
	args[2] = m.Id
	return args
}

// func (m *MortarRow) Prepare() []interface{} {
// 	args := make([]interface{}, 3)
// 	args[0] = m.TS
// 	args[1] = m.value
// 	args[2] = m.id
// 	return args
// }

// globally accessible
var ctx *PGContext
var credentials *Credentials

func main() {
	// start an http server on port
	// set logging flags
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// parse postgres flags
	// command line stuff
	host := flag.String("host", "node1.local", "host of the postgres db")
	port := flag.Int("port", 5432, "port of the postgres instance")
	user := flag.String("user", "postgres", "username with which to access the db")
	password := flag.String("pswd", "password", "password of the username provided")
	dbname := flag.String("db", "data", "name of the database to access")
	table := flag.String("table", "data", "name of the table to query")
	colStr := flag.String("cols", "time,value,id", "column names seperated by commas")
	credPath := flag.String("f", "credentials.json", "location of the credentials json")
	flag.Parse()

	// parse the credentials from the file
	err := ReadCredentials(*credPath, credentials)
	if err != nil {
		log.Printf(err.Error())
	}

	// pass context as a PGContext struct
	ctx = &PGContext{
		host:     *host,
		port:     *port,
		user:     *user,
		password: *password,
		dbname:   *dbname,
		table:    *table,
		colStr:   *colStr,
	}

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

	// assume flow format and convert to frost AND mortar
	// FROST-server conversion

	// unmarshall the raw msg
	flowIn := FlowIn{}
	json.Unmarshal([]byte(body), &flowIn)

	// convert to frost format
	flowOut := Flow2Frost(flowIn)
	frostBytes, _ := json.Marshal(flowOut)

	// credentials
	frost_url := "http://chaosbox.princeton.edu/frost/v1.1/Observations"
	frost_user := "chaos"
	frost_password := "Zer0exergy"

	// send to frost server
	client := &http.Client{}
	fmt.Printf("Posting to FROST Server:\n%s\n", frostBytes)
	req, _ = http.NewRequest("POST", frost_url, bytes.NewBuffer(frostBytes))
	req.Header.Add("Authorization", "Basic "+basicAuth(frost_user, frost_password))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("unable to post to FROST server: %s", err.Error())
	}

	fmt.Println("Response Status:", resp.Status)
	// end FROST-server

	// send to local mortar server
	// convert FlowOut to MortarRow
	// PGInsert(TSRow{
	// 	TS:    flowOut.TS,
	// 	Value: float32(flowOut.Result),
	// 	Id:    "DO SOMETHING HERE",
	// })

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

type PGContext struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
	table    string
	colStr   string
}

// func PGInsert(m MortarRow) error {
// 	columns := strings.Split(ctx.colStr, ",")
// 	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 		ctx.host, ctx.port, ctx.user, ctx.password, ctx.dbname)

// 	// open the database client
// 	db, err := sql.Open("postgres", conn)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	statement := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`,
// 		ctx.table, strings.Join(columns, ", "), PlaceHolders(len(columns)))
// 	// fmt.Println(statement)

// 	values := m.Prepare()
// 	_, err = db.Exec(statement, values...)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return err
// }

func PGInsert(m TSRow, c *Credentials) error {
	columns := strings.Split(c.PgColumns, ",")
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.PgAddr, c.PgPort, c.PgUser, c.PgPassword, c.PgDbName)

	// open the database client
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// compose sql injection statement
	statement := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`,
		c.PgTable, strings.Join(columns, ", "), PlaceHolders(len(columns)))
	// fmt.Println(statement)

	values := m.Prepare()
	_, err = db.Exec(statement, values...)
	if err != nil {
		panic(err)
	}

	return err
}

func PlaceHolders(numPlaceholders int) string {
	placeholders := ""
	for i := 1; i <= numPlaceholders; i++ {
		if i == numPlaceholders {
			placeholders += fmt.Sprintf("$%d", i)
			break
		}
		placeholders += fmt.Sprintf("$%d, ", i)
	}
	return placeholders
}

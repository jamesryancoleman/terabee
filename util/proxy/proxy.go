package main

import (
	"database/sql"
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

// globals
var ColumnNames []string
var credentials Credentials

// This struct is used to parse credentials from a JSON file
type Credentials struct {
	PgAddr     string `json:"pg_addr"`
	PgUser     string `json:"pg_user"`
	PgPassword string `json:"pg_password"`
	PgPort     int    `json:"pg_port"`
	PgDbName   string `json:"pg_dbname"`
	PgTable    string `json:"pg_table"`
	PgColumns  string `json:"column_names"`
}

// payload from Terabee LXL Flow message
type FlowPayload struct {
	UnixTS int64  `json:"at"`
	Serial string `json:"serial_number"`
	Flow   Flow   `json:"value"`
}

// value object from Terabee LXL flow payload
type Flow struct {
	In  int `json:"in"`
	Out int `json:"out"`
}

func (f *FlowPayload) GetTime() time.Time {
	return time.Unix(f.UnixTS, 0)
}

// columns of a minimal Timescale table
type TSRow struct {
	Time  string
	Value float32
	Id    string
}

// converts TSRow into a slice of interfaces
func (m *TSRow) ToInterfaces() []interface{} {
	args := make([]interface{}, 3)
	args[0] = m.Time
	args[1] = m.Value
	args[2] = m.Id
	return args
}

// load postgres credentials from a JSON file
func ReadCredentials(path string, credPointer *Credentials) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(raw, credPointer)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		// set the column names variable
		ColumnNames = strings.Split(credPointer.PgColumns, ",")
	}
	// fmt.Printf("%+v\n", *credPointer)
	return err
}

// insert to postgres
func PGInsert(m TSRow, c *Credentials) error {
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
		c.PgTable, strings.Join(ColumnNames, ", "), PlaceHolders(len(ColumnNames)))

	values := m.ToInterfaces()
	_, err = db.Exec(statement, values...)
	if err != nil {
		panic(err)
	}

	return err
}

// helper function for formatting sql insert
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

func HandleDefaultEndpoint(w http.ResponseWriter, req *http.Request) {
	body := ReadBody(req)

	var flow FlowPayload
	err := json.Unmarshal(body, &flow)
	if err != nil {
		log.Println(err.Error())
	}

	// uncomment to validate payload
	fmt.Printf("%s %s %+v\n", flow.GetTime().Format(time.RFC3339), flow.Serial, flow.Flow)

	// insert the in and out values separately
	in := TSRow{
		Time:  flow.GetTime().Format(time.RFC3339),
		Value: float32(flow.Flow.In),
		Id:    fmt.Sprintf("%s.%s", flow.Serial, "In"),
	}
	PGInsert(in, &credentials)
	if err != nil {
		log.Println(err.Error())
	}

	out := TSRow{
		Time:  flow.GetTime().Format(time.RFC3339),
		Value: float32(flow.Flow.Out),
		Id:    fmt.Sprintf("%s.%s", flow.Serial, "Out"),
	}
	PGInsert(out, &credentials)
	if err != nil {
		log.Println(err.Error())
	}

}

func ReadBody(req *http.Request) []byte {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic("error parsing body of request")
	}
	return body
}

func main() {
	// set logging flags
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// command line stuff
	credPath := flag.String("f", "credentials.json", "location of the credentials json")
	flag.Parse()

	// parse the credentials from the file
	fmt.Printf("%v %s\n", credPath, *credPath)
	err := ReadCredentials(*credPath, &credentials)
	if err != nil {
		log.Println(err.Error())
	}

	// define handlers
	http.HandleFunc("/", HandleDefaultEndpoint)

	// start the server
	default_port := 8080
	log.Println("== Starting HTTP Server ==")
	log.Printf("== Listening on port %d ==", default_port)
	http.ListenAndServe(fmt.Sprintf(":%d", default_port), nil)
}

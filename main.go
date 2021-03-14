package main

import (
	"bytes"
	"crypto/subtle"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/olekukonko/tablewriter"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var host = ":8080"
var conn = "root:password@tcp(127.0.0.1:3306)/"
var username = "admin"
var password = "admin"
var title = "MySQL Console"

type Env struct {
	Title string
}

//go:embed templates
var templateFiles embed.FS

//go:embed static
var staticFiles embed.FS

var db *sql.DB
var validDB = false

func jsonResponse(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(i)
}


func errorConnecting(err error) {
	validDB = false
	log.Printf("Error connecting to database: %s", err)
	log.Printf("Trying again in 1 minute...")
	time.Sleep(1 * time.Minute)
}

func keepMySQLAlive() {
	mysql.SetLogger(log.New(ioutil.Discard, "", 0))

	for {
		log.Println("Connecting to database...")

		var err error
		db, err = sql.Open("mysql", conn)
		if err != nil {
			errorConnecting(err)
			continue
		}
		err = db.Ping()
		if err != nil {
			errorConnecting(err)
			continue
		}
		log.Println("Connected to database!")
		validDB = true

		for {
			err := db.Ping()
			if err != nil {
				errorConnecting(err)
				break
			}
			time.Sleep(2 * time.Second)
		}
	}
}

type ErrResp struct {
	Error string `json:"err"`
}

type MsgResp struct {
	Message string `json:"message"`
}

func isExec(statement string) bool {
	statement = strings.TrimSpace(strings.ToLower(statement))
	return strings.HasPrefix(statement, "update") || strings.HasPrefix(statement, "insert") || strings.HasPrefix(statement, "delete")
}

func handleQuery(w http.ResponseWriter, statement string) {
	rows, err := db.Query(statement)
	if err != nil {
		jsonResponse(w, ErrResp{Error: err.Error()})
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		jsonResponse(w, ErrResp{
			Error: err.Error(),
		})
		return
	}

	buffer := bytes.NewBufferString("")
	table := tablewriter.NewWriter(buffer)
	table.SetHeader(columns)
	table.SetAutoFormatHeaders(false)

	if len(columns) == 0 {
		jsonResponse(w, MsgResp{
			Message: "OK",
		})
		return
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i, _ := range columns {
			values[i] = new(sql.RawBytes)
		}



		err := rows.Scan(values...)
		if err != nil {
			jsonResponse(w, ErrResp{Error: err.Error()})
			return
		}

		valueStrings := []string{}
		for _, value := range values {
			valueStrings = append(valueStrings, fmt.Sprintf("%s", value))
		}
		table.Append(valueStrings)
	}
	table.Render()
	jsonResponse(w, MsgResp{
		Message: buffer.String(),
	})
}

func handleExec(w http.ResponseWriter, statement string) {
	res, err := db.Exec(statement)
	if err != nil {
		jsonResponse(w, ErrResp{
			Error: err.Error(),
		})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		jsonResponse(w, ErrResp{
			Error: err.Error(),
		})
		return
	}

	message := fmt.Sprintf("%d row(s) affected", rowsAffected)
	jsonResponse(w, MsgResp{
		Message: message,
	})
}

func basicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="login"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}

func main() {
	if val := os.Getenv("HOST"); val != "" {
		host = val
	}

	if val := os.Getenv("DB_CONN"); val != "" {
		conn = val
	}

	if val := os.Getenv("USERNAME"); val != "" {
		username = val
	}

	if val := os.Getenv("PASSWORD"); val != "" {
		password = val
	}

	if val := os.Getenv("TITLE"); val != "" {
		title = val
	}

	go keepMySQLAlive()

	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)
	http.Handle("/static/", fs)

	http.HandleFunc("/", basicAuth(func(w http.ResponseWriter, r *http.Request) {
		// t, err := template.ParseFS(templateFiles, "templates/index.html")
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			return
		}

		t.Execute(w, Env{
			Title: title,
		})
	}))

	http.HandleFunc("/query", basicAuth(func(w http.ResponseWriter, r *http.Request) {
		if !validDB || db == nil {
			jsonResponse(w, ErrResp{Error: "Not connected to database"})
			return
		}

		statement := r.URL.Query().Get("statement")
		if isExec(statement) {
			handleExec(w, statement)
		} else {
			handleQuery(w, statement)
		}
	}))

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if !validDB || db == nil {
			jsonResponse(w, ErrResp{Error: "Not connected to database"})
			return
		}

		jsonResponse(w, MsgResp{Message: "Connected to the database"})
	})

	log.Printf("Starting server at: %s", host)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		panic(err)
		return
	}
}
package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"github.com/felixge/httpsnoop"
	"github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	mux := http.NewServeMux()

	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)
	mux.Handle("/static/", fs)

	mux.HandleFunc("/login", login)
	mux.Handle("/", verify(index))
	mux.Handle("/query", verify(query))
	mux.Handle("/status", verify(status))

	loggingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(mux, w, r)
		log.Printf("%s %s %d %s",
			r.Method,
			r.URL,
			m.Code,
			m.Duration)
	})

	log.Printf("Starting server at: %s", host)
	err := http.ListenAndServe(host, loggingHandler)
	if err != nil {
		panic(err)
		return
	}
}
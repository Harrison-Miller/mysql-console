# MySQL Console

A simple web console for querying your mysql database.


### Features

* mysql-client like console
* Simple login page and auth protection
* Database connection status
* Rerun query from history

![Screenshot](media/screenshot.png?raw=true)

## Requirements 

Go 1.16

## Build & Run

`go build . -o mysql-console`
`go run . `

## Run with Docker

`docker run --rm -it -p8080:8080 harrisonmiller/mysql-console`

## Environment Variables

| Variable | Description | Default | Notes |
| -------- | ----------- | ------- | ----- |
| HOST | address for the web server | `:8080` | |
| DB_CONN | db connection string | `root:password@tcp(127.0.0.1:3306)/`| [uses DSN format](https://github.com/go-sql-driver/mysql#dsn-data-source-name) |
| USERNAME | username | `admin` | |
| PASSWORD | password as bcrypt hash | `admin` | [format](https://en.wikipedia.org/wiki/Bcrypt#Description) |
| TITLE | custom page title | `MySQL Console` | |
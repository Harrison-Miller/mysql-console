# MySQL Querier

A simple web page to run queries against your mysql database.

**Note**: This was created as a fun evening project. You should not use this for any serious purpose. I am not responsible for any harm that may come of using this software.

![Screenshot](media/screenshot.png?raw=true)

## Requirements 

Go 1.16

## Build

`go build . -o mysql-querier`

## Docker

`docker build . -tmysql-querier`

`docker run --rm -it -p8080:8080 mysql-querier`

## Environment Variables

| Variable | Description | Default | Notes |
| -------- | ----------- | ------- | ----- |
| HOST | host for webserver | `:8080` | |
| DB_CONN | db connection string | `root:password@tcp(127.0.0.1:3306)/`| [uses DSN format](https://github.com/go-sql-driver/mysql#dsn-data-source-name) |
| USERNAME | basic auth username | `admin` | |
| PASSWORD | basic auth password | `admin` | |
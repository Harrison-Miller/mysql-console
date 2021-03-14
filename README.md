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

* **HOST**: the host for the webserver - default: `:8080`
* **DB_CONN**: the db connection string
* **USERNAME**: basic auth username - default: `admin`
* **PASSWORD**: basic auth password - default: `password`

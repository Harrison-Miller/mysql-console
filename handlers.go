package main

import (
	"html/template"
	"net/http"
)

var indexTemplate *template.Template

func init() {
	indexTemplate = template.Must(template.ParseFS(templateFiles, "templates/base.html", "templates/index.html"))
}

func status(w http.ResponseWriter, r *http.Request) {
	if !validDB || db == nil {
		jsonResponse(w, ErrResp{Error: "Not connected to database"})
		return
	}

	jsonResponse(w, MsgResp{Message: "Connected to the database"})
}

func index(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, Env{
		Title: title,
	})
}

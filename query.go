package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"net/http"
	"strings"
)

func query(w http.ResponseWriter, r *http.Request) {
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


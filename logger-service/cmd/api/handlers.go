package main

import (
	"net/http"

	"log-service/data"
)

type JsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json to var
	var requestPayload JsonPayload
	_ = app.readJson(w, r, &requestPayload)

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJson(w, http.StatusAccepted, resp)
}

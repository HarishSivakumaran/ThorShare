package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error`
	Message string `json:"message`
	Data    any    `json:"data, omitempty"`
}

func (app *Config) readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 //1 MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body must have only one JSON value")
	}

	return nil
}

func (app *Config) writeJson(w http.ResponseWriter, status int, data any, header ...http.Header) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(header) > 0 {
		for key, val := range header[0] {
			w.Header()[key] = val
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(res)

	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse

	payload.Error = true
	payload.Message = err.Error()

	return app.writeJson(w, statusCode, payload)

}

func (app *Config) readJsonToStruct(w http.ResponseWriter, jsonData io.ReadCloser, data any) error {
	maxBytes := 1048576 //1 MB

	jsonData = http.MaxBytesReader(w, jsonData, int64(maxBytes))

	dec := json.NewDecoder(jsonData)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body must have only one JSON value")
	}

	return nil
}

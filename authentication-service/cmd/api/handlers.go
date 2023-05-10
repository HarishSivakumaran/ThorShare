package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

)

const (
	logServiceUrl = "http://logger-service/log"
)

func (app *Config) authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &requestPayload)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.errorJson(w, errors.New("Invalid Credentials"), http.StatusBadRequest)
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		app.errorJson(w, errors.New("Invalid Credentials"), http.StatusBadRequest)
	}

	//log authentication
	err = app.LogRequest("Authentication", fmt.Sprintf("%s logged in ", user.Email))
	if err != nil {
		log.Println("Something went wrong while logging")
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in used %s", user.Email),
		Data:    user,
	}

	app.writeJson(w, http.StatusAccepted, payload)

}

func (app *Config) LogRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jwt/pkg/rest"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var authToken string

// Testing real-life user flow on seed data and live server.
func Test_UserFlow(t *testing.T) {
	// Prepare db.
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sqlx.Open("mysql", os.Getenv("MYSQL_DSN_TEST"))
	if err != nil {
		log.Fatal("failed to connect database")
	}
	_, err = db.Exec("TRUNCATE TABLE vote;")
	if err != nil {
		log.Fatal("could not truncate table vote")
	}

	// Login.
	payload := map[string]string{
		"email":    "user@example.com",
		"password": "password123",
	}
	login(t, payload)

	// Success on upvote user.
	upvote1 := map[string]string{
		"from_user": "e0e5ba28-19fc-4c65-8692-f61266608e3e",
		"vote":      "1",
	}
	assert.Equal(t, http.StatusOK, upvoteUser(t, "e0e5ba28-19fc-4c65-8692-f61266608d4n", upvote1).StatusCode)

	// Error on the same upvote in less than 24 hrs.
	upvote2 := map[string]string{
		"from_user": "e0e5ba28-19fc-4c65-8692-f61266608e3e",
		"vote":      "1",
	}
	assert.Equal(t, http.StatusInternalServerError, upvoteUser(t, "e0e5ba28-19fc-4c65-8692-f61266608e3e", upvote2).StatusCode)

	// Error on the same user upvote.
	upvote3 := map[string]string{
		"from_user": "e0e5ba28-19fc-4c65-8692-f61266608e3e",
		"vote":      "1",
	}
	assert.Equal(t, http.StatusInternalServerError, upvoteUser(t, "e0e5ba28-19fc-4c65-8692-f61266608e3e", upvote3).StatusCode)

	// Error on vote validation.
	upvote4 := map[string]string{
		"from_user": "xxx",
		"vote":      "10",
	}
	assert.Equal(t, http.StatusBadRequest, upvoteUser(t, "xxx", upvote4).StatusCode)

	// View user with upvotes.
	assert.Equal(t, http.StatusOK, viewUser(t, "e0e5ba28-19fc-4c65-8692-f61266608d4n").StatusCode)
}

func login(t *testing.T, payload map[string]string) {

	resp, err := sendRequest("POST", "http://localhost:8080/login", nil, payload)
	if err != nil {
		t.Error(err)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(respBody))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var r rest.Response
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		t.Error(err)
	}

	// Retrive token from response.
	if loginResp, ok := r.Payload.(map[string]interface{}); ok {
		authToken = loginResp["token"].(string)
	} else {
		t.Error("Could not retrieve token from response.")
	}

	assert.NotEmpty(t, authToken)
}

func upvoteUser(t *testing.T, uuid string, payload map[string]string) *http.Response {

	resp, err := sendRequest("POST", "http://localhost:8080/user/"+uuid+"/vote", &authToken, payload)
	if err != nil {
		t.Error(err)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(respBody))

	return resp

}
func viewUser(t *testing.T, uuid string) *http.Response {

	resp, err := sendRequest("GET", "http://localhost:8080/user/"+uuid, &authToken, nil)
	if err != nil {
		t.Error(err)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(respBody))

	return resp

}

// sendRequest sends an HTTP request with the given method and JSON payload.
func sendRequest(method string, url string, token *string, payload interface{}) (*http.Response, error) {
	// Marshal the payload into JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)

	// Create the HTTP request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("Authorization", *token)
	}

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

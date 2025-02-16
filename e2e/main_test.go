package main_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"
)

const MERCH_APP_URL = "http://app:8080"

func auth(t *testing.T, username, password string) string {
    authRequest := map[string]string{
        "username": username,
        "password": password,
    }
    body, err := json.Marshal(authRequest)
    assert.NoError(t, err)

    resp, err := http.Post(MERCH_APP_URL + "/api/auth", "application/json", bytes.NewBuffer(body))
    assert.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode, "expected successful status")

    var authResponse struct {
        Token string `json:"token"`
    }
    err = json.NewDecoder(resp.Body).Decode(&authResponse)
    assert.NoError(t, err)

    return authResponse.Token
}
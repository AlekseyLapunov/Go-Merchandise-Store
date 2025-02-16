package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlekseyLapunov/Go-Merchandise-Store/e2e/structs"
)

func TestSendCoin_Valid(t *testing.T) {
    type Emps struct {
        client      *http.Client
        token       string
        coinsToSend int
    }

    emps := []Emps{
        {&http.Client{}, auth(t, "test_user3", "some_password"), 111},
        {&http.Client{}, auth(t, "test_user4", "some_password"), 587},
    }

    // --- checking that both employees have 1000 coins on start
    for _, e := range emps {
        _, infoResponse := getInfo(t, e.client, e.token)
        assert.Equal(t, infoResponse.Coins, 1000)
    }

    // --- first employee sends coins
    e := emps[0]
    sendCoin(t, e.client, e.token, "test_user4", e.coinsToSend)

    // --- second employee also sends coins
    e = emps[1]
    sendCoin(t, e.client, e.token, "test_user3", e.coinsToSend)

    // --- both users are checking their balance after those transfers
    for i, e := range emps {
        _, infoResponse := getInfo(t, e.client, e.token)

        expectedAmount := 1000 - e.coinsToSend
        if i == 0 {
            expectedAmount += emps[1].coinsToSend
        } else {
            expectedAmount += emps[0].coinsToSend
        }

        assert.Equal(t, infoResponse.Coins, expectedAmount)
    }
}

func TestSendCoin_Verbose(t *testing.T) {
    type Emps struct {
        client      *http.Client
        token       string
    }

    emps := []Emps{
        {&http.Client{}, auth(t, "test_user90", "some_password")},
        {&http.Client{}, auth(t, "test_user45", "some_password")},
        {&http.Client{}, auth(t, "test_user34", "some_password")},
    }

    // --- first employee sends coins to other employees
    e := emps[0]
    sendCoin(t, e.client, e.token, "test_user45", 421)
    sendCoin(t, e.client, e.token, "test_user34", 123)

    // --- second employee sends coins to other employees
    e = emps[1]
    sendCoin(t, e.client, e.token, "test_user90", 10)
    sendCoin(t, e.client, e.token, "test_user34", 5)

    // --- third employee sends coins to other employees
    e = emps[2]
    sendCoin(t, e.client, e.token, "test_user90", 111)
    sendCoin(t, e.client, e.token, "test_user45", 222)

    // --- checking exchange history
    e = emps[0]
    checkHistory(t, e.client, e.token, "test_user45", 10)

    e = emps[1]
    checkHistory(t, e.client, e.token, "test_user90", 421)

    e = emps[2]
    checkHistory(t, e.client, e.token, "test_user45", 5)
}

func sendCoin(t *testing.T, client *http.Client, token, toUser string, amount int) *http.Response {
    body := fmt.Sprintf(`{"toUser": "%s", "amount": %d}`, toUser, amount)

    req, err := http.NewRequest("POST", MERCH_APP_URL + "/api/sendCoin", strings.NewReader(body))
    assert.NoError(t, err)
    req.Header.Set("Authorization", "BearerAuth " + token)

    resp, err := client.Do(req)
    assert.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode, "expected successful status")

    return resp
}

func getInfo(t *testing.T, client *http.Client, token string) (*http.Response, structs.InfoResponse) {
    req, err := http.NewRequest("GET", MERCH_APP_URL + "/api/info", nil)
    assert.NoError(t, err)
    req.Header.Set("Authorization", "BearerAuth " + token)

    resp, err := client.Do(req)
    assert.NoError(t, err)
    defer resp.Body.Close()
    assert.Equal(t, http.StatusOK, resp.StatusCode, "expected successful status")

    var infoResponse structs.InfoResponse

    err = json.NewDecoder(resp.Body).Decode(&infoResponse)
    assert.NoError(t, err, "can't decode infoResponse")

    return resp, infoResponse
}

func checkHistory(t *testing.T, client *http.Client, token string, fromUser string, amount int) {
    _, info := getInfo(t, client, token)

    isFound    := false
    realAmount := 0
    for _, entry := range info.CoinHistory.Received {
        if entry.FromUser == fromUser {
            isFound    = true
            realAmount = entry.Amount
        }
    }

    assert.Equal(t, isFound, true)
    assert.Equal(t, realAmount, amount)
    assert.Equal(t, len(info.CoinHistory.Received), 2)
    assert.Equal(t, len(info.CoinHistory.Sent), 2)
}

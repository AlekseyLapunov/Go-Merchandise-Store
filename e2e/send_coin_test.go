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
        req, err := http.NewRequest("GET", MERCH_APP_URL + "/api/info", nil)
        assert.NoError(t, err)
        req.Header.Set("Authorization", "BearerAuth " + e.token)
    
        resp, err := e.client.Do(req)
        assert.NoError(t, err)
        defer resp.Body.Close()
        assert.Equal(t, http.StatusOK, resp.StatusCode, "expected successful status")

        var infoResponse structs.InfoResponse

        err = json.NewDecoder(resp.Body).Decode(&infoResponse)
        assert.NoError(t, err, "can't decode infoResponse")

        assert.Equal(t, infoResponse.Coins, 1000)
    }

    // --- first employee sends coins
    e := emps[0]
    {
        body := fmt.Sprintf(`{"toUser": "test_user4", "amount": %d}`, e.coinsToSend)

        req, err := http.NewRequest("POST", MERCH_APP_URL + "/api/sendCoin", strings.NewReader(body))
        assert.NoError(t, err)
        req.Header.Set("Authorization", "BearerAuth " + e.token)
    
        resp, err := e.client.Do(req)
        assert.NoError(t, err)
        defer resp.Body.Close()

        assert.Equal(t, http.StatusOK, resp.StatusCode, "expected successful status")
    }

    // --- second employee also sends coins
    e = emps[1]
    {
        body := fmt.Sprintf(`{"toUser": "test_user3", "amount": %d}`, e.coinsToSend)

        req, err := http.NewRequest("POST", MERCH_APP_URL + "/api/sendCoin", strings.NewReader(body))
        assert.NoError(t, err)
        req.Header.Set("Authorization", "BearerAuth " + e.token)
    
        resp, err := e.client.Do(req)
        assert.NoError(t, err)
        defer resp.Body.Close()

        assert.Equal(t, http.StatusOK, resp.StatusCode, "expected successful status")
    }

    // --- both users are checking their balance after those transfers
    for i, e := range emps {
        req, err := http.NewRequest("GET", MERCH_APP_URL + "/api/info", nil)
        assert.NoError(t, err)
        req.Header.Set("Authorization", "BearerAuth " + e.token)
    
        resp, err := e.client.Do(req)
        assert.NoError(t, err)
        defer resp.Body.Close()
        assert.Equal(t, http.StatusOK, resp.StatusCode, "expected successful status")

        var infoResponse structs.InfoResponse

        err = json.NewDecoder(resp.Body).Decode(&infoResponse)
        assert.NoError(t, err, "can't decode infoResponse")

        expectedAmount := 1000 - e.coinsToSend
        if i == 0 {
            expectedAmount += emps[1].coinsToSend
        } else {
            expectedAmount += emps[0].coinsToSend
        }

        assert.Equal(t, infoResponse.Coins, expectedAmount)
    }
}
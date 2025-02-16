package main_test

import (
    "encoding/json"
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/AlekseyLapunov/Go-Merchandise-Store/e2e/structs"
)

func TestBuyMerch_Valid(t *testing.T) {
    client := &http.Client{}

    token := auth(t, "test_user1", "some_password")

    // --- checking that we have 1000 coins on start
    {
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

        assert.Equal(t, infoResponse.Coins, 1000)
    }

    // --- buying powerbanks
    const powerbanks = 3
    const cost = 200
    {
        req, err := http.NewRequest("GET", MERCH_APP_URL + "/api/buy/powerbank", nil)
        assert.NoError(t, err)
        req.Header.Set("Authorization", "BearerAuth " + token)
    
        for i := 0; i < powerbanks; i++ {
            resp, err := client.Do(req)
            assert.NoError(t, err)
            defer resp.Body.Close()
    
            assert.Equal(t, http.StatusOK, resp.StatusCode, "expected successful status")
        }
    }

    // --- checking info after these purchases
    {
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

        assert.Equal(t, infoResponse.Inventory[0].Type, "powerbank")
        assert.Equal(t, infoResponse.Inventory[0].Quantity, 3)
        assert.Equal(t, infoResponse.Coins, 1000 - powerbanks*cost)
    }
}
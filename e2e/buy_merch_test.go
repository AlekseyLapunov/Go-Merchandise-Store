package main_test

import (
	"fmt"
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
        resp, info := getInfo(t, client, token)

        assert.Equal(t, http.StatusOK, resp.StatusCode, "expected successful status")
        assert.Equal(t, info.Coins, 1000)
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
        resp, info := getInfo(t, client, token)

        assert.Equal(t, http.StatusOK, resp.StatusCode, "expected successful status")
        assert.Equal(t, info.Inventory[0].Type, "powerbank")
        assert.Equal(t, info.Inventory[0].Quantity, 3)
        assert.Equal(t, info.Coins, 1000 - powerbanks*cost)
    }
}

func TestBuyMerch_WrongToken(t *testing.T) {
    client := &http.Client{}

    falseToken := "l;lgldfg.jkkasdlk.690gflkk5"

    // --- trying to buy something
    req, err := http.NewRequest("GET", MERCH_APP_URL + "/api/buy/hoody", nil)
    assert.NoError(t, err)
    req.Header.Set("Authorization", "BearerAuth " + falseToken)

    resp, err := client.Do(req)
    assert.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "expected unauthorized")
}

func TestBuyMerch_InsufficientCoins(t *testing.T) {
    client := &http.Client{}

    token := auth(t, "test_user5", "some_password")

    // --- trying to buy something expensive a lot of times
    const hoodies = 6
    req, err := http.NewRequest("GET", MERCH_APP_URL + "/api/buy/hoody", nil)
    assert.NoError(t, err)
    req.Header.Set("Authorization", "BearerAuth " + token)

    var resp *http.Response
    for i := 0; i < hoodies; i++ {
        resp, err := client.Do(req)
        assert.NoError(t, err)
        defer resp.Body.Close()
    }

    if resp == nil {
        return
    }
    
    // --- last reponse should contain status code 400
    assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "expected bad request")
}

func TestBuyMerch_Verbose(t *testing.T) {
    client := &http.Client{}

    token := auth(t, "test_user8", "some_password")

    itemsToBuy := []structs.InventoryItem {
        {"hoody",    1},  // 300
        {"pen",      15}, // 150
        {"umbrella", 1},  // 200
        {"socks",    2},  // 20
        {"wallet",   2},  // 50
        {"book",     4},  // 200
        {"cup",      1},  // 20

    }
    resultCost := 990

    // --- buying a bunch of different items
    for _, item := range itemsToBuy {
        req, err := http.NewRequest("GET", MERCH_APP_URL + fmt.Sprintf("/api/buy/%s", item.Type), nil)
        assert.NoError(t, err)
        req.Header.Set("Authorization", "BearerAuth " + token)
    
        for i := 0; i < item.Quantity; i++ {
            resp, err := client.Do(req)
            assert.NoError(t, err)
            defer resp.Body.Close()

            assert.Equal(t, http.StatusOK, resp.StatusCode, "expected successful operation")
        }
    }
    
    // --- checking info after these purchases
    {
        resp, info := getInfo(t, client, token)

        assert.Equal(t, http.StatusOK, resp.StatusCode, "expected successful status")

        nameToQuantity := make(map[string]int)
        for _, item := range itemsToBuy {
            nameToQuantity[item.Type] += item.Quantity
        }

        for _, item := range info.Inventory {
            quantity := nameToQuantity[item.Type]
            assert.Equal(t, item.Quantity, quantity)
        }

        assert.Equal(t, info.Coins, 1000 - resultCost)
    }
}

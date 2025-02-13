package entity

type InfoResponse struct {
    Coins       int             `json:"coins"`
    Inventory   []InventoryItem `json:"inventory"`
    CoinHistory CoinHistory     `json:"coinHistory"`
}

type InventoryItem struct {
    Type     string `json:"type"`
    Quantity int    `json:"quantity"`
}

type CoinHistory struct {
    Received []Transaction `json:"received"`
    Sent     []Transaction `json:"sent"`
}

type Transaction struct {
    FromUser string `json:"fromUser"`
    ToUser   string `json:"toUser"`
    Amount   int    `json:"amount"`
}

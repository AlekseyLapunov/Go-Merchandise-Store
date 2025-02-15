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
    Received []RecvEntry `json:"received"`
    Sent     []SentEntry `json:"sent"`
}

type RecvEntry struct {
    FromUser string `json:"fromUser"`
    Amount   int    `json:"amount"`
}

type SentEntry struct {
    ToUser string `json:"toUser"`
    Amount int    `json:"amount"`
}


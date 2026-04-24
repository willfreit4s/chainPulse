package entity

import "time"

type TokenMetric struct {
	ID               string     `json:"id"`
	TokenID          string     `json:"token_id"`
	Token_address    string     `json:"token_address"`
	Price            float64    `json:"price"`
	Price_change_24h float64    `json:"price_change_24h"`
	Volume_24h       float64    `json:"volume_24h"`
	Market_cap       float64    `json:"market_cap"`
	Holders          int64      `json:"holders"`
	Transactions_24h int64      `json:"transactions_24h"`
	Timestamp        time.Time  `json:"timestamp"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
	Deleted          bool       `json:"deleted"`
}

func (tm *TokenMetric) GetTokenMetric() *TokenMetric {
	return tm
}

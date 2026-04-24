package entity

import "time"

type Token struct {
	ID          string     `json:"id"`
	Address     string     `json:"address"`
	Name        string     `json:"name"`
	Symbol      string     `json:"symbol"`
	Decimals    uint8      `json:"decimals"`
	TotalSupply float64    `json:"total_supply"`
	Chain       string     `json:"chain"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Deleted     bool       `json:"deleted"`
}

func (t *Token) GetToken() *Token {
	return t
}

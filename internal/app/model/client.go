package model

import "time"

//Client struct represent client data
type Client struct {
	ID               string           `json:"client_id"`
	ClientName       string           `json:"client_name"`
	ClientsRefTokens []ClientRefToken `json:"-"`
}

//ClientRefToken struct
type ClientRefToken struct {
	SessionID string    `json:"-"`
	RefToken  string    `json:"refToken"`
	ExpIn     time.Time `json:"exp_in"`
	CreatedAt time.Time `json:"-"`
}

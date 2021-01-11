package model

import "time"

type Token struct {
	ExpIn time.Time `json:"exp_in,omitempty"`
	Token string    `json:"token,omitempty"`
}

type Identity struct {
	UserID   string `json:"user_id,omitempty"`
	UserName string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	//
	AuthToken Token `json:"auth_token,omitempty"`
	RefToken  Token `json:"ref_token,omitempty"`
}

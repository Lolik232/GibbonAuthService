//Package model is package of data models
package model

import (
	"time"
)

//User struct represent user data
type User struct {
	ID             string            `json:"user_id,omitempty"`
	UserName       string            `json:"username,omitempty"`
	Email          string            `json:"email,omitempty"`
	EmailConfirmed bool              `json:"-"`
	PasswordHash   string            `json:"-"`
	UserInfo       map[string]string `json:"user_info,omitempty"`
	UserSessions   []UserSession     `json:"user_sessions,omitempty"`
	CreatedAt      *time.Time        `json:"created_at,omitempty"`
	Roles          []UserRole        `json:"roles,omitempty"`
}

//UserSession struct
type UserSession struct {
	SessionID      string    `json:"session_id,omitempty"`
	ClientName     string    `json:"client_name,omitempty"`
	Device         string    `json:"device,omitempty"`
	LastActiveTime time.Time `json:"last_active_time,omitempty"`
}
type UserRole struct {
	ClientName string   `bson:"client_name,omitempty"`
	Roles      []string `bson:"roles,omitempty"`
}

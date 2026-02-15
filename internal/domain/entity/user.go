package entity

import "time"

// User represents the domain entity for a user.
// This is the core business object.
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}
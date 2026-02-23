package entity

import (
    "fmt"
    "time"
)

type User struct {
    ID           uint      `json:"id"`
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    Password     string    `json:"-"`
    PasswordHash string    `json:"-"`
    Role         string    `json:"role"`
    IsActive     bool      `json:"is_active"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) GetID() string {
    return fmt.Sprintf("%d", u.ID)
}

func (u *User) MigratePassword() {
    if u.Password != "" && u.PasswordHash == "" {
        u.PasswordHash = u.Password
        u.Password = ""
    }
}
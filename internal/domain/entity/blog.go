package entity

import "time"

type Blog struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Slug      string    `json:"slug"`
    Content   string    `json:"content"`
    AuthorID  uint      `json:"author_id"`
    Status    string    `json:"status"` 
    Tags      string    `json:"tags"`   
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
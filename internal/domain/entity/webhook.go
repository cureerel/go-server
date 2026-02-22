package entity

import "time"

type WebhookEvent struct {
    ID        string    `json:"id"`
    Provider  string    `json:"provider"`
    EventType string    `json:"event_type"`
    Payload   []byte    `json:"-"`
    Signature string    `json:"-"`
    Processed bool      `json:"processed"`
    CreatedAt time.Time `json:"created_at"`
}
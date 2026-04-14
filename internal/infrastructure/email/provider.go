// internal/infrastructure/email/provider.go
package email

import "context"

// Email is a single email message.
type Email struct {
	From    string   
	To      []string
	Subject string
	HTML    string 
	Text    string 
}


type Provider interface {
	Send(ctx context.Context, email Email) error
}
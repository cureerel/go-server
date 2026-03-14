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

// Provider is the single interface every email adapter must satisfy.
// To switch from Resend to SES / Postmark / SendGrid — implement this
// interface and swap the wire-up in main.go. Nothing else changes.
type Provider interface {
	Send(ctx context.Context, email Email) error
}
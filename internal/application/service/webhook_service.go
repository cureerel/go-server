package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"
)

type WebhookService struct {
	webhookRepo    repository.WebhookRepository
	stripeSecret   string
	razorpaySecret string
}

type WebhookConfig struct {
	StripeSecret   string
	RazorpaySecret string
}

func NewWebhookService(repo repository.WebhookRepository, cfg WebhookConfig) *WebhookService {
	return &WebhookService{
		webhookRepo:    repo,
		stripeSecret:   cfg.StripeSecret,
		razorpaySecret: cfg.RazorpaySecret,
	}
}

func (s *WebhookService) ProcessStripeWebhook(ctx context.Context, payload []byte, signature string) error {
	if !s.verifyStripeSignature(payload, signature) {
		return fmt.Errorf("invalid stripe signature")
	}

	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to parse stripe webhook: %w", err)
	}

	eventType, _ := event["type"].(string)
	eventID, _ := event["id"].(string)

	if err := s.storeEvent(ctx, "stripe", eventID, eventType, payload, signature); err != nil {
		return err
	}

	switch eventType {
	case "payment_intent.succeeded":
		return s.handleStripePaymentSuccess(ctx, event)
	case "payment_intent.payment_failed":
		return s.handleStripePaymentFailed(ctx, event)
	case "charge.refunded":
		return s.handleStripeRefund(ctx, event)
	default:
		return nil
	}
}

func (s *WebhookService) ProcessRazorpayWebhook(ctx context.Context, payload []byte, signature string) error {
	if !s.verifyRazorpaySignature(payload, signature) {
		return fmt.Errorf("invalid razorpay signature")
	}

	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to parse razorpay webhook: %w", err)
	}

	eventID, _ := event["id"].(string)
	eventType, _ := event["event"].(string)

	if err := s.storeEvent(ctx, "razorpay", eventID, eventType, payload, signature); err != nil {
		return err
	}

	switch eventType {
	case "payment.captured":
		return s.handleRazorpayPaymentCaptured(ctx, event)
	case "payment.failed":
		return s.handleRazorpayPaymentFailed(ctx, event)
	case "refund.processed":
		return s.handleRazorpayRefund(ctx, event)
	default:
		return nil
	}
}

func (s *WebhookService) verifyStripeSignature(payload []byte, signature string) bool {
	expected := hmac.New(sha256.New, []byte(s.stripeSecret))
	expected.Write(payload)
	return hex.EncodeToString(expected.Sum(nil)) == signature
}

func (s *WebhookService) verifyRazorpaySignature(payload []byte, signature string) bool {
	expected := hmac.New(sha256.New, []byte(s.razorpaySecret))
	expected.Write(payload)
	return hex.EncodeToString(expected.Sum(nil)) == signature
}

func (s *WebhookService) storeEvent(ctx context.Context, provider, eventID, eventType string, payload []byte, signature string) error {
	webhookEvent := &entity.WebhookEvent{
		ID:        eventID,
		Provider:  provider,
		EventType: eventType,
		Payload:   payload,
		Signature: signature,
		CreatedAt: time.Now(),
	}
	return s.webhookRepo.StoreEvent(ctx, webhookEvent)
}

func (s *WebhookService) handleStripePaymentSuccess(ctx context.Context, event map[string]interface{}) error {
	return s.createPayment(ctx, "stripe", event)
}

func (s *WebhookService) handleStripePaymentFailed(ctx context.Context, event map[string]interface{}) error {
	return nil
}

func (s *WebhookService) handleStripeRefund(ctx context.Context, event map[string]interface{}) error {
	return nil
}

func (s *WebhookService) handleRazorpayPaymentCaptured(ctx context.Context, event map[string]interface{}) error {
	return s.createPayment(ctx, "razorpay", event)
}

func (s *WebhookService) handleRazorpayPaymentFailed(ctx context.Context, event map[string]interface{}) error {
	return nil
}

func (s *WebhookService) handleRazorpayRefund(ctx context.Context, event map[string]interface{}) error {
	return nil
}

func (s *WebhookService) createPayment(ctx context.Context, provider string, event map[string]interface{}) error {
	payment := &entity.Payment{
		ID:        generateID(),
		Provider:  provider,
		Status:    "completed",
		CreatedAt: time.Now(),
	}
	return s.webhookRepo.CreatePayment(ctx, payment)
}

func generateID() string {
	return fmt.Sprintf("pay_%d", time.Now().UnixNano())
}


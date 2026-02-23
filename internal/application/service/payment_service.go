package service

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/cureerel/gotemplate/internal/domain/entity"
    "github.com/cureerel/gotemplate/internal/domain/repository"
)

type PaymentService struct {
    paymentRepo repository.WebhookRepository // reuses webhook repo since Payment is stored there
    orderRepo   repository.OrderRepository
}

func NewPaymentService(paymentRepo repository.WebhookRepository, orderRepo repository.OrderRepository) *PaymentService {
    return &PaymentService{
        paymentRepo: paymentRepo,
        orderRepo:   orderRepo,
    }
}

type RecordPaymentInput struct {
    UserID        uint
    OrderID       string
    Amount        int64
    Currency      entity.Currency
    Provider      entity.PaymentProvider
    ProviderTxnID string
    CustomerEmail string
    Description   string
}

func (s *PaymentService) Record(ctx context.Context, input RecordPaymentInput) (*entity.Payment, error) {
    payment := &entity.Payment{
        ID:            fmt.Sprintf("pay_%d", time.Now().UnixNano()),
        UserID:        input.UserID,
        OrderID:       input.OrderID,
        Amount:        input.Amount,
        Currency:      input.Currency,
        Status:        entity.PaymentPending,
        Provider:      input.Provider,
        ProviderTxnID: input.ProviderTxnID,
        CustomerEmail: input.CustomerEmail,
        Description:   input.Description,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    if err := s.paymentRepo.SavePayment(ctx, payment); err != nil {
        return nil, err
    }
    return payment, nil
}

func (s *PaymentService) MarkCompleted(ctx context.Context, id string) error {
    payment, err := s.paymentRepo.GetPaymentByID(ctx, id)
    if err != nil {
        return err
    }
    if payment == nil {
        return errors.New("payment not found")
    }
    if payment.Status != entity.PaymentPending {
        return errors.New("only pending payments can be marked completed")
    }
    return s.paymentRepo.UpdatePaymentStatus(ctx, id, string(entity.PaymentCompleted))
}

func (s *PaymentService) MarkFailed(ctx context.Context, id string) error {
    payment, err := s.paymentRepo.GetPaymentByID(ctx, id)
    if err != nil {
        return err
    }
    if payment == nil {
        return errors.New("payment not found")
    }
    return s.paymentRepo.UpdatePaymentStatus(ctx, id, string(entity.PaymentFailed))
}

func (s *PaymentService) Refund(ctx context.Context, id string) error {
    payment, err := s.paymentRepo.GetPaymentByID(ctx, id)
    if err != nil {
        return err
    }
    if payment == nil {
        return errors.New("payment not found")
    }
    if payment.Status != entity.PaymentCompleted {
        return errors.New("only completed payments can be refunded")
    }
    return s.paymentRepo.UpdatePaymentStatus(ctx, id, string(entity.PaymentRefunded))
}

func (s *PaymentService) GetByID(ctx context.Context, id string) (*entity.Payment, error) {
    payment, err := s.paymentRepo.GetPaymentByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if payment == nil {
        return nil, errors.New("payment not found")
    }
    return payment, nil
}
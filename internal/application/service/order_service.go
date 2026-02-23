package service

import (
    "context"
    "errors"

    "github.com/cureerel/gotemplate/internal/domain/entity"
    "github.com/cureerel/gotemplate/internal/domain/repository"
)

type OrderService struct {
    orderRepo   repository.OrderRepository
    productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) *OrderService {
    return &OrderService{
        orderRepo:   orderRepo,
        productRepo: productRepo,
    }
}

type OrderItemInput struct {
    ProductID uint
    Quantity  int
}

type CreateOrderInput struct {
    UserID   uint
    Currency entity.Currency
    Items    []OrderItemInput
}

func (s *OrderService) Create(ctx context.Context, input CreateOrderInput) (*entity.Order, error) {
    if len(input.Items) == 0 {
        return nil, errors.New("order must contain at least one item")
    }

    // Resolve products and build order items
    var orderItems []entity.OrderItem
    for _, item := range input.Items {
        if item.Quantity <= 0 {
            return nil, errors.New("item quantity must be greater than zero")
        }

        product, err := s.productRepo.GetByID(ctx, item.ProductID)
        if err != nil {
            return nil, err
        }
        if product == nil {
            return nil, errors.New("product not found")
        }
        if !product.IsActive {
            return nil, errors.New("product is not available")
        }

        orderItems = append(orderItems, entity.OrderItem{
            ProductID: product.ID,
            Type:      product.Type,
            Quantity:  item.Quantity,
            UnitPrice: product.Price,
        })
    }

    order := &entity.Order{
        UserID:   input.UserID,
        Status:   entity.OrderPending,
        Currency: input.Currency,
        Items:    orderItems,
    }
    order.TotalAmount = order.CalculateTotal()

    if err := s.orderRepo.Create(ctx, order); err != nil {
        return nil, err
    }
    return order, nil
}

func (s *OrderService) GetByID(ctx context.Context, id uint) (*entity.Order, error) {
    order, err := s.orderRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if order == nil {
        return nil, errors.New("order not found")
    }
    return order, nil
}

func (s *OrderService) GetByUser(ctx context.Context, userID uint, page, limit int) ([]entity.Order, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }
    return s.orderRepo.GetByUser(ctx, userID, page, limit)
}

func (s *OrderService) UpdateStatus(ctx context.Context, id uint, status entity.OrderStatus) error {
    order, err := s.orderRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if order == nil {
        return errors.New("order not found")
    }

    if !isValidOrderTransition(order.Status, status) {
        return errors.New("invalid order status transition")
    }

    return s.orderRepo.UpdateStatus(ctx, id, string(status))
}

// Enforce valid state transitions
func isValidOrderTransition(current, next entity.OrderStatus) bool {
    allowed := map[entity.OrderStatus][]entity.OrderStatus{
        entity.OrderPending:    {entity.OrderConfirmed, entity.OrderCancelled},
        entity.OrderConfirmed:  {entity.OrderDispatched, entity.OrderCancelled},
        entity.OrderDispatched: {entity.OrderDelivered},
        entity.OrderDelivered:  {entity.OrderCompleted},
    }
    for _, s := range allowed[current] {
        if s == next {
            return true
        }
    }
    return false
}
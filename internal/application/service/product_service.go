package service

import (
    "context"
    "errors"

    "github.com/cureerel/gotemplate/internal/domain/entity"
    "github.com/cureerel/gotemplate/internal/domain/repository"
)

type ProductService struct {
    productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
    return &ProductService{productRepo: productRepo}
}

type CreateProductInput struct {
    Name        string
    Description string
    Type        entity.ProductType
    Price       int64
    Currency    entity.Currency
}

type UpdateProductInput struct {
    ID          uint
    Name        *string
    Description *string
    Price       *int64
    Currency    *entity.Currency
    IsActive    *bool
}

func (s *ProductService) Create(ctx context.Context, input CreateProductInput) (*entity.Product, error) {
    if input.Type != entity.ProductPhysical && input.Type != entity.ProductDigital {
        return nil, errors.New("invalid product type: must be physical or digital")
    }
    if input.Price <= 0 {
        return nil, errors.New("price must be greater than zero")
    }

    product := &entity.Product{
        Name:        input.Name,
        Description: input.Description,
        Type:        input.Type,
        Price:       input.Price,
        Currency:    input.Currency,
        IsActive:    true,
    }

    if err := s.productRepo.Create(ctx, product); err != nil {
        return nil, err
    }
    return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*entity.Product, error) {
    product, err := s.productRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if product == nil {
        return nil, errors.New("product not found")
    }
    return product, nil
}

func (s *ProductService) GetAll(ctx context.Context, page, limit int) ([]entity.Product, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }
    return s.productRepo.GetAll(ctx, page, limit)
}

func (s *ProductService) Update(ctx context.Context, input UpdateProductInput) (*entity.Product, error) {
    product, err := s.productRepo.GetByID(ctx, input.ID)
    if err != nil {
        return nil, err
    }
    if product == nil {
        return nil, errors.New("product not found")
    }

    if input.Name != nil {
        product.Name = *input.Name
    }
    if input.Description != nil {
        product.Description = *input.Description
    }
    if input.Price != nil {
        if *input.Price <= 0 {
            return nil, errors.New("price must be greater than zero")
        }
        product.Price = *input.Price
    }
    if input.Currency != nil {
        product.Currency = *input.Currency
    }
    if input.IsActive != nil {
        product.IsActive = *input.IsActive
    }

    if err := s.productRepo.Update(ctx, product); err != nil {
        return nil, err
    }
    return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
    return s.productRepo.Delete(ctx, id)
}
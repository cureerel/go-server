// internal/application/service/product_service.go
package service

import (
	"context"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/pkg/apperror"
	"github.com/cureerel/cserver/pkg/idgen"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

type CreateProductInput struct {
	Name        string
	Description string
	Type        string // physical | digital
	Price       int64
	Currency    string
	Stock       int // -1 for unlimited (auto-set for digital)
	ImageURL    string
}

type UpdateProductInput struct {
	Name        *string
	Description *string
	Price       *int64
	Stock       *int
	ImageURL    *string
	IsActive    *bool
}

func (s *ProductService) Create(ctx context.Context, in CreateProductInput) (*entity.Product, error) {
	pType := entity.ProductType(in.Type)
	if pType != entity.ProductPhysical && pType != entity.ProductDigital {
		return nil, apperror.NewBadRequest("type must be physical or digital")
	}

	stock := in.Stock
	if pType == entity.ProductDigital {
		stock = -1 // digital products are always unlimited
	}

	cur := entity.Currency(in.Currency)
	if cur == "" {
		cur = entity.CurrencyUSD
	}

	p := &entity.Product{
		SKU:         idgen.New(idgen.PrefixProduct),
		Name:        in.Name,
		Description: in.Description,
		Type:        pType,
		Price:       in.Price,
		Currency:    cur,
		Stock:       stock,
		ImageURL:    in.ImageURL,
		IsActive:    true,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, apperror.NewInternal(err, "failed to create product")
	}
	return p, nil
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*entity.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch product")
	}
	if p == nil {
		return nil, apperror.NewNotFound("product not found")
	}
	return p, nil
}

func (s *ProductService) GetAll(ctx context.Context, page, limit int, productType string) ([]entity.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.GetAll(ctx, repository.ProductFilter{
		Page:  page,
		Limit: limit,
		Type:  productType,
	})
}

func (s *ProductService) Update(ctx context.Context, id uint, in UpdateProductInput) (*entity.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch product")
	}
	if p == nil {
		return nil, apperror.NewNotFound("product not found")
	}

	if in.Name != nil {
		p.Name = *in.Name
	}
	if in.Description != nil {
		p.Description = *in.Description
	}
	if in.Price != nil {
		p.Price = *in.Price
	}
	if in.Stock != nil && p.Type != entity.ProductDigital {
		p.Stock = *in.Stock
	}
	if in.ImageURL != nil {
		p.ImageURL = *in.ImageURL
	}
	if in.IsActive != nil {
		p.IsActive = *in.IsActive
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, apperror.NewInternal(err, "failed to update product")
	}
	return p, nil
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperror.NewInternal(err, "failed to fetch product")
	}
	if p == nil {
		return apperror.NewNotFound("product not found")
	}
	return s.repo.Delete(ctx, id)
}

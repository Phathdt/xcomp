package services

import (
	"context"
	"math"
	"time"

	"example/modules/product/application/dto"
	"example/modules/product/domain/entities"
	"example/modules/product/domain/repositories"

	"github.com/google/uuid"
)

type ProductService struct {
	ProductRepo      repositories.ProductRepository      `inject:"ProductRepository"`
	ProductCacheRepo repositories.ProductCacheRepository `inject:"ProductCacheRepository"`
}

func (ps *ProductService) GetServiceName() string {
	return "ProductService"
}

func (ps *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := ps.ProductCacheRepo.Get(ctx, id)
	if err != nil {
		product, err = ps.ProductRepo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		if setErr := ps.ProductCacheRepo.Set(ctx, product, 5*time.Minute); setErr != nil {
		}
	} else if product == nil {
		product, err = ps.ProductRepo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		if setErr := ps.ProductCacheRepo.Set(ctx, product, 5*time.Minute); setErr != nil {
		}
	}

	return ps.toProductResponse(product), nil
}

func (ps *ProductService) ListProducts(ctx context.Context, page, pageSize int32) (*dto.ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	products, err := ps.ProductRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := ps.ProductRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := int32(math.Ceil(float64(totalCount) / float64(pageSize)))

	response := &dto.ProductListResponse{
		Products:   make([]*dto.ProductResponse, len(products)),
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	for i, product := range products {
		response.Products[i] = ps.toProductResponse(product)
	}

	return response, nil
}

func (ps *ProductService) ListProductsByCategory(ctx context.Context, category string, page, pageSize int32) (*dto.ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	products, err := ps.ProductRepo.ListByCategory(ctx, category, pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := ps.ProductRepo.CountByCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	totalPages := int32(math.Ceil(float64(totalCount) / float64(pageSize)))

	response := &dto.ProductListResponse{
		Products:   make([]*dto.ProductResponse, len(products)),
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	for i, product := range products {
		response.Products[i] = ps.toProductResponse(product)
	}

	return response, nil
}

func (ps *ProductService) SearchProducts(ctx context.Context, searchReq *dto.ProductSearchRequest) (*dto.ProductListResponse, error) {
	if searchReq.Page < 1 {
		searchReq.Page = 1
	}
	if searchReq.PageSize < 1 || searchReq.PageSize > 100 {
		searchReq.PageSize = 10
	}

	offset := (searchReq.Page - 1) * searchReq.PageSize

	products, err := ps.ProductRepo.Search(ctx, searchReq.Query, searchReq.PageSize, offset)
	if err != nil {
		return nil, err
	}

	totalCount := int64(len(products))
	totalPages := int32(math.Ceil(float64(totalCount) / float64(searchReq.PageSize)))

	response := &dto.ProductListResponse{
		Products:   make([]*dto.ProductResponse, len(products)),
		TotalCount: totalCount,
		Page:       searchReq.Page,
		PageSize:   searchReq.PageSize,
		TotalPages: totalPages,
	}

	for i, product := range products {
		response.Products[i] = ps.toProductResponse(product)
	}

	return response, nil
}

func (ps *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := &entities.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Category:      req.Category,
		IsActive:      true,
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}

	createdProduct, err := ps.ProductRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return ps.toProductResponse(createdProduct), nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, id uuid.UUID, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	existingProduct, err := ps.ProductRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existingProduct.Name = req.Name
	existingProduct.Description = req.Description
	existingProduct.Price = req.Price
	existingProduct.StockQuantity = req.StockQuantity
	existingProduct.Category = req.Category

	if err := existingProduct.Validate(); err != nil {
		return nil, err
	}

	updatedProduct, err := ps.ProductRepo.Update(ctx, existingProduct)
	if err != nil {
		return nil, err
	}

	ps.ProductCacheRepo.Delete(ctx, id)

	return ps.toProductResponse(updatedProduct), nil
}

func (ps *ProductService) UpdateProductStock(ctx context.Context, id uuid.UUID, req *dto.UpdateStockRequest) (*dto.ProductResponse, error) {
	updatedProduct, err := ps.ProductRepo.UpdateStock(ctx, id, req.StockQuantity)
	if err != nil {
		return nil, err
	}

	ps.ProductCacheRepo.Delete(ctx, id)

	return ps.toProductResponse(updatedProduct), nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	err := ps.ProductRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	ps.ProductCacheRepo.Delete(ctx, id)
	return nil
}

func (ps *ProductService) toProductResponse(product *entities.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		Category:      product.Category,
		IsActive:      product.IsActive,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}
}

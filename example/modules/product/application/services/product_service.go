package services

import (
	"context"
	"math"
	"time"

	"example/modules/product/application/dto"
	"example/modules/product/domain/entities"
	"example/modules/product/domain/interfaces"

	"xcomp"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepo      interfaces.ProductRepository      // lowercase - manual injection
	productCacheRepo interfaces.ProductCacheRepository // lowercase - manual injection
	Logger           xcomp.Logger                      `inject:"Logger"` // uppercase - auto injection
}

func NewProductService() *ProductService {
	return &ProductService{}
}

// Method injection for lowercase fields
func (ps *ProductService) SetDependencies(
	productRepo interfaces.ProductRepository,
	productCacheRepo interfaces.ProductCacheRepository,
) {
	ps.productRepo = productRepo
	ps.productCacheRepo = productCacheRepo
}

func (ps *ProductService) GetServiceName() string {
	return "ProductService"
}

func (ps *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	ps.Logger.Debug("Getting product", xcomp.Field("product_id", id))

	product, err := ps.productCacheRepo.Get(ctx, id)
	if err != nil {
		ps.Logger.Debug("Product not found in cache, fetching from database",
			xcomp.Field("product_id", id),
			xcomp.Field("cache_error", err))

		product, err = ps.productRepo.GetByID(ctx, id)
		if err != nil {
			ps.Logger.Error("Failed to get product from database",
				xcomp.Field("product_id", id),
				xcomp.Field("error", err))
			return nil, err
		}

		if setErr := ps.productCacheRepo.Set(ctx, product, 5*time.Minute); setErr != nil {
			ps.Logger.Warn("Failed to cache product",
				xcomp.Field("product_id", id),
				xcomp.Field("error", setErr))
		}
	} else if product == nil {
		ps.Logger.Debug("Product cache miss, fetching from database",
			xcomp.Field("product_id", id))

		product, err = ps.productRepo.GetByID(ctx, id)
		if err != nil {
			ps.Logger.Error("Failed to get product from database",
				xcomp.Field("product_id", id),
				xcomp.Field("error", err))
			return nil, err
		}

		if setErr := ps.productCacheRepo.Set(ctx, product, 5*time.Minute); setErr != nil {
			ps.Logger.Warn("Failed to cache product",
				xcomp.Field("product_id", id),
				xcomp.Field("error", setErr))
		}
	} else {
		ps.Logger.Debug("Product found in cache", xcomp.Field("product_id", id))
	}

	ps.Logger.Info("Product retrieved successfully",
		xcomp.Field("product_id", id),
		xcomp.Field("product_name", product.Name))

	return ps.toProductResponse(product), nil
}

func (ps *ProductService) ListProducts(ctx context.Context, page, pageSize int32) (*dto.ProductListResponse, error) {
	ps.Logger.Debug("Getting product", xcomp.Field("page", page), xcomp.Field("page_size", pageSize))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	products, err := ps.productRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := ps.productRepo.Count(ctx)
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

	products, err := ps.productRepo.ListByCategory(ctx, category, pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := ps.productRepo.CountByCategory(ctx, category)
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

	products, err := ps.productRepo.Search(ctx, searchReq.Query, searchReq.PageSize, offset)
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
	ps.Logger.Info("Creating new product",
		xcomp.Field("product_name", req.Name),
		xcomp.Field("price", req.Price),
		xcomp.Field("stock_quantity", req.StockQuantity))

	product := &entities.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Category:      req.Category,
		IsActive:      true,
	}

	if err := product.Validate(); err != nil {
		ps.Logger.Error("Product validation failed",
			xcomp.Field("product_name", req.Name),
			xcomp.Field("error", err))
		return nil, err
	}

	createdProduct, err := ps.productRepo.Create(ctx, product)
	if err != nil {
		ps.Logger.Error("Failed to create product",
			xcomp.Field("product_name", req.Name),
			xcomp.Field("error", err))
		return nil, err
	}

	ps.Logger.Info("Product created successfully",
		xcomp.Field("product_id", createdProduct.ID),
		xcomp.Field("product_name", createdProduct.Name))

	return ps.toProductResponse(createdProduct), nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, id uuid.UUID, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	existingProduct, err := ps.productRepo.GetByID(ctx, id)
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

	updatedProduct, err := ps.productRepo.Update(ctx, existingProduct)
	if err != nil {
		return nil, err
	}

	ps.productCacheRepo.Delete(ctx, id)

	return ps.toProductResponse(updatedProduct), nil
}

func (ps *ProductService) UpdateProductStock(ctx context.Context, id uuid.UUID, req *dto.UpdateStockRequest) (*dto.ProductResponse, error) {
	updatedProduct, err := ps.productRepo.UpdateStock(ctx, id, req.StockQuantity)
	if err != nil {
		return nil, err
	}

	ps.productCacheRepo.Delete(ctx, id)

	return ps.toProductResponse(updatedProduct), nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	err := ps.productRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	ps.productCacheRepo.Delete(ctx, id)
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

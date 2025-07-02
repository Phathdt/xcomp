package persistence

import (
	"context"
	"fmt"
	"time"

	"example/infrastructure/database"
	"example/modules/product/domain/entities"
	"example/modules/product/domain/interfaces"
	"example/modules/product/infrastructure/query/gen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProductRepositoryImpl struct {
	DbConnection *database.DatabaseConnection `inject:"DatabaseConnection"`
	queries      *gen.Queries
}

func (pr *ProductRepositoryImpl) GetServiceName() string {
	return "ProductRepository"
}

func (pr *ProductRepositoryImpl) Initialize() {
	if pr.DbConnection != nil && pr.DbConnection.GetDB() != nil {
		pr.queries = gen.New(pr.DbConnection.GetDB())
	}
}

func (pr *ProductRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	pgID := pgtype.UUID{}
	if err := pgID.Scan(id.String()); err != nil {
		return nil, fmt.Errorf("failed to convert UUID: %w", err)
	}

	result, err := pr.queries.GetProduct(ctx, pgID)
	if err != nil {
		return nil, pr.convertError(err)
	}

	return pr.convertToEntity(result), nil
}

func (pr *ProductRepositoryImpl) List(ctx context.Context, limit, offset int32) ([]*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	results, err := pr.queries.ListProducts(ctx, gen.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	products := make([]*entities.Product, len(results))
	for i, result := range results {
		products[i] = pr.convertToEntity(result)
	}

	return products, nil
}

func (pr *ProductRepositoryImpl) ListByCategory(ctx context.Context, category string, limit, offset int32) ([]*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	results, err := pr.queries.ListProductsByCategory(ctx, gen.ListProductsByCategoryParams{
		Category: &category,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list products by category: %w", err)
	}

	products := make([]*entities.Product, len(results))
	for i, result := range results {
		products[i] = pr.convertToEntity(result)
	}

	return products, nil
}

func (pr *ProductRepositoryImpl) Search(ctx context.Context, searchQuery string, limit, offset int32) ([]*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	results, err := pr.queries.SearchProducts(ctx, gen.SearchProductsParams{
		Column1: &searchQuery,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	products := make([]*entities.Product, len(results))
	for i, result := range results {
		products[i] = pr.convertToEntity(result)
	}

	return products, nil
}

func (pr *ProductRepositoryImpl) Create(ctx context.Context, product *entities.Product) (*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	pgPrice := pgtype.Numeric{}
	if err := pgPrice.Scan(fmt.Sprintf("%.2f", product.Price)); err != nil {
		return nil, fmt.Errorf("failed to convert price: %w", err)
	}

	result, err := pr.queries.CreateProduct(ctx, gen.CreateProductParams{
		Name:          product.Name,
		Description:   product.Description,
		Price:         pgPrice,
		StockQuantity: product.StockQuantity,
		Category:      product.Category,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return pr.convertToEntity(result), nil
}

func (pr *ProductRepositoryImpl) Update(ctx context.Context, product *entities.Product) (*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	pgID := pgtype.UUID{}
	if err := pgID.Scan(product.ID.String()); err != nil {
		return nil, fmt.Errorf("failed to convert UUID: %w", err)
	}

	pgPrice := pgtype.Numeric{}
	if err := pgPrice.Scan(fmt.Sprintf("%.2f", product.Price)); err != nil {
		return nil, fmt.Errorf("failed to convert price: %w", err)
	}

	result, err := pr.queries.UpdateProduct(ctx, gen.UpdateProductParams{
		ID:            pgID,
		Name:          product.Name,
		Description:   product.Description,
		Price:         pgPrice,
		StockQuantity: product.StockQuantity,
		Category:      product.Category,
	})
	if err != nil {
		return nil, pr.convertError(err)
	}

	return pr.convertToEntity(result), nil
}

func (pr *ProductRepositoryImpl) UpdateStock(ctx context.Context, id uuid.UUID, stockQuantity int32) (*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	pgID := pgtype.UUID{}
	if err := pgID.Scan(id.String()); err != nil {
		return nil, fmt.Errorf("failed to convert UUID: %w", err)
	}

	result, err := pr.queries.UpdateProductStock(ctx, gen.UpdateProductStockParams{
		ID:            pgID,
		StockQuantity: stockQuantity,
	})
	if err != nil {
		return nil, pr.convertError(err)
	}

	return pr.convertToEntity(result), nil
}

func (pr *ProductRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if pr.queries == nil {
		pr.Initialize()
	}

	pgID := pgtype.UUID{}
	if err := pgID.Scan(id.String()); err != nil {
		return fmt.Errorf("failed to convert UUID: %w", err)
	}

	return pr.queries.DeleteProduct(ctx, pgID)
}

func (pr *ProductRepositoryImpl) Count(ctx context.Context) (int64, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	return pr.queries.CountProducts(ctx)
}

func (pr *ProductRepositoryImpl) CountByCategory(ctx context.Context, category string) (int64, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	return pr.queries.CountProductsByCategory(ctx, &category)
}

func (pr *ProductRepositoryImpl) convertToEntity(sqlcProduct *gen.Product) *entities.Product {
	var id uuid.UUID
	if sqlcProduct.ID.Valid {
		id = uuid.UUID(sqlcProduct.ID.Bytes)
	}

	var price float64
	if sqlcProduct.Price.Valid {
		if f, err := sqlcProduct.Price.Float64Value(); err == nil {
			price = f.Float64
		}
	}

	var createdAt, updatedAt time.Time
	if sqlcProduct.CreatedAt.Valid {
		createdAt = sqlcProduct.CreatedAt.Time
	}
	if sqlcProduct.UpdatedAt.Valid {
		updatedAt = sqlcProduct.UpdatedAt.Time
	}

	return &entities.Product{
		ID:            id,
		Name:          sqlcProduct.Name,
		Description:   sqlcProduct.Description,
		Price:         price,
		StockQuantity: sqlcProduct.StockQuantity,
		Category:      sqlcProduct.Category,
		IsActive:      sqlcProduct.IsActive,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func (pr *ProductRepositoryImpl) convertError(err error) error {
	if err.Error() == "no rows in result set" {
		return entities.ErrProductNotFound
	}
	return fmt.Errorf("database error: %w", err)
}

var _ interfaces.ProductRepository = (*ProductRepositoryImpl)(nil)

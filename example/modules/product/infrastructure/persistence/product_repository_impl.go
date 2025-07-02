package persistence

import (
	"context"
	"fmt"
	"time"

	"example/infrastructure/database"
	"example/modules/product/domain/entities"
	"example/modules/product/domain/interfaces"
	"example/modules/product/infrastructure/persistence/queries"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProductRepositoryImpl struct {
	DbConnection *database.DatabaseConnection `inject:"DatabaseConnection"`
	queries      *queries.Queries
}

func (pr *ProductRepositoryImpl) GetServiceName() string {
	return "ProductRepository"
}

func (pr *ProductRepositoryImpl) Initialize() {
	if pr.DbConnection != nil && pr.DbConnection.GetDB() != nil {
		pr.queries = queries.New(pr.DbConnection.GetDB())
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

	return pr.convertToEntity(&result), nil
}

func (pr *ProductRepositoryImpl) List(ctx context.Context, limit, offset int32) ([]*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	results, err := pr.queries.ListProducts(ctx, queries.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	products := make([]*entities.Product, len(results))
	for i, result := range results {
		products[i] = pr.convertToEntity(&result)
	}

	return products, nil
}

func (pr *ProductRepositoryImpl) ListByCategory(ctx context.Context, category string, limit, offset int32) ([]*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	pgCategory := pgtype.Text{}
	if err := pgCategory.Scan(category); err != nil {
		return nil, fmt.Errorf("failed to convert category: %w", err)
	}

	results, err := pr.queries.ListProductsByCategory(ctx, queries.ListProductsByCategoryParams{
		Category: pgCategory,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list products by category: %w", err)
	}

	products := make([]*entities.Product, len(results))
	for i, result := range results {
		products[i] = pr.convertToEntity(&result)
	}

	return products, nil
}

func (pr *ProductRepositoryImpl) Search(ctx context.Context, searchQuery string, limit, offset int32) ([]*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	pgQuery := pgtype.Text{}
	if err := pgQuery.Scan(searchQuery); err != nil {
		return nil, fmt.Errorf("failed to convert search query: %w", err)
	}

	results, err := pr.queries.SearchProducts(ctx, queries.SearchProductsParams{
		Column1: pgQuery,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	products := make([]*entities.Product, len(results))
	for i, result := range results {
		products[i] = pr.convertToEntity(&result)
	}

	return products, nil
}

func (pr *ProductRepositoryImpl) Create(ctx context.Context, product *entities.Product) (*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	pgDescription := pgtype.Text{}
	if product.Description != nil {
		if err := pgDescription.Scan(*product.Description); err != nil {
			return nil, fmt.Errorf("failed to convert description: %w", err)
		}
	}

	pgPrice := pgtype.Numeric{}
	if err := pgPrice.Scan(fmt.Sprintf("%.2f", product.Price)); err != nil {
		return nil, fmt.Errorf("failed to convert price: %w", err)
	}

	pgCategory := pgtype.Text{}
	if product.Category != nil {
		if err := pgCategory.Scan(*product.Category); err != nil {
			return nil, fmt.Errorf("failed to convert category: %w", err)
		}
	}

	result, err := pr.queries.CreateProduct(ctx, queries.CreateProductParams{
		Name:          product.Name,
		Description:   pgDescription,
		Price:         pgPrice,
		StockQuantity: product.StockQuantity,
		Category:      pgCategory,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return pr.convertToEntity(&result), nil
}

func (pr *ProductRepositoryImpl) Update(ctx context.Context, product *entities.Product) (*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	pgID := pgtype.UUID{}
	if err := pgID.Scan(product.ID.String()); err != nil {
		return nil, fmt.Errorf("failed to convert UUID: %w", err)
	}

	pgDescription := pgtype.Text{}
	if product.Description != nil {
		if err := pgDescription.Scan(*product.Description); err != nil {
			return nil, fmt.Errorf("failed to convert description: %w", err)
		}
	}

	pgPrice := pgtype.Numeric{}
	if err := pgPrice.Scan(fmt.Sprintf("%.2f", product.Price)); err != nil {
		return nil, fmt.Errorf("failed to convert price: %w", err)
	}

	pgCategory := pgtype.Text{}
	if product.Category != nil {
		if err := pgCategory.Scan(*product.Category); err != nil {
			return nil, fmt.Errorf("failed to convert category: %w", err)
		}
	}

	result, err := pr.queries.UpdateProduct(ctx, queries.UpdateProductParams{
		ID:            pgID,
		Name:          product.Name,
		Description:   pgDescription,
		Price:         pgPrice,
		StockQuantity: product.StockQuantity,
		Category:      pgCategory,
	})
	if err != nil {
		return nil, pr.convertError(err)
	}

	return pr.convertToEntity(&result), nil
}

func (pr *ProductRepositoryImpl) UpdateStock(ctx context.Context, id uuid.UUID, stockQuantity int32) (*entities.Product, error) {
	if pr.queries == nil {
		pr.Initialize()
	}

	pgID := pgtype.UUID{}
	if err := pgID.Scan(id.String()); err != nil {
		return nil, fmt.Errorf("failed to convert UUID: %w", err)
	}

	result, err := pr.queries.UpdateProductStock(ctx, queries.UpdateProductStockParams{
		ID:            pgID,
		StockQuantity: stockQuantity,
	})
	if err != nil {
		return nil, pr.convertError(err)
	}

	return pr.convertToEntity(&result), nil
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

	pgCategory := pgtype.Text{}
	if err := pgCategory.Scan(category); err != nil {
		return 0, fmt.Errorf("failed to convert category: %w", err)
	}

	return pr.queries.CountProductsByCategory(ctx, pgCategory)
}

func (pr *ProductRepositoryImpl) convertToEntity(sqlcProduct *queries.Product) *entities.Product {
	var id uuid.UUID
	if sqlcProduct.ID.Valid {
		id = uuid.UUID(sqlcProduct.ID.Bytes)
	}

	var description *string
	if sqlcProduct.Description.Valid {
		description = &sqlcProduct.Description.String
	}

	var price float64
	if sqlcProduct.Price.Valid {
		if f, err := sqlcProduct.Price.Float64Value(); err == nil {
			price = f.Float64
		}
	}

	var category *string
	if sqlcProduct.Category.Valid {
		category = &sqlcProduct.Category.String
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
		Description:   description,
		Price:         price,
		StockQuantity: sqlcProduct.StockQuantity,
		Category:      category,
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

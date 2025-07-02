package repositories

import (
	"context"
	"fmt"
	"time"

	"example/modules/customer/domain/entities"
	"example/modules/customer/infrastructure/query/gen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepositoryImpl struct {
	DB      *pgxpool.Pool `inject:"DatabaseConnection"`
	Queries *gen.Queries
}

func (r *CustomerRepositoryImpl) GetServiceName() string {
	return "CustomerRepositoryImpl"
}

func (r *CustomerRepositoryImpl) Initialize() {
	if r.DB != nil {
		r.Queries = gen.New(r.DB)
	}
}

func (r *CustomerRepositoryImpl) Create(ctx context.Context, customer *entities.Customer) (*entities.Customer, error) {
	if r.Queries == nil {
		r.Initialize()
	}

	result, err := r.Queries.CreateCustomer(ctx, gen.CreateCustomerParams{
		Username: customer.Username,
		Email:    customer.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return r.convertToEntity(result), nil
}

func (r *CustomerRepositoryImpl) Update(ctx context.Context, customer *entities.Customer) (*entities.Customer, error) {
	if r.Queries == nil {
		r.Initialize()
	}

	pgID := pgtype.UUID{}
	if err := pgID.Scan(customer.ID.String()); err != nil {
		return nil, fmt.Errorf("failed to convert UUID: %w", err)
	}

	result, err := r.Queries.UpdateCustomer(ctx, gen.UpdateCustomerParams{
		ID:       pgID,
		Username: customer.Username,
		Email:    customer.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return r.convertToEntity(result), nil
}

func (r *CustomerRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if r.Queries == nil {
		r.Initialize()
	}

	pgID := pgtype.UUID{}
	if err := pgID.Scan(id.String()); err != nil {
		return fmt.Errorf("failed to convert UUID: %w", err)
	}

	return r.Queries.DeleteCustomer(ctx, pgID)
}

func (r *CustomerRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.Customer, error) {
	if r.Queries == nil {
		r.Initialize()
	}

	pgID := pgtype.UUID{}
	if err := pgID.Scan(id.String()); err != nil {
		return nil, fmt.Errorf("failed to convert UUID: %w", err)
	}

	result, err := r.Queries.GetCustomer(ctx, pgID)
	if err != nil {
		return nil, r.convertError(err)
	}

	return r.convertToEntity(result), nil
}

func (r *CustomerRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entities.Customer, error) {
	if r.Queries == nil {
		r.Initialize()
	}

	result, err := r.Queries.GetCustomerByUsername(ctx, username)
	if err != nil {
		return nil, r.convertError(err)
	}

	return r.convertToEntity(result), nil
}

func (r *CustomerRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.Customer, error) {
	if r.Queries == nil {
		r.Initialize()
	}

	result, err := r.Queries.GetCustomerByEmail(ctx, email)
	if err != nil {
		return nil, r.convertError(err)
	}

	return r.convertToEntity(result), nil
}

func (r *CustomerRepositoryImpl) List(ctx context.Context, limit, offset int32) ([]*entities.Customer, error) {
	if r.Queries == nil {
		r.Initialize()
	}

	results, err := r.Queries.ListCustomers(ctx, gen.ListCustomersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	customers := make([]*entities.Customer, len(results))
	for i, result := range results {
		customers[i] = r.convertToEntity(result)
	}

	return customers, nil
}

func (r *CustomerRepositoryImpl) Search(ctx context.Context, query string, limit, offset int32) ([]*entities.Customer, error) {
	if r.Queries == nil {
		r.Initialize()
	}

	results, err := r.Queries.SearchCustomers(ctx, gen.SearchCustomersParams{
		Column1: &query,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}

	customers := make([]*entities.Customer, len(results))
	for i, result := range results {
		customers[i] = r.convertToEntity(result)
	}

	return customers, nil
}

func (r *CustomerRepositoryImpl) Count(ctx context.Context) (int64, error) {
	if r.Queries == nil {
		r.Initialize()
	}

	return r.Queries.CountCustomers(ctx)
}

func (r *CustomerRepositoryImpl) convertToEntity(sqlcCustomer *gen.Customer) *entities.Customer {
	var id uuid.UUID
	if sqlcCustomer.ID.Valid {
		id = uuid.UUID(sqlcCustomer.ID.Bytes)
	}

	var createdAt, updatedAt time.Time
	if sqlcCustomer.CreatedAt.Valid {
		createdAt = sqlcCustomer.CreatedAt.Time
	}
	if sqlcCustomer.UpdatedAt.Valid {
		updatedAt = sqlcCustomer.UpdatedAt.Time
	}

	return &entities.Customer{
		ID:        id,
		Username:  sqlcCustomer.Username,
		Email:     sqlcCustomer.Email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (r *CustomerRepositoryImpl) convertError(err error) error {
	if err.Error() == "no rows in result set" {
		return entities.ErrCustomerNotFound
	}
	return fmt.Errorf("database error: %w", err)
}

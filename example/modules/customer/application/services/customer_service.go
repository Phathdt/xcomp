package services

import (
	"context"
	"time"

	"example/modules/customer/application/dto"
	"example/modules/customer/domain/entities"
	"example/modules/customer/domain/interfaces"

	"github.com/google/uuid"
)

type CustomerService struct {
	CustomerRepository      interfaces.CustomerRepository      `inject:"CustomerRepository"`
	CustomerCacheRepository interfaces.CustomerCacheRepository `inject:"CustomerCacheRepository"`
}

func (cs *CustomerService) GetServiceName() string {
	return "CustomerService"
}

func (cs *CustomerService) CreateCustomer(ctx context.Context, req *dto.CreateCustomerRequest) (*dto.CustomerResponse, error) {
	customer := &entities.Customer{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := customer.Validate(); err != nil {
		return nil, err
	}

	existingCustomer, _ := cs.CustomerRepository.GetByUsername(ctx, req.Username)
	if existingCustomer != nil {
		return nil, entities.ErrCustomerUsernameExists
	}

	existingCustomer, _ = cs.CustomerRepository.GetByEmail(ctx, req.Email)
	if existingCustomer != nil {
		return nil, entities.ErrCustomerEmailExists
	}

	createdCustomer, err := cs.CustomerRepository.Create(ctx, customer)
	if err != nil {
		return nil, err
	}

	cs.CustomerCacheRepository.Set(ctx, cs.CustomerCacheRepository.GetCustomerCacheKey(createdCustomer.ID), createdCustomer, 30*time.Minute)
	cs.CustomerCacheRepository.Set(ctx, cs.CustomerCacheRepository.GetCustomerUsernameCacheKey(createdCustomer.Username), createdCustomer, 30*time.Minute)
	cs.CustomerCacheRepository.Set(ctx, cs.CustomerCacheRepository.GetCustomerEmailCacheKey(createdCustomer.Email), createdCustomer, 30*time.Minute)

	return cs.mapToCustomerResponse(createdCustomer), nil
}

func (cs *CustomerService) UpdateCustomer(ctx context.Context, id uuid.UUID, req *dto.UpdateCustomerRequest) (*dto.CustomerResponse, error) {
	existingCustomer, err := cs.CustomerRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingCustomer == nil {
		return nil, entities.ErrCustomerNotFound
	}

	if req.Username != existingCustomer.Username {
		usernameTaken, _ := cs.CustomerRepository.GetByUsername(ctx, req.Username)
		if usernameTaken != nil {
			return nil, entities.ErrCustomerUsernameExists
		}
	}

	if req.Email != existingCustomer.Email {
		emailTaken, _ := cs.CustomerRepository.GetByEmail(ctx, req.Email)
		if emailTaken != nil {
			return nil, entities.ErrCustomerEmailExists
		}
	}

	existingCustomer.Username = req.Username
	existingCustomer.Email = req.Email

	if err := existingCustomer.Validate(); err != nil {
		return nil, err
	}

	updatedCustomer, err := cs.CustomerRepository.Update(ctx, existingCustomer)
	if err != nil {
		return nil, err
	}

	cs.CustomerCacheRepository.Delete(ctx, cs.CustomerCacheRepository.GetCustomerCacheKey(updatedCustomer.ID))
	cs.CustomerCacheRepository.Delete(ctx, cs.CustomerCacheRepository.GetCustomerUsernameCacheKey(updatedCustomer.Username))
	cs.CustomerCacheRepository.Delete(ctx, cs.CustomerCacheRepository.GetCustomerEmailCacheKey(updatedCustomer.Email))

	return cs.mapToCustomerResponse(updatedCustomer), nil
}

func (cs *CustomerService) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	existingCustomer, err := cs.CustomerRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingCustomer == nil {
		return entities.ErrCustomerNotFound
	}

	if err := cs.CustomerRepository.Delete(ctx, id); err != nil {
		return err
	}

	cs.CustomerCacheRepository.Delete(ctx, cs.CustomerCacheRepository.GetCustomerCacheKey(id))
	cs.CustomerCacheRepository.Delete(ctx, cs.CustomerCacheRepository.GetCustomerUsernameCacheKey(existingCustomer.Username))
	cs.CustomerCacheRepository.Delete(ctx, cs.CustomerCacheRepository.GetCustomerEmailCacheKey(existingCustomer.Email))

	return nil
}

func (cs *CustomerService) GetCustomer(ctx context.Context, id uuid.UUID) (*dto.CustomerResponse, error) {
	cacheKey := cs.CustomerCacheRepository.GetCustomerCacheKey(id)
	if cachedCustomer, _ := cs.CustomerCacheRepository.Get(ctx, cacheKey); cachedCustomer != nil {
		return cs.mapToCustomerResponse(cachedCustomer), nil
	}

	customer, err := cs.CustomerRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, entities.ErrCustomerNotFound
	}

	cs.CustomerCacheRepository.Set(ctx, cacheKey, customer, 30*time.Minute)

	return cs.mapToCustomerResponse(customer), nil
}

func (cs *CustomerService) GetCustomerByUsername(ctx context.Context, username string) (*dto.CustomerResponse, error) {
	cacheKey := cs.CustomerCacheRepository.GetCustomerUsernameCacheKey(username)
	if cachedCustomer, _ := cs.CustomerCacheRepository.Get(ctx, cacheKey); cachedCustomer != nil {
		return cs.mapToCustomerResponse(cachedCustomer), nil
	}

	customer, err := cs.CustomerRepository.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, entities.ErrCustomerNotFound
	}

	cs.CustomerCacheRepository.Set(ctx, cacheKey, customer, 30*time.Minute)

	return cs.mapToCustomerResponse(customer), nil
}

func (cs *CustomerService) GetCustomerByEmail(ctx context.Context, email string) (*dto.CustomerResponse, error) {
	cacheKey := cs.CustomerCacheRepository.GetCustomerEmailCacheKey(email)
	if cachedCustomer, _ := cs.CustomerCacheRepository.Get(ctx, cacheKey); cachedCustomer != nil {
		return cs.mapToCustomerResponse(cachedCustomer), nil
	}

	customer, err := cs.CustomerRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, entities.ErrCustomerNotFound
	}

	cs.CustomerCacheRepository.Set(ctx, cacheKey, customer, 30*time.Minute)

	return cs.mapToCustomerResponse(customer), nil
}

func (cs *CustomerService) ListCustomers(ctx context.Context, page, pageSize int32) (*dto.CustomerListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	customers, err := cs.CustomerRepository.List(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := cs.CustomerRepository.Count(ctx)
	if err != nil {
		return nil, err
	}

	customerResponses := make([]*dto.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = cs.mapToCustomerResponse(customer)
	}

	totalPages := int32((totalCount + int64(pageSize) - 1) / int64(pageSize))

	return &dto.CustomerListResponse{
		Customers:  customerResponses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (cs *CustomerService) SearchCustomers(ctx context.Context, req *dto.CustomerSearchRequest) (*dto.CustomerListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize
	customers, err := cs.CustomerRepository.Search(ctx, req.Query, req.PageSize, offset)
	if err != nil {
		return nil, err
	}

	totalCount := int64(len(customers))

	customerResponses := make([]*dto.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = cs.mapToCustomerResponse(customer)
	}

	totalPages := int32((totalCount + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &dto.CustomerListResponse{
		Customers:  customerResponses,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (cs *CustomerService) mapToCustomerResponse(customer *entities.Customer) *dto.CustomerResponse {
	return &dto.CustomerResponse{
		ID:        customer.ID,
		Username:  customer.Username,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}

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
	customerRepository      interfaces.CustomerRepository      // lowercase - manual injection
	customerCacheRepository interfaces.CustomerCacheRepository // lowercase - manual injection
}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

// Method injection for lowercase fields
func (cs *CustomerService) SetDependencies(
	customerRepository interfaces.CustomerRepository,
	customerCacheRepository interfaces.CustomerCacheRepository,
) {
	cs.customerRepository = customerRepository
	cs.customerCacheRepository = customerCacheRepository
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

	existingCustomer, _ := cs.customerRepository.GetByUsername(ctx, req.Username)
	if existingCustomer != nil {
		return nil, entities.ErrCustomerUsernameExists
	}

	existingCustomer, _ = cs.customerRepository.GetByEmail(ctx, req.Email)
	if existingCustomer != nil {
		return nil, entities.ErrCustomerEmailExists
	}

	createdCustomer, err := cs.customerRepository.Create(ctx, customer)
	if err != nil {
		return nil, err
	}

	cs.customerCacheRepository.Set(ctx, cs.customerCacheRepository.GetCustomerCacheKey(createdCustomer.ID), createdCustomer, 30*time.Minute)
	cs.customerCacheRepository.Set(ctx, cs.customerCacheRepository.GetCustomerUsernameCacheKey(createdCustomer.Username), createdCustomer, 30*time.Minute)
	cs.customerCacheRepository.Set(ctx, cs.customerCacheRepository.GetCustomerEmailCacheKey(createdCustomer.Email), createdCustomer, 30*time.Minute)

	return cs.mapToCustomerResponse(createdCustomer), nil
}

func (cs *CustomerService) UpdateCustomer(ctx context.Context, id uuid.UUID, req *dto.UpdateCustomerRequest) (*dto.CustomerResponse, error) {
	existingCustomer, err := cs.customerRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingCustomer == nil {
		return nil, entities.ErrCustomerNotFound
	}

	if req.Username != existingCustomer.Username {
		usernameTaken, _ := cs.customerRepository.GetByUsername(ctx, req.Username)
		if usernameTaken != nil {
			return nil, entities.ErrCustomerUsernameExists
		}
	}

	if req.Email != existingCustomer.Email {
		emailTaken, _ := cs.customerRepository.GetByEmail(ctx, req.Email)
		if emailTaken != nil {
			return nil, entities.ErrCustomerEmailExists
		}
	}

	existingCustomer.Username = req.Username
	existingCustomer.Email = req.Email

	if err := existingCustomer.Validate(); err != nil {
		return nil, err
	}

	updatedCustomer, err := cs.customerRepository.Update(ctx, existingCustomer)
	if err != nil {
		return nil, err
	}

	cs.customerCacheRepository.Delete(ctx, cs.customerCacheRepository.GetCustomerCacheKey(updatedCustomer.ID))
	cs.customerCacheRepository.Delete(ctx, cs.customerCacheRepository.GetCustomerUsernameCacheKey(updatedCustomer.Username))
	cs.customerCacheRepository.Delete(ctx, cs.customerCacheRepository.GetCustomerEmailCacheKey(updatedCustomer.Email))

	return cs.mapToCustomerResponse(updatedCustomer), nil
}

func (cs *CustomerService) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	existingCustomer, err := cs.customerRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingCustomer == nil {
		return entities.ErrCustomerNotFound
	}

	if err := cs.customerRepository.Delete(ctx, id); err != nil {
		return err
	}

	cs.customerCacheRepository.Delete(ctx, cs.customerCacheRepository.GetCustomerCacheKey(id))
	cs.customerCacheRepository.Delete(ctx, cs.customerCacheRepository.GetCustomerUsernameCacheKey(existingCustomer.Username))
	cs.customerCacheRepository.Delete(ctx, cs.customerCacheRepository.GetCustomerEmailCacheKey(existingCustomer.Email))

	return nil
}

func (cs *CustomerService) GetCustomer(ctx context.Context, id uuid.UUID) (*dto.CustomerResponse, error) {
	cacheKey := cs.customerCacheRepository.GetCustomerCacheKey(id)
	if cachedCustomer, _ := cs.customerCacheRepository.Get(ctx, cacheKey); cachedCustomer != nil {
		return cs.mapToCustomerResponse(cachedCustomer), nil
	}

	customer, err := cs.customerRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, entities.ErrCustomerNotFound
	}

	cs.customerCacheRepository.Set(ctx, cacheKey, customer, 30*time.Minute)

	return cs.mapToCustomerResponse(customer), nil
}

func (cs *CustomerService) GetCustomerByUsername(ctx context.Context, username string) (*dto.CustomerResponse, error) {
	cacheKey := cs.customerCacheRepository.GetCustomerUsernameCacheKey(username)
	if cachedCustomer, _ := cs.customerCacheRepository.Get(ctx, cacheKey); cachedCustomer != nil {
		return cs.mapToCustomerResponse(cachedCustomer), nil
	}

	customer, err := cs.customerRepository.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, entities.ErrCustomerNotFound
	}

	cs.customerCacheRepository.Set(ctx, cacheKey, customer, 30*time.Minute)

	return cs.mapToCustomerResponse(customer), nil
}

func (cs *CustomerService) GetCustomerByEmail(ctx context.Context, email string) (*dto.CustomerResponse, error) {
	cacheKey := cs.customerCacheRepository.GetCustomerEmailCacheKey(email)
	if cachedCustomer, _ := cs.customerCacheRepository.Get(ctx, cacheKey); cachedCustomer != nil {
		return cs.mapToCustomerResponse(cachedCustomer), nil
	}

	customer, err := cs.customerRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, entities.ErrCustomerNotFound
	}

	cs.customerCacheRepository.Set(ctx, cacheKey, customer, 30*time.Minute)

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
	customers, err := cs.customerRepository.List(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := cs.customerRepository.Count(ctx)
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
	customers, err := cs.customerRepository.Search(ctx, req.Query, req.PageSize, offset)
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

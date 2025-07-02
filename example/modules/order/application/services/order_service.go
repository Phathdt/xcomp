package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"example/modules/order/application/dto"
	"example/modules/order/domain/entities"
	"example/modules/order/domain/interfaces"

	"xcomp"

	"github.com/google/uuid"
)

type OrderService struct {
	OrderRepo      interfaces.OrderRepository      `inject:"OrderRepository"`
	OrderItemRepo  interfaces.OrderItemRepository  `inject:"OrderItemRepository"`
	OrderCacheRepo interfaces.OrderCacheRepository `inject:"OrderCacheRepository"`
	Logger         xcomp.Logger                    `inject:"Logger"`
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	s.Logger.Info("Creating order",
		xcomp.Field("customer_id", req.CustomerID),
		xcomp.Field("items_count", len(req.Items)))

	order := entities.NewOrder(req.CustomerID)
	order.ShippingAddress = req.ShippingAddress
	order.BillingAddress = req.BillingAddress
	order.Notes = req.Notes

	for _, itemReq := range req.Items {
		err := order.AddItem(itemReq.ProductID, itemReq.ProductName, itemReq.Quantity, itemReq.UnitPrice)
		if err != nil {
			return nil, err
		}
	}

	order.CalculateTotal()

	if err := order.Validate(); err != nil {
		return nil, err
	}

	if err := s.OrderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	for _, item := range order.OrderItems {
		if err := s.OrderItemRepo.Create(ctx, item); err != nil {
			return nil, err
		}
	}

	response := dto.ToOrderResponse(order)
	return &response, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	s.Logger.Info("Getting order by ID", xcomp.Field("order_id", id))

	order, err := s.OrderCacheRepo.Get(ctx, id)
	if err != nil {
		order, err = s.OrderRepo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		items, err := s.OrderItemRepo.GetByOrderID(ctx, id)
		if err != nil {
			return nil, err
		}
		order.OrderItems = items

		if setErr := s.OrderCacheRepo.Set(ctx, order, 5*time.Minute); setErr != nil {
			s.Logger.Warn("Failed to cache order",
				xcomp.Field("order_id", id),
				xcomp.Field("error", setErr))
		}
	} else if order == nil {
		order, err = s.OrderRepo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		items, err := s.OrderItemRepo.GetByOrderID(ctx, id)
		if err != nil {
			return nil, err
		}
		order.OrderItems = items

		if setErr := s.OrderCacheRepo.Set(ctx, order, 5*time.Minute); setErr != nil {
			log.Printf("Failed to cache order: %v", setErr)
		}
	}

	response := dto.ToOrderResponse(order)
	return &response, nil
}

func (s *OrderService) GetOrdersByCustomerID(ctx context.Context, customerID uuid.UUID, page, pageSize int32) (*dto.OrderListResponse, error) {
	s.Logger.Info("Getting orders for customer",
		xcomp.Field("customer_id", customerID),
		xcomp.Field("page", page),
		xcomp.Field("page_size", pageSize))

	offset := (page - 1) * pageSize
	orders, err := s.OrderRepo.GetByCustomerID(ctx, customerID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		items, err := s.OrderItemRepo.GetByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		order.OrderItems = items
	}

	total, err := s.OrderRepo.CountByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	response := dto.ToOrderListResponse(orders, total, page, pageSize)
	return &response, nil
}

func (s *OrderService) GetAllOrders(ctx context.Context, page, pageSize int32) (*dto.OrderListResponse, error) {
	log.Printf("OrderService: Getting all orders")

	offset := (page - 1) * pageSize
	orders, err := s.OrderRepo.GetAll(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		items, err := s.OrderItemRepo.GetByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		order.OrderItems = items
	}

	total, err := s.OrderRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	response := dto.ToOrderListResponse(orders, total, page, pageSize)
	return &response, nil
}

func (s *OrderService) GetOrdersByStatus(ctx context.Context, status entities.OrderStatus, page, pageSize int32) (*dto.OrderListResponse, error) {
	log.Printf("OrderService: Getting orders by status %s", status)

	// Check if dependencies are properly injected
	if s.OrderRepo == nil {
		return nil, fmt.Errorf("orderRepo is not injected")
	}
	if s.OrderItemRepo == nil {
		return nil, fmt.Errorf("orderItemRepo is not injected")
	}

	offset := (page - 1) * pageSize
	orders, err := s.OrderRepo.GetByStatus(ctx, status, pageSize, offset)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		if order == nil {
			continue
		}
		items, err := s.OrderItemRepo.GetByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		order.OrderItems = items
	}

	total, err := s.OrderRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	response := dto.ToOrderListResponse(orders, total, page, pageSize)
	return &response, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, id uuid.UUID, req dto.UpdateOrderRequest) (*dto.OrderResponse, error) {
	s.Logger.Info("Updating order", xcomp.Field("order_id", id))

	order, err := s.OrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		order.Status = *req.Status
	}
	if req.ShippingCost != nil {
		order.ShippingCost = *req.ShippingCost
	}
	if req.TaxAmount != nil {
		order.TaxAmount = *req.TaxAmount
	}
	if req.DiscountAmount != nil {
		order.DiscountAmount = *req.DiscountAmount
	}
	if req.ShippingAddress != nil {
		order.ShippingAddress = req.ShippingAddress
	}
	if req.BillingAddress != nil {
		order.BillingAddress = req.BillingAddress
	}
	if req.Notes != nil {
		order.Notes = req.Notes
	}

	order.CalculateTotal()

	if err := order.Validate(); err != nil {
		return nil, err
	}

	if err := s.OrderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	items, err := s.OrderItemRepo.GetByOrderID(ctx, id)
	if err != nil {
		return nil, err
	}
	order.OrderItems = items

	response := dto.ToOrderResponse(order)
	return &response, nil
}

func (s *OrderService) ConfirmOrder(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	log.Printf("OrderService: Confirming order %s", id)

	order, err := s.OrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := order.ConfirmOrder(); err != nil {
		return nil, err
	}

	if err := s.OrderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	items, err := s.OrderItemRepo.GetByOrderID(ctx, id)
	if err != nil {
		return nil, err
	}
	order.OrderItems = items

	response := dto.ToOrderResponse(order)
	return &response, nil
}

func (s *OrderService) ShipOrder(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	log.Printf("OrderService: Shipping order %s", id)

	order, err := s.OrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := order.ShipOrder(); err != nil {
		return nil, err
	}

	if err := s.OrderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	items, err := s.OrderItemRepo.GetByOrderID(ctx, id)
	if err != nil {
		return nil, err
	}
	order.OrderItems = items

	response := dto.ToOrderResponse(order)
	return &response, nil
}

func (s *OrderService) DeliverOrder(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	log.Printf("OrderService: Delivering order %s", id)

	order, err := s.OrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := order.DeliverOrder(); err != nil {
		return nil, err
	}

	if err := s.OrderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	items, err := s.OrderItemRepo.GetByOrderID(ctx, id)
	if err != nil {
		return nil, err
	}
	order.OrderItems = items

	response := dto.ToOrderResponse(order)
	return &response, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	log.Printf("OrderService: Cancelling order %s", id)

	order, err := s.OrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := order.CancelOrder(); err != nil {
		return nil, err
	}

	if err := s.OrderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	items, err := s.OrderItemRepo.GetByOrderID(ctx, id)
	if err != nil {
		return nil, err
	}
	order.OrderItems = items

	response := dto.ToOrderResponse(order)
	return &response, nil
}

func (s *OrderService) AddOrderItem(ctx context.Context, orderID uuid.UUID, req dto.AddOrderItemRequest) (*dto.OrderResponse, error) {
	log.Printf("OrderService: Adding item to order %s", orderID)

	order, err := s.OrderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if err := order.AddItem(req.ProductID, req.ProductName, req.Quantity, req.UnitPrice); err != nil {
		return nil, err
	}

	order.CalculateTotal()

	if err := s.OrderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	newItem := order.OrderItems[len(order.OrderItems)-1]
	if err := s.OrderItemRepo.Create(ctx, newItem); err != nil {
		return nil, err
	}

	response := dto.ToOrderResponse(order)
	return &response, nil
}

func (s *OrderService) UpdateOrderItemQuantity(ctx context.Context, orderID, productID uuid.UUID, req dto.UpdateOrderItemQuantityRequest) (*dto.OrderResponse, error) {
	log.Printf("OrderService: Updating item quantity in order %s", orderID)

	order, err := s.OrderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	items, err := s.OrderItemRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order.OrderItems = items

	if err := order.UpdateItemQuantity(productID, req.Quantity); err != nil {
		return nil, err
	}

	order.CalculateTotal()

	if err := s.OrderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	for _, item := range order.OrderItems {
		if item.ProductID == productID {
			if err := s.OrderItemRepo.Update(ctx, item); err != nil {
				return nil, err
			}
			break
		}
	}

	response := dto.ToOrderResponse(order)
	return &response, nil
}

func (s *OrderService) RemoveOrderItem(ctx context.Context, orderID, productID uuid.UUID) (*dto.OrderResponse, error) {
	log.Printf("OrderService: Removing item from order %s", orderID)

	order, err := s.OrderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	items, err := s.OrderItemRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order.OrderItems = items

	var itemToRemove *entities.OrderItem
	for _, item := range order.OrderItems {
		if item.ProductID == productID {
			itemToRemove = item
			break
		}
	}

	if itemToRemove == nil {
		return nil, entities.ErrOrderItemNotFound
	}

	if err := order.RemoveItem(productID); err != nil {
		return nil, err
	}

	order.CalculateTotal()

	if err := s.OrderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	if err := s.OrderItemRepo.Delete(ctx, itemToRemove.ID); err != nil {
		return nil, err
	}

	response := dto.ToOrderResponse(order)
	return &response, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	log.Printf("OrderService: Deleting order %s", id)

	if err := s.OrderItemRepo.DeleteByOrderID(ctx, id); err != nil {
		return err
	}

	return s.OrderRepo.Delete(ctx, id)
}

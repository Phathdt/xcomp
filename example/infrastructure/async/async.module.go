package async

import (
	"context"
	"example/jobs"
	"example/modules/customer/domain/interfaces"
	orderInterfaces "example/modules/order/domain/interfaces"
	"example/processors"
	"example/schedulers"

	"xcomp"

	"fmt"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/redis/go-redis/v9"
)

type AsyncService struct {
	scheduler *schedulers.CheckPendingOrderScheduler
	server    *asynq.Server
	monitor   *asynqmon.HTTPHandler
	logger    xcomp.Logger
	processor *processors.CheckPendingOrderProcessor
}

func NewAsyncService(
	redisClient *redis.Client,
	orderService orderInterfaces.OrderService,
	customerService interfaces.CustomerService,
	logger xcomp.Logger,
) *AsyncService {
	redisAddr := redisClient.Options().Addr

	scheduler := schedulers.NewCheckPendingOrderScheduler(redisAddr, logger)

	processor := processors.NewCheckPendingOrderProcessor(
		orderService,
		customerService,
		logger,
	)

	redisOpt := asynq.RedisClientOpt{Addr: redisAddr}

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	monitor := asynqmon.New(asynqmon.Options{
		RootPath:     "/monitoring",
		RedisConnOpt: redisOpt,
	})

	return &AsyncService{
		scheduler: scheduler,
		server:    server,
		monitor:   monitor,
		logger:    logger,
		processor: processor,
	}
}

func (a *AsyncService) Start(ctx context.Context) error {
	a.logger.Info("Starting async service")

	mux := asynq.NewServeMux()
	mux.HandleFunc(jobs.TypeCheckPendingOrder, a.processor.ProcessCheckPendingOrder)

	go func() {
		if err := a.server.Run(mux); err != nil {
			a.logger.Error("Asynq server failed", xcomp.Field("error", err))
		}
	}()

	if err := a.scheduler.Start(ctx); err != nil {
		return err
	}

	a.logger.Info("Async service started successfully")
	return nil
}

func (a *AsyncService) Stop() {
	a.logger.Info("Stopping async service")

	if a.scheduler != nil {
		a.scheduler.Stop()
	}

	if a.server != nil {
		a.server.Shutdown()
	}

	a.logger.Info("Async service stopped")
}

func (a *AsyncService) GetMonitorHandler() *asynqmon.HTTPHandler {
	return a.monitor
}

func CreateAsyncModule() xcomp.Module {
	return xcomp.NewModule().
		AddFactory("AsyncService", func(c *xcomp.Container) any {
			redisClient, ok := c.Get("RedisClient").(*redis.Client)
			if !ok || redisClient == nil {
				panic("RedisClient not found or invalid type in container")
			}

			logger, ok := c.Get("Logger").(xcomp.Logger)
			if !ok || logger == nil {
				panic("Logger not found or invalid type in container")
			}

			// Debug: List all available services
			services := c.ListServices()
			logger.Info("Available services in container for AsyncService creation",
				xcomp.Field("services", services))

			// Get OrderService from container (from orderModule)
			orderServiceRaw := c.Get("OrderService")
			if orderServiceRaw == nil {
				logger.Error("OrderService is nil in container")
				panic("OrderService not found in container")
			}

			orderService, ok := orderServiceRaw.(orderInterfaces.OrderService)
			if !ok {
				logger.Error("OrderService has wrong type",
					xcomp.Field("actual_type", fmt.Sprintf("%T", orderServiceRaw)))
				panic("OrderService has invalid type in container")
			}

			logger.Info("Successfully retrieved OrderService from container",
				xcomp.Field("orderService_type", fmt.Sprintf("%T", orderService)),
				xcomp.Field("orderService_pointer", fmt.Sprintf("%p", orderService)))

			// Debug: Check if we can access the concrete type to inspect its fields
			logger.Info("Investigating OrderService dependency injection status")

			customerService, ok := c.Get("CustomerService").(interfaces.CustomerService)
			if !ok || customerService == nil {
				panic("CustomerService not found or invalid type in container")
			}

			logger.Info("Creating AsyncService with dependencies",
				xcomp.Field("redisAddr", redisClient.Options().Addr))

			asyncService := NewAsyncService(redisClient, orderService, customerService, logger)
			return asyncService
		}).
		Build()
}

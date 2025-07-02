package processors

import (
	"context"
	"encoding/json"
	"example/jobs"
	"example/modules/customer/domain/interfaces"
	orderInterfaces "example/modules/order/domain/interfaces"

	"xcomp"

	"fmt"

	"github.com/hibiken/asynq"
)

type CheckPendingOrderProcessor struct {
	orderService    orderInterfaces.OrderService
	customerService interfaces.CustomerService
	logger          xcomp.Logger
}

func NewCheckPendingOrderProcessor(
	orderService orderInterfaces.OrderService,
	customerService interfaces.CustomerService,
	logger xcomp.Logger,
) *CheckPendingOrderProcessor {
	return &CheckPendingOrderProcessor{
		orderService:    orderService,
		customerService: customerService,
		logger:          logger,
	}
}

func (p *CheckPendingOrderProcessor) ProcessCheckPendingOrder(ctx context.Context, t *asynq.Task) error {
	var job jobs.CheckPendingOrderJob
	if err := json.Unmarshal(t.Payload(), &job); err != nil {
		p.logger.Error("Failed to unmarshal check pending order job",
			xcomp.Field("error", err))
		return err
	}

	p.logger.Info("Processing check pending order job",
		xcomp.Field("job_created_at", job.CreatedAt),
		xcomp.Field("orderService_pointer", fmt.Sprintf("%p", p.orderService)))

	return nil
}

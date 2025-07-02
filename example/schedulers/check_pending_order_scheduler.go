package schedulers

import (
	"context"
	"example/jobs"
	"time"

	"xcomp"

	"github.com/hibiken/asynq"
)

type CheckPendingOrderScheduler struct {
	client asynq.Client
	logger xcomp.Logger
	ticker *time.Ticker
	done   chan bool
}

func NewCheckPendingOrderScheduler(redisAddr string, logger xcomp.Logger) *CheckPendingOrderScheduler {
	redisOpt := asynq.RedisClientOpt{Addr: redisAddr}
	client := asynq.NewClient(redisOpt)

	return &CheckPendingOrderScheduler{
		client: *client,
		logger: logger,
		done:   make(chan bool),
	}
}

func (s *CheckPendingOrderScheduler) Start(ctx context.Context) error {
	s.logger.Info("Starting CheckPendingOrderScheduler")

	s.ticker = time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.logger.Info("CheckPendingOrderScheduler stopped due to context cancellation")
				return
			case <-s.done:
				s.logger.Info("CheckPendingOrderScheduler stopped")
				return
			case <-s.ticker.C:
				if err := s.enqueueCheckPendingOrderJob(); err != nil {
					s.logger.Error("Failed to enqueue check pending order job",
						xcomp.Field("error", err))
				}
			}
		}
	}()

	return nil
}

func (s *CheckPendingOrderScheduler) Stop() {
	s.logger.Info("Stopping CheckPendingOrderScheduler")

	if s.ticker != nil {
		s.ticker.Stop()
	}

	close(s.done)
	s.client.Close()
}

func (s *CheckPendingOrderScheduler) enqueueCheckPendingOrderJob() error {
	job := jobs.NewCheckPendingOrderJob()
	payload, err := job.Payload()
	if err != nil {
		return err
	}

	task := asynq.NewTask(jobs.TypeCheckPendingOrder, payload)
	info, err := s.client.Enqueue(task)
	if err != nil {
		return err
	}

	s.logger.Debug("Enqueued check pending order job",
		xcomp.Field("task_id", info.ID),
		xcomp.Field("queue", info.Queue),
		xcomp.Field("created_at", job.CreatedAt))

	return nil
}

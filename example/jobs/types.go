package jobs

import (
	"encoding/json"
	"time"
)

const (
	TypeCheckPendingOrder = "check_pending_order"
)

type CheckPendingOrderJob struct {
	CreatedAt time.Time `json:"created_at"`
}

func NewCheckPendingOrderJob() *CheckPendingOrderJob {
	return &CheckPendingOrderJob{
		CreatedAt: time.Now(),
	}
}

func (j *CheckPendingOrderJob) Payload() ([]byte, error) {
	return json.Marshal(j)
}

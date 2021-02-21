package queue

import (
	"context"

	queueapiv1 "github.com/WinnersonKharsunai/GraduationProject/server/pkg/queue-api/go"
)

// Queue ...
type Queue struct {
	qSvc queueapiv1.QueueServiceClient
}

// ImqQueueIF ...
type ImqQueueIF interface {
	SendMessage(ctx context.Context, msg string) error
}

// NewQueue ...
func NewQueue(qSvc queueapiv1.QueueServiceClient) ImqQueueIF {
	return &Queue{
		qSvc: qSvc,
	}
}

// SendMessage ...
func (q *Queue) SendMessage(ctx context.Context, msg string) error {
	_, err := q.qSvc.SendMessage(ctx, &queueapiv1.SendMessageRequest{})
	if err != nil {
		return err
	}
	return nil
}

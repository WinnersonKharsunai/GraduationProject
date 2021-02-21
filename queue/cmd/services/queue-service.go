package services

import (
	"context"

	queueapi "github.com/WinnersonKharsunai/GraduationProject/queue/pkg/queue-api/go"
	"github.com/sirupsen/logrus"
)

// QueueService is the concrete implementation for the Queue Service
type QueueService struct {
	log *logrus.Logger
}

// NewQueueService is the factory function for the QueueService
func NewQueueService(log *logrus.Logger) *QueueService {
	return &QueueService{
		log: log,
	}
}

// CreateTopic ...
func (qs *QueueService) CreateTopic(ctx context.Context, in *queueapi.CreateTopicRequest) (*queueapi.CreateTopicResponse, error) {
	panic("not implemented")
}

// DeleteTopic ...
func (qs *QueueService) DeleteTopic(ctx context.Context, in *queueapi.DeleteTopicRequest) (*queueapi.DeleteTopicResponse, error) {
	panic("not implemented")
}

// SendMessage ...
func (qs *QueueService) SendMessage(ctx context.Context, in *queueapi.SendMessageRequest) (*queueapi.SendMessageResponse, error) {
	panic("not implemented")
}

// RetrieveMessage ...
func (qs *QueueService) RetrieveMessage(ctx context.Context, in *queueapi.RetrieveMessageRequest) (*queueapi.RetrieveMessageResponse, error) {
	panic("not implemented")
}

// DeleteMessage ...
func (qs *QueueService) DeleteMessage(ctx context.Context, in *queueapi.DeleteMessageRequest) (*queueapi.RetrieveMessageResponse, error) {
	panic("not implemented")
}

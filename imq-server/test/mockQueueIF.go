package test

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/queue"
)

// MockQueueIF is a struct for mocking ImqQueueIF
type MockQueueIF struct {
	Mock
	queue.ImqQueueIF
}

// SendMessage mocks on ImqQueueIF.SendMessage
func (mk *MockQueueIF) SendMessage(ctx context.Context, message queue.SendMessageRequest) error {
	args := mk.Called(ctx, message)
	return args.Error(0)
}

// RetrieveMessage mocks on ImqQueueIF.RetrieveMessage
func (mk *MockQueueIF) RetrieveMessage(ctx context.Context, topicID string) (*queue.Message, error) {
	args := mk.Called(ctx, topicID)
	return args.Get(0).(*queue.Message), args.Error(1)
}

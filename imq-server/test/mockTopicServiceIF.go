package test

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/domain"
)

// MockTopicServiceIF is a struct for mocking ImqQueueIF
type MockTopicServiceIF struct {
	Mock
	domain.TopicServicesIF
}

// GetTopics mocks on TopicServiceIF.GetTopics
func (m *MockTopicServiceIF) GetTopics(ctx context.Context, publisherID int) (*[]string, error) {
	args := m.Called(ctx, publisherID)
	return args.Get(0).(*[]string), args.Error(1)
}

// RegisterPublisherToTopic mocks on TopicServiceIF.RegisterPublisherToTopic
func (m *MockTopicServiceIF) RegisterPublisherToTopic(ctx context.Context, publisherID int, topicName string) error {
	args := m.Called(ctx, publisherID)
	return args.Error(0)
}

// DeregisterPublisherFromTopic mocks on TopicServiceIF.DeregisterPublisherFromTopic
func (m *MockTopicServiceIF) DeregisterPublisherFromTopic(ctx context.Context, publisherID int) error {
	args := m.Called(ctx, publisherID)
	return args.Error(0)
}

// AddMessageToTopic mocks on TopicServiceIF.AddMessageToTopic
func (m *MockTopicServiceIF) AddMessageToTopic(ctx context.Context, publisherID int, message domain.Message) error {
	args := m.Called(ctx, publisherID, message)
	return args.Error(0)
}

<<<<<<< HEAD
=======
// GetMessageStatus mocks on TopicServiceIF.GetMessageStatus
func (m *MockTopicServiceIF) GetMessageStatus(ctx context.Context, topic domain.Topic) (string, error) {
	args := m.Called(ctx, topic)
	return args.String(0), args.Error(1)
}

>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469
// GetMessage mocks on TopicServiceIF.GetMessage
func (m *MockTopicServiceIF) GetMessage(ctx context.Context, subscriberID int, topicName string) (*domain.Message, error) {
	args := m.Called(ctx, subscriberID, topicName)
	return args.Get(0).(*domain.Message), args.Error(1)
}

// RegisterSubscriberToTopic mocks on TopicServiceIF.RegisterSubscriberToTopic
func (m *MockTopicServiceIF) RegisterSubscriberToTopic(ctx context.Context, subscriberID int, topicName string) error {
	args := m.Called(ctx, subscriberID, topicName)
	return args.Error(0)
}

// DeregisterSubscriberFromTopic mocks on TopicServiceIF.DeregisterSubscriberFromTopic
func (m *MockTopicServiceIF) DeregisterSubscriberFromTopic(ctx context.Context, subscriberID int, topicName string) error {
	args := m.Called(ctx, subscriberID, topicName)
	return args.Error(0)
}

// GetRegisteredTopic mocks on TopicServiceIF.GetRegisteredTopic
func (m *MockTopicServiceIF) GetRegisteredTopic(ctx context.Context, subscriberID int) (*[]string, error) {
	args := m.Called(ctx, subscriberID)
	return args.Get(0).(*[]string), args.Error(1)
}

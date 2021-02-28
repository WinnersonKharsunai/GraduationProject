package test

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/publisher"
)

// MockPublisherIF is a struct for mocking PublisherIF
type MockPublisherIF struct {
	Mock
	publisher.PublisherIF
}

// ShowTopics mocks on PublisherIF.ShowTopics
func (m *MockPublisherIF) ShowTopics(ctx context.Context, in *publisher.ShowTopicRequest) (*publisher.ShowTopicResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*publisher.ShowTopicResponse), args.Error(1)
}

// ConnectToTopic mocks on PublisherIF.ConnectToTopic
func (m *MockPublisherIF) ConnectToTopic(ctx context.Context, in *publisher.ConnectToTopicRequest) (*publisher.ConnectToTopicResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*publisher.ConnectToTopicResponse), args.Error(1)
}

// DisconnectFromTopic mocks on PublisherIF.DisconnectFromTopic
func (m *MockPublisherIF) DisconnectFromTopic(ctx context.Context, in *publisher.DisconnectFromTopicRequest) (*publisher.DisconnectFromTopicResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*publisher.DisconnectFromTopicResponse), args.Error(1)
}

// PublishMessage mocks on PublisherIF.PublishMessage
func (m *MockPublisherIF) PublishMessage(ctx context.Context, in *publisher.PublishMessageRequest) (*publisher.PublishMessageResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*publisher.PublishMessageResponse), args.Error(1)
}

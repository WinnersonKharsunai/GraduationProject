package test

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/subscriber"
)

// MockSubscriberIF is a struct for mocking MSubscriberIF
type MockSubscriberIF struct {
	Mock
	subscriber.SubscriberIF
}

// ShowTopics mocks on SubscriberIF.ShowTopics
func (m *MockSubscriberIF) ShowTopics(ctx context.Context, in *subscriber.ShowTopicRequest) (*subscriber.ShowTopicResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*subscriber.ShowTopicResponse), args.Error(1)
}

// SubscribeToTopic mocks on SubscriberIF.SubscribeToTopic
func (m *MockSubscriberIF) SubscribeToTopic(ctx context.Context, in *subscriber.SubscribeToTopicRequest) (*subscriber.SubscribeToTopicResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*subscriber.SubscribeToTopicResponse), args.Error(1)
}

// UnsubscribeFromTopic mocks on SubscriberIF.UnsubscribeFromTopic
func (m *MockSubscriberIF) UnsubscribeFromTopic(ctx context.Context, in *subscriber.UnsubscribeFromTopicRequest) (*subscriber.UnsubscribeFromTopicResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*subscriber.UnsubscribeFromTopicResponse), args.Error(1)
}

// GetSubscribedTopics mocks on SubscriberIF.GetSubscribedTopics
func (m *MockSubscriberIF) GetSubscribedTopics(ctx context.Context, in *subscriber.GetSubscribedTopicsRequest) (*subscriber.GetSubscribedTopicsResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*subscriber.GetSubscribedTopicsResponse), args.Error(1)
}

// GetMessageFromTopic mocks on SubscriberIF.GetMessageFromTopic
func (m *MockSubscriberIF) GetMessageFromTopic(ctx context.Context, in *subscriber.GetMessageFromTopicRequest) (*subscriber.GetMessageFromTopicResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*subscriber.GetMessageFromTopicResponse), args.Error(1)
}

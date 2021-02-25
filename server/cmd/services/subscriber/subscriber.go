package subscriber

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/domain"
	"github.com/sirupsen/logrus"
)

// Subscriber is the concrete implementation for the Subscriber
type Subscriber struct {
	log      *logrus.Logger
	topicSvc domain.TopicServicesIF
}

// SubscriberIF is the interface to be for the Subscriber service
type SubscriberIF interface {
	SubscribeToTopic(ctx context.Context, in *SubscribeToTopicRequest) (*SubscribeToTopicResponse, error)
	UnsubscribeFromTopic(ctx context.Context, in *UnsubscribeFromTopicRequest) (*UnsubscribeFromTopicResponse, error)
	GetSubscribedTopics(ctx context.Context, in *GetSubscribedTopicsRequest) (*GetSubscribedTopicsResponse, error)
	GetMessageFromTopic(ctx context.Context, in *GetMessageFromTopicRequest) (*GetMessageFromTopicResponse, error)
}

// NewSubscriber is the factory function for the Subscriber
func NewSubscriber(log *logrus.Logger, topic domain.TopicServicesIF) SubscriberIF {
	return &Subscriber{
		log:      log,
		topicSvc: topic,
	}
}

// SubscribeToTopic ....
func (s *Subscriber) SubscribeToTopic(ctx context.Context, in *SubscribeToTopicRequest) (*SubscribeToTopicResponse, error) {
	return &SubscribeToTopicResponse{}, nil
}

// UnsubscribeFromTopic ...
func (s *Subscriber) UnsubscribeFromTopic(ctx context.Context, in *UnsubscribeFromTopicRequest) (*UnsubscribeFromTopicResponse, error) {
	return &UnsubscribeFromTopicResponse{}, nil
}

// GetSubscribedTopics ...
func (s *Subscriber) GetSubscribedTopics(ctx context.Context, in *GetSubscribedTopicsRequest) (*GetSubscribedTopicsResponse, error) {
	return &GetSubscribedTopicsResponse{}, nil
}

// GetMessageFromTopic ...
func (s *Subscriber) GetMessageFromTopic(ctx context.Context, in *GetMessageFromTopicRequest) (*GetMessageFromTopicResponse, error) {
	return &GetMessageFromTopicResponse{}, nil
}

package subscriber

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/domain"
	"github.com/sirupsen/logrus"
)

// Subscriber is the concrete implementation for the Subscriber
type Subscriber struct {
	log          *logrus.Logger
	topicService domain.TopicServicesIF
}

// SubscriberIF is the interface to be for the Subscriber service
type SubscriberIF interface {
	SubscribeToTopic(ctx context.Context, in *SubscribeToTopicRequest) (*SubscribeToTopicResponse, error)
	UnsubscribeFromTopic(ctx context.Context, in *UnsubscribeFromTopicRequest) (*UnsubscribeFromTopicResponse, error)
	GetSubscribedTopics(ctx context.Context, in *GetSubscribedTopicsRequest) (*GetSubscribedTopicsResponse, error)
	GetMessageFromTopic(ctx context.Context, in *GetMessageFromTopicRequest) (*GetMessageFromTopicResponse, error)
}

// NewSubscriber is the factory function for the Subscriber
func NewSubscriber(log *logrus.Logger, topicService domain.TopicServicesIF) SubscriberIF {
	return &Subscriber{
		log:          log,
		topicService: topicService,
	}
}

// SubscribeToTopic subscribes given subscriber to topic
func (s *Subscriber) SubscribeToTopic(ctx context.Context, in *SubscribeToTopicRequest) (*SubscribeToTopicResponse, error) {

	err := s.topicService.RegisterSubscriberToTopic(ctx, in.SubscriberID, in.TopicName)
	if err != nil {
		s.log.WithField("subscriberId", in.SubscriberID).Errorf("SubscribeToTopic: failed to register subscriber to topic: %v", err)
		return nil, err
	}

	return &SubscribeToTopicResponse{}, nil
}

// UnsubscribeFromTopic ubsubscribes given client from topic
func (s *Subscriber) UnsubscribeFromTopic(ctx context.Context, in *UnsubscribeFromTopicRequest) (*UnsubscribeFromTopicResponse, error) {

	err := s.topicService.DeregisterSubscriberFromTopic(ctx, in.SubscriberID, in.TopicName)
	if err != nil {
		s.log.WithField("subscriberId", in.SubscriberID).Errorf("UnsubscribeFromTopic: failed to deregister subscriber from topic: %v", err)
		return nil, err
	}

	return &UnsubscribeFromTopicResponse{}, nil
}

// GetSubscribedTopics fetches all the topics subscribed by given client
func (s *Subscriber) GetSubscribedTopics(ctx context.Context, in *GetSubscribedTopicsRequest) (*GetSubscribedTopicsResponse, error) {

	var getSubscribedTopicsResponse *GetSubscribedTopicsResponse

	topics, err := s.topicService.GetRegisteredTopic(ctx, in.SubscriberID)
	if err != nil {
		s.log.WithField("subscriberId", in.SubscriberID).Errorf("GetSubscribedTopics: failed to get subscribed topic: %v", err)
		return nil, err
	}

	for _, topic := range *topics {
		getSubscribedTopicsResponse.Topics = append(getSubscribedTopicsResponse.Topics, topic)
	}

	return getSubscribedTopicsResponse, nil
}

// GetMessageFromTopic fetches message publised to a given topic
func (s *Subscriber) GetMessageFromTopic(ctx context.Context, in *GetMessageFromTopicRequest) (*GetMessageFromTopicResponse, error) {

	var getMessageFromTopicResponse *GetMessageFromTopicResponse

	message, err := s.topicService.GetMessage(ctx, in.SubscriberID, in.TopicName)
	if err != nil {
		s.log.WithField("subscriberId", in.SubscriberID).Errorf("GetMessageFromTopic: failed to get message from topic: %v", err)
		return nil, err
	}

	getMessageFromTopicResponse.Message = Message{
		Data:      message.Data,
		CretedAt:  message.CretedAt,
		ExpiresAt: message.ExpiresAt,
	}

	return getMessageFromTopicResponse, nil
}

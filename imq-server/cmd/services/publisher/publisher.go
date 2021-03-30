package publisher

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Publisher is the concrete implementation for the Publisher
type Publisher struct {
	log          *logrus.Logger
	topicService domain.TopicServicesIF
}

// PublisherIF is the interface for the Publisher service
type PublisherIF interface {
	ShowTopics(ctx context.Context, in *ShowTopicRequest) (*ShowTopicResponse, error)
	ConnectToTopic(ctx context.Context, in *ConnectToTopicRequest) (*ConnectToTopicResponse, error)
	DisconnectFromTopic(ctx context.Context, in *DisconnectFromTopicRequest) (*DisconnectFromTopicResponse, error)
	PublishMessage(ctx context.Context, in *PublishMessageRequest) (*PublishMessageResponse, error)
}

// NewPublisher is the factory function for the Publisher type
func NewPublisher(log *logrus.Logger, topicService domain.TopicServicesIF) PublisherIF {
	return &Publisher{
		log:          log,
		topicService: topicService,
	}
}

// ShowTopics fetch all the topics that are available
func (p *Publisher) ShowTopics(ctx context.Context, in *ShowTopicRequest) (*ShowTopicResponse, error) {
	showTopicResponse := &ShowTopicResponse{}

	topics, err := p.topicService.GetTopics(ctx, in.PublisherID)
	if err != nil {
		p.log.WithField("publisherId", in.PublisherID).Errorf("ShowTopics: failed to get topics: %v", err)
		return nil, err
	}

	for _, topic := range *topics {
		showTopicResponse.Topics = append(showTopicResponse.Topics, topic)
	}

	return showTopicResponse, nil
}

// ConnectToTopic register publisher to topic
func (p *Publisher) ConnectToTopic(ctx context.Context, in *ConnectToTopicRequest) (*ConnectToTopicResponse, error) {
	connectToTopicResponse := &ConnectToTopicResponse{}

	err := p.topicService.RegisterPublisherToTopic(ctx, in.PublisherID, in.TopicName)
	if err != nil {
		p.log.WithField("publisherId", in.PublisherID).Errorf("ConnectToTopic: failed to add publisher to topic: %v", err)
		return nil, err
	}

	connectToTopicResponse.Status = sucessConnected

	return connectToTopicResponse, nil
}

// DisconnectFromTopic deregister publisher from topic
func (p *Publisher) DisconnectFromTopic(ctx context.Context, in *DisconnectFromTopicRequest) (*DisconnectFromTopicResponse, error) {
	disconnectFromTopicResponse := &DisconnectFromTopicResponse{}

	err := p.topicService.DeregisterPublisherFromTopic(ctx, in.PublisherID)
	if err != nil {
		p.log.WithField("publisherId", in.PublisherID).Errorf("DisconnectFromTopic: failed to remove publisher from topic: %v", err)
		return nil, err
	}

	disconnectFromTopicResponse.Status = statusDisconnected

	return disconnectFromTopicResponse, nil
}

// PublishMessage publishes new message to topic
func (p *Publisher) PublishMessage(ctx context.Context, in *PublishMessageRequest) (*PublishMessageResponse, error) {
	publishMessageResponse := &PublishMessageResponse{}

	msg := domain.Message{
		MessageID: uuid.New().String(),
		Data:      in.Message.Data,
		CretedAt:  in.Message.CretedAt,
		ExpiresAt: in.Message.ExpiresAt,
	}

	err := p.topicService.AddMessageToTopic(ctx, in.PublisherID, msg)
	if err != nil {
		p.log.WithField("publisherId", in.PublisherID).Errorf("PublishMessage: failed to add message to topic: %v", err)
		return nil, err
	}

	publishMessageResponse.Status = statusSuccessful

	return publishMessageResponse, nil
}

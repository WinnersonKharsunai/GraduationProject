package services

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/domain"
	"github.com/sirupsen/logrus"
)

// Publisher is the concrete implementation for the Publisher
type Publisher struct {
	log      *logrus.Logger
	topicSvc domain.TopicServicesIF
}

// PublisherIF is the interface for the Publisher service
type PublisherIF interface {
	ShowTopics(ctx context.Context, in *ShowTopicRequest) (*ShowTopicResponse, error)
	ConnectToTopic(ctx context.Context, in *ConnectToTopicRequest) (*ConnectToTopicResponse, error)
	DisconnectFromTopic(ctx context.Context, in *DisconnectFromTopicRequest) (*DisconnectFromTopicResponse, error)
	PublishMessage(ctx context.Context, in *PublishMessageRequest) (*PublishMessageResponse, error)
	CheckMessageStatus(ctx context.Context, in *CheckMessageStatusRequest) (*CheckMessageStatusResponse, error)
}

// NewPublisher is the factory function for the Publisher type
func NewPublisher(log *logrus.Logger, topic domain.TopicServicesIF) PublisherIF {
	return &Publisher{
		log:      log,
		topicSvc: topic,
	}
}

// ShowTopics ...
func (p *Publisher) ShowTopics(ctx context.Context, in *ShowTopicRequest) (*ShowTopicResponse, error) {
	return &ShowTopicResponse{}, nil
}

// ConnectToTopic ...
func (p *Publisher) ConnectToTopic(ctx context.Context, in *ConnectToTopicRequest) (*ConnectToTopicResponse, error) {
	return &ConnectToTopicResponse{}, nil
}

// DisconnectFromTopic ...
func (p *Publisher) DisconnectFromTopic(ctx context.Context, in *DisconnectFromTopicRequest) (*DisconnectFromTopicResponse, error) {
	return &DisconnectFromTopicResponse{}, nil
}

// PublishMessage ...
func (p *Publisher) PublishMessage(ctx context.Context, in *PublishMessageRequest) (*PublishMessageResponse, error) {
	return &PublishMessageResponse{}, nil
}

// CheckMessageStatus ...
func (p *Publisher) CheckMessageStatus(ctx context.Context, in *CheckMessageStatusRequest) (*CheckMessageStatusResponse, error) {
	return &CheckMessageStatusResponse{}, nil
}

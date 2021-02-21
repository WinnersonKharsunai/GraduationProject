package services

import (
	"context"
	"fmt"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/domain"
	"github.com/sirupsen/logrus"
)

// Publisher ...
type Publisher struct {
	log      *logrus.Logger
	topicSvc domain.TopicServicesIF
}

// PublisherIF ...
type PublisherIF interface {
	ShowTopics(ctx context.Context) ([]string, error)
	ConnectToTopic(ctx context.Context, publisherID int, topicName string) error
	DisconnectFromTopic(ctx context.Context, id int) error
	PublishMessage(ctx context.Context, msg string) error
	CheckMessageStatus(ctx context.Context, msg string) (string, error)
}

// NewPublisher ...
func NewPublisher(log *logrus.Logger, topic domain.TopicServicesIF) PublisherIF {
	return &Publisher{
		log:      log,
		topicSvc: topic,
	}
}

// ShowTopics ...
func (p *Publisher) ShowTopics(ctx context.Context) ([]string, error) {

	fmt.Println("HEllo")
	topics, err := p.topicSvc.GetTopics(ctx)
	if err != nil {
		return []string{}, err
	}
	return topics, nil
}

// ConnectToTopic ...
func (p *Publisher) ConnectToTopic(ctx context.Context, publisherID int, topicName string) error {
	err := p.topicSvc.AddPublisherToTopic(ctx, publisherID, topicName)
	if err != nil {
		return err
	}
	return nil
}

// DisconnectFromTopic ...
func (p *Publisher) DisconnectFromTopic(ctx context.Context, id int) error {
	err := p.topicSvc.RemovePublisherFromTopic(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// PublishMessage ...
func (p *Publisher) PublishMessage(ctx context.Context, msg string) error {
	err := p.topicSvc.AddMessageToTopic(ctx, domain.Topic{})
	if err != nil {
		return err
	}
	return nil
}

// CheckMessageStatus ...
func (p *Publisher) CheckMessageStatus(ctx context.Context, msg string) (string, error) {
	status, err := p.topicSvc.GetMessageStatus(ctx, domain.Topic{})
	if err != nil {
		return "", err
	}
	return status, nil
}

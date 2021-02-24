package services

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/domain"
	"github.com/sirupsen/logrus"
)

// Subscriber ...
type Subscriber struct {
	log      *logrus.Logger
	topicSvc domain.TopicServicesIF
}

// SubscriberIF ...
type SubscriberIF interface {
}

// NewSubscriber ...
func NewSubscriber(log *logrus.Logger, topic domain.TopicServicesIF) SubscriberIF {
	return &Subscriber{
		log:      log,
		topicSvc: topic,
	}
}

// PublishMessage ...
func (Subscriber) PublishMessage(ctx context.Context, msg string) error {

	return nil
}

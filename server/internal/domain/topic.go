package domain

import (
	"context"
	"errors"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/queue"
	"github.com/WinnersonKharsunai/GraduationProject/server/internal/storage"
	"github.com/sirupsen/logrus"
)

// TopicService implement the TopicService interface
type TopicService struct {
	log   *logrus.Logger
	db    storage.DatabaseIF
	queue queue.ImqQueueIF
}

// TopicServicesIF is the interaface of topic service
type TopicServicesIF interface {
	GetTopics(ctx context.Context, publisherID int) (*[]string, error)
	RegisterPublisherToTopic(ctx context.Context, publisherID int, topicName string) error
	DeregisterPublisherFromTopic(ctx context.Context, publisherID int) error
	AddMessageToTopic(ctx context.Context, publisherID int, message Message) error
	GetMessageStatus(ctx context.Context, topic Topic) (string, error)
	GetMessage(ctx context.Context, subscriberID int, topicName string) (*Message, error)
	RegisterSubscriberToTopic(ctx context.Context, subscriberID int, topicName string) error
	DeregisterSubscriberFromTopic(ctx context.Context, subscriberID int, topicName string) error
	GetRegisteredTopic(ctx context.Context, subscriberID int) (*[]string, error)
}

// NewTopic is the factory function for the TopicService type
func NewTopic(log *logrus.Logger, db storage.DatabaseIF, queue queue.ImqQueueIF) TopicServicesIF {
	return &TopicService{
		log:   log,
		db:    db,
		queue: queue,
	}
}

// GetTopics fetches all the available topics
func (t *TopicService) GetTopics(ctx context.Context, publisherID int) (*[]string, error) {

	topics, err := t.db.FetchAllTopics(ctx, publisherID)
	if err != nil {
		t.log.WithField("publisherId", publisherID).Errorf("GetTopics: failed to fetch all topics: %v", err)
		return nil, err
	}

	return topics, nil
}

// RegisterPublisherToTopic register publisher to topic
func (t *TopicService) RegisterPublisherToTopic(ctx context.Context, publisherID int, topicName string) error {

	currentTopicID, notFound, err := t.db.GetTopicIDFromPublisher(ctx, publisherID)
	if err != nil {
		t.log.WithField("publisherId", publisherID).Errorf("RegisterPublisherToTopic: failed to get topicId from Publisher: %v", err)
		return err
	}

	topicID, err := t.db.GetTopicIDFromTopic(ctx, topicName)
	if err != nil {
		t.log.WithField("publisherId", publisherID).Errorf("RegisterPublisherToTopic: failed to get topicId from Topic: %v", err)
		return err
	}

	if notFound {
		err = t.db.InsertPublisher(ctx, publisherID, topicName)
		if err != nil {
			t.log.WithField("publisherId", publisherID).Errorf("RegisterPublisherToTopic: failed to insert publisher to topic: %v", err)
			return err
		}
	}

	if currentTopicID != topicID {
		err := t.db.UpdatePublisherIDIntoPublisher(ctx, publisherID, topicID)
		if err != nil {
			t.log.WithField("publisherId", publisherID).Errorf("RegisterPublisherToTopic: failed to update publisher to topic: %v", err)
			return err
		}
	}

	return nil
}

// DeregisterPublisherFromTopic  deregister publisher from topic
func (t *TopicService) DeregisterPublisherFromTopic(ctx context.Context, publisherID int) error {

	err := t.db.RemoveTopicIDFromPublisher(ctx, publisherID)
	if err != nil {
		t.log.WithField("publisherId", publisherID).Errorf("DeregisterPublisherFromTopic: failed to insert publisher to topic: %v", err)
		return nil
	}

	return nil
}

// AddMessageToTopic publish the messaget to given topic
func (t *TopicService) AddMessageToTopic(ctx context.Context, publisherID int, message Message) error {

	topicID, notFound, err := t.db.GetTopicIDFromPublisher(ctx, publisherID)
	if err != nil {
		return err
	}

	if notFound {
		return errors.New("not subscribed to any topic")
	}

	sendMessageRequest := queue.SendMessageRequest{
		TopicID: topicID,
		Message: queue.Message{
			MessageID: message.MessageID,
			Data:      message.Data,
			CretedAt:  message.CretedAt,
			ExpiresAt: message.ExpiresAt,
		},
	}

	if err := t.queue.SendMessage(ctx, sendMessageRequest); err != nil {
		return err
	}

	msg := storage.Message{
		MessageID: message.MessageID,
		Data:      message.Data,
		CretedAt:  message.CretedAt,
		ExpiresAt: message.ExpiresAt,
	}

	err = t.db.InsertMessageIntoMessage(ctx, publisherID, topicID, msg)
	if err != nil {
		return err
	}

	return nil
}

// RegisterSubscriberToTopic add subscriber to the given topic
func (t *TopicService) RegisterSubscriberToTopic(ctx context.Context, subscriberID int, topicName string) error {

	topics, err := t.db.GetSubscribedTopics(ctx, subscriberID)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		if topicName == topic {
			return nil
		}
	}

	topicID, err := t.db.GetTopicIDFromTopic(ctx, topicName)
	if err != nil {
		return err
	}

	err = t.db.InsertSubscriberIDIntoSubscriber(ctx, subscriberID)
	if err != nil {
		return err
	}

	err = t.db.InsertIntoSubscriberTopicMap(ctx, subscriberID, topicID)
	if err != nil {
		return err
	}

	return nil
}

// DeregisterSubscriberFromTopic remove subscriber from topic
func (t *TopicService) DeregisterSubscriberFromTopic(ctx context.Context, subscriberID int, topicName string) error {

	topicID, err := t.db.GetTopicIDFromTopic(ctx, topicName)
	if err != nil {
		return err
	}

	err = t.db.RemoveTopicIDFromSubscriberTopicMap(ctx, subscriberID, topicID)
	if err != nil {
		return err
	}

	return nil
}

// GetRegisteredTopic fetches all the registered topics
func (t *TopicService) GetRegisteredTopic(ctx context.Context, subscriberID int) (*[]string, error) {

	topics, err := t.db.GetSubscribedTopics(ctx, subscriberID)
	if err != nil {
		return nil, err
	}

	return &topics, nil
}

// GetMessageStatus ...
func (t *TopicService) GetMessageStatus(ctx context.Context, topic Topic) (string, error) {
	return "", nil
}

// GetMessage fetch the message from the queue based on the given subscriberId
func (t *TopicService) GetMessage(ctx context.Context, subscriberID int, topicName string) (*Message, error) {

	msg := t.queue.RetrieveMessage(ctx, topicName)

	return &Message{
		MessageID: msg.MessageID,
		Data:      msg.Data,
		CretedAt:  msg.CretedAt,
		ExpiresAt: msg.ExpiresAt,
	}, nil
}

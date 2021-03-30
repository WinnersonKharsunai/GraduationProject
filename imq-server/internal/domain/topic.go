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
func (t *TopicService) GetTopics(ctx context.Context, id int) (*[]string, error) {
	topics, err := t.db.FetchAllTopics(ctx, id)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

// RegisterPublisherToTopic register publisher to topic
func (t *TopicService) RegisterPublisherToTopic(ctx context.Context, publisherID int, topicName string) error {
	topicID, notFound, err := t.db.GetTopicIDFromPublisher(ctx, publisherID)
	if err != nil {
		t.log.WithField("publisherId", publisherID).Errorf("RegisterPublisherToTopic: failed to get topicId from Publisher: %v", err)
		return err
	}

	if topicID != "" {
		err := errors.New("cannot register more than one topic at a time")
		t.log.WithField("publisherId", publisherID).Errorf("RegisterPublisherToTopic: already registered to topics: %v", err)
		return err
	}

	topicID, err = t.db.GetTopicIDFromTopic(ctx, topicName)
	if err != nil {
		t.log.WithField("publisherId", publisherID).Errorf("RegisterPublisherToTopic: failed to get topicId from Topic: %v", err)
		return err
	}

	if topicID == "" {
		err := errors.New("topic not found")
		t.log.WithField("publisherId", publisherID).Errorf("RegisterPublisherToTopic: failed to get topicId: %v", err)
		return errors.New("topic not found")
	}

	if notFound {
		err = t.db.InsertPublisher(ctx, publisherID, topicID)
		if err != nil {
			t.log.WithField("publisherId", publisherID).Errorf("RegisterPublisherToTopic: failed to insert publisher to topic: %v", err)
			return err
		}
	} else {
		err = t.db.UpdateTopicIDIntoPublisher(ctx, publisherID, topicID)
		if err != nil {
			t.log.WithField("publisherId", publisherID).Errorf("RegisterPublisherToTopic: failed to update topicId into Publisher: %v", err)
			return err
		}
	}

	return nil
}

// DeregisterPublisherFromTopic  deregister publisher from topic
func (t *TopicService) DeregisterPublisherFromTopic(ctx context.Context, publisherID int) error {
	_, notFound, err := t.db.GetTopicIDFromPublisher(ctx, publisherID)
	if err != nil {
		t.log.WithField("publisherId", publisherID).Errorf("RegisterPublisherToTopic: failed to get topicId from Publisher: %v", err)
		return err
	}

	if notFound {
		err := errors.New("you are not registered with any topic")
		t.log.WithField("publisherId", publisherID).Errorf("DeregisterPublisherFromTopic: record not found in Publisher: %v", err)
		return err
	}

	err = t.db.RemoveTopicIDFromPublisher(ctx, publisherID)
	if err != nil {
		t.log.WithField("publisherId", publisherID).Errorf("DeregisterPublisherFromTopic: failed to insert publisher to topic: %v", err)
		return err
	}

	return nil
}

// AddMessageToTopic publish the messaget to given topic
func (t *TopicService) AddMessageToTopic(ctx context.Context, publisherID int, message Message) error {
	topicID, notFound, err := t.db.GetTopicIDFromPublisher(ctx, publisherID)
	if err != nil {
		t.log.WithField("publisherId", publisherID).Errorf("AddMessageToTopic: failed to get topicId from publisher: %v", err)
		return err
	}

	if notFound {
		err := errors.New("you are not register to any topic")
		t.log.WithField("publisherId", publisherID).Errorf("AddMessageToTopic: failed to publish, not registered to any topic: %v", err)
		return err
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
		t.log.WithField("publisherId", publisherID).Errorf("AddMessageToTopic: failed to send message to queue: %v", err)
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
		t.log.WithField("publisherId", publisherID).Errorf("AddMessageToTopic: failed to store message: %v", err)
		return err
	}

	return nil
}

// RegisterSubscriberToTopic add subscriber to the given topic
func (t *TopicService) RegisterSubscriberToTopic(ctx context.Context, subscriberID int, topicName string) error {
	topics, err := t.db.GetSubscribedTopics(ctx, subscriberID)
	if err != nil {
		t.log.WithField("subscriberId", subscriberID).Errorf("RegisterSubscriberToTopic: failed to get subscribed topics: %v", err)
		return err
	}

	for _, topic := range topics {
		if topicName == topic {
			err := errors.New("you are already subscribed to this topic")
			t.log.WithField("subscriberId", subscriberID).Errorf("RegisterSubscriberToTopic: found subscribed topic: %v", err)
			return err
		}
	}

	topicID, err := t.db.GetTopicIDFromTopic(ctx, topicName)
	if err != nil {
		t.log.WithField("subscriberId", subscriberID).Errorf("RegisterSubscriberToTopic: failed to get topicId from topics: %v", err)
		return err
	}

	err = t.db.InsertSubscriberIDIntoSubscriber(ctx, subscriberID)
	if err != nil {
		t.log.WithField("subscriberId", subscriberID).Errorf("RegisterSubscriberToTopic: failed to save subscriberId: %v", err)
		return err
	}

	err = t.db.InsertIntoSubscriberTopicMap(ctx, subscriberID, topicID)
	if err != nil {
		t.log.WithField("subscriberId", subscriberID).Errorf("RegisterSubscriberToTopic: failed to save mapped topic for subscriber: %v", err)
		return err
	}

	return nil
}

// DeregisterSubscriberFromTopic remove subscriber from topic
func (t *TopicService) DeregisterSubscriberFromTopic(ctx context.Context, subscriberID int, topicName string) error {
	topicID, err := t.db.GetTopicIDFromTopic(ctx, topicName)
	if err != nil {
		t.log.WithField("subscriberId", subscriberID).Errorf("DeregisterSubscriberFromTopic: failed to get topicId from topic: %v", err)
		return err
	}

	if topicID == "" {
		err := errors.New("topic not found")
		t.log.WithField("subscriberId", subscriberID).Errorf("DeregisterSubscriberFromTopic: no record found for the given topic name: %v", err)
		return err
	}

	err = t.db.RemoveTopicIDFromSubscriberTopicMap(ctx, subscriberID, topicID)
	if err != nil {
		t.log.WithField("subscriberId", subscriberID).Errorf("DeregisterSubscriberFromTopic: failed to remove subscriber mapping to topic: %v", err)
		return err
	}

	return nil
}

// GetRegisteredTopic fetches all the registered topics
func (t *TopicService) GetRegisteredTopic(ctx context.Context, subscriberID int) (*[]string, error) {
	topics, err := t.db.GetSubscribedTopics(ctx, subscriberID)
	if err != nil {
		t.log.WithField("subscriberId", subscriberID).Errorf("GetRegisteredTopic: failed to get subscribed topics: %v", err)
		return nil, err
	}

	return &topics, nil
}

// GetMessage fetch the message from the queue based on the given subscriberId
func (t *TopicService) GetMessage(ctx context.Context, subscriberID int, topicName string) (*Message, error) {
	topicID, err := t.db.GetTopicIDFromTopic(ctx, topicName)
	if err != nil {
		t.log.WithField("subscriberId", subscriberID).Errorf("GetMessage: failed to get topicid from topic: %v", err)
		return nil, err
	}

	msg, err := t.queue.RetrieveMessage(ctx, topicID)
	if err != nil {
		t.log.WithField("subscriberId", subscriberID).Errorf("GetMessage: failed to retrieve message from queue: %v", err)
		return nil, err
	}

	message := Message{
		MessageID: msg.MessageID,
		Data:      msg.Data,
		CretedAt:  msg.CretedAt,
		ExpiresAt: msg.ExpiresAt,
	}

	return &message, nil
}

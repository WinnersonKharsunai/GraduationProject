package domain

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/queue"
	"github.com/WinnersonKharsunai/GraduationProject/server/internal/storage"
)

// TopicService ...
type TopicService struct {
	db    storage.DatabaseIF
	queue queue.ImqQueueIF
}

// TopicServicesIF ...
type TopicServicesIF interface {
	GetTopics(ctx context.Context) ([]string, error)
	AddPublisherToTopic(ctx context.Context, publisherID int, topicname string) error
	RemovePublisherFromTopic(ctx context.Context, id int) error
	AddMessageToTopic(ctx context.Context, topic Topic) error
	GetMessageStatus(ctx context.Context, topic Topic) (string, error)
}

// NewTopic ...
func NewTopic(db storage.DatabaseIF, queue queue.ImqQueueIF) TopicServicesIF {
	return &TopicService{
		db:    db,
		queue: queue,
	}
}

// GetTopics ...
func (t *TopicService) GetTopics(ctx context.Context) ([]string, error) {
	topics, err := t.db.FetchAllTopics(ctx)
	if err != nil {
		return []string{}, err
	}
	return topics, nil
}

// AddPublisherToTopic ...
func (t *TopicService) AddPublisherToTopic(ctx context.Context, publisherID int, topicName string) error {
	err := t.db.InsertPublisher(ctx, publisherID, topicName)
	if err != nil {
		return err
	}
	return nil
}

// RemovePublisherFromTopic ...
func (t *TopicService) RemovePublisherFromTopic(ctx context.Context, id int) error {
	err := t.db.RemovePublisher(ctx, id)
	if err != nil {
		return nil
	}
	return nil
}

// AddMessageToTopic ...
func (t *TopicService) AddMessageToTopic(ctx context.Context, topic Topic) error {

	return nil
}

// GetMessageStatus ...
func (t *TopicService) GetMessageStatus(ctx context.Context, topic Topic) (string, error) {
	return "", nil
}

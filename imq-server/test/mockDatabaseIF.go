package test

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/storage"
)

// MockDatabaseIF is a struct for mocking DatabaseIF
type MockDatabaseIF struct {
	Mock
	storage.DatabaseIF
}

// FetchAllTopics mocks on DatabaseIF.FetchAllTopics
func (m *MockDatabaseIF) FetchAllTopics(ctx context.Context, id int) (*[]string, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*[]string), args.Error(1)
}

// InsertPublisher mocks on DatabaseIF.InsertPublisher
func (m *MockDatabaseIF) InsertPublisher(ctx context.Context, publisherID int, topicName string) error {
	args := m.Called(ctx, publisherID, topicName)
	return args.Error(0)
}

// UpdateTopicIDIntoPublisher mocks on DatabaseIF.UpdateTopicIDIntoPublisher
func (m *MockDatabaseIF) UpdateTopicIDIntoPublisher(ctx context.Context, publisherID int, topicID string) error {
	args := m.Called(ctx, publisherID, topicID)
	return args.Error(0)
}

// RemoveTopicIDFromPublisher mocks on DatabaseIF.RemoveTopicIDFromPublisher
func (m *MockDatabaseIF) RemoveTopicIDFromPublisher(ctx context.Context, publisherID int) error {
	args := m.Called(ctx, publisherID)
	return args.Error(0)
}

// FetchQueues mocks on DatabaseIF.FetchQueues
func (m *MockDatabaseIF) FetchQueues(ctx context.Context) (*storage.Queue, error) {
	args := m.Called(ctx)
	return &storage.Queue{}, args.Error(1)
}

// GetTopicIDFromPublisher mocks on DatabaseIF.GetTopicIDFromPublisher
func (m *MockDatabaseIF) GetTopicIDFromPublisher(ctx context.Context, publisherID int) (string, bool, error) {
	args := m.Called(ctx, publisherID)
	return args.String(0), args.Bool(1), args.Error(2)
}

// GetTopicIDFromTopic mocks on DatabaseIF.GetTopicIDFromTopic
func (m *MockDatabaseIF) GetTopicIDFromTopic(ctx context.Context, topicName string) (string, error) {
	args := m.Called(ctx, topicName)
	return args.String(0), args.Error(1)
}

// InsertMessageIntoMessage mocks on DatabaseIF.InsertMessageIntoMessage
func (m *MockDatabaseIF) InsertMessageIntoMessage(ctx context.Context, publisherID int, topicID string, message storage.Message) error {
	args := m.Called(ctx, publisherID, topicID, message)
	return args.Error(0)
}

// GetSubscribedTopics mocks on DatabaseIF.GetSubscribedTopics
func (m *MockDatabaseIF) GetSubscribedTopics(ctx context.Context, subscriberID int) ([]string, error) {
	args := m.Called(ctx, subscriberID)
	return []string{}, args.Error(1)
}

// InsertSubscriberIDIntoSubscriber mocks on DatabaseIF.InsertSubscriberIDIntoSubscriber
func (m *MockDatabaseIF) InsertSubscriberIDIntoSubscriber(ctx context.Context, subscriberID int) error {
	args := m.Called(ctx, subscriberID)
	return args.Error(0)
}

// InsertIntoSubscriberTopicMap mocks on DatabaseIF.InsertIntoSubscriberTopicMap
func (m *MockDatabaseIF) InsertIntoSubscriberTopicMap(ctx context.Context, subscriberID int, topicID string) error {
	args := m.Called(ctx, subscriberID, topicID)
	return args.Error(0)
}

// RemoveTopicIDFromSubscriberTopicMap mocks on DatabaseIF.RemoveTopicIDFromSubscriberTopicMap
func (m *MockDatabaseIF) RemoveTopicIDFromSubscriberTopicMap(ctx context.Context, subscriberID int, topicID string) error {
	args := m.Called(ctx, subscriberID, topicID)
	return args.Error(0)
}

// SaveQueues mocks on DatabaseIF.SaveQueues
func (m *MockDatabaseIF) SaveQueues(ctx context.Context, queue *[]storage.StoreQueue, isLiveQueue bool) error {
	args := m.Called(ctx, queue, isLiveQueue)
	return args.Error(0)
}

// RemoveMessagesFromQueue mocks on DatabaseIF.RemoveMessagesFromQueue
func (m *MockDatabaseIF) RemoveMessagesFromQueue(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

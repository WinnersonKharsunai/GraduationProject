package domain_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/domain"
	"github.com/WinnersonKharsunai/GraduationProject/server/internal/storage"
	"github.com/WinnersonKharsunai/GraduationProject/server/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestGetTopics_Fail(t *testing.T) {
	publisherID := 5000
	expectedErr := errors.New("failed to store new client")

	mockDb := &test.MockDatabaseIF{}
	mockQueue := &test.MockQueueIF{}
	mockDb.Given(storage.DatabaseIF.FetchAllTopics).When(mock.Anything, publisherID).Return(&[]string{}, expectedErr)

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)
	_, err := topic.GetTopics(context.Background(), publisherID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestGetTopics_Pass(t *testing.T) {
	publisherID := 5000
	topics := &[]string{"golang", "java"}

	mockDb := &test.MockDatabaseIF{}
	mockQueue := &test.MockQueueIF{}
	mockDb.Given(storage.DatabaseIF.FetchAllTopics).When(mock.Anything, publisherID).Return(topics, nil)

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)
	resp, err := topic.GetTopics(context.Background(), publisherID)
	if err != nil {
		t.Fatalf("\nexpected: nil \n\t got: %v", err)
	}

	if !reflect.DeepEqual(resp, topics) {
		t.Fatalf("\nexpected: %v \n\t got: %v", topics, resp)
	}
}

func TestRegisterPublisherToTopic_GetTopicIDFromPublisher_Fail(t *testing.T) {
	publisherID := 5000
	topicName := "golang"
	expectedErr := errors.New("failed to get topicId")

	mockDb := &test.MockDatabaseIF{}
	mockQueue := &test.MockQueueIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromPublisher).When(mock.Anything, publisherID).Return("", false, expectedErr)

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.RegisterPublisherToTopic(context.Background(), publisherID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestRegisterPublisherToTopic_GetTopicIDFromPublisher_AlreadyRegisteredFail(t *testing.T) {
	publisherID := 5000
	topicName := "golang"
	topicID := "12345"
	expectedErr := errors.New("cannot register more than one topic at a time")

	mockDb := &test.MockDatabaseIF{}
	mockQueue := &test.MockQueueIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromPublisher).When(mock.Anything, publisherID).Return(topicID, false, nil)

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.RegisterPublisherToTopic(context.Background(), publisherID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestRegisterPublisherToTopic_GetTopicIDFromTopic_Fail(t *testing.T) {
	publisherID := 5000
	topicName := "golang"
	expectedErr := errors.New("failed to get topicId")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromPublisher).When(mock.Anything, publisherID).Return("", false, nil)
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return("", expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.RegisterPublisherToTopic(context.Background(), publisherID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestRegisterPublisherToTopic_GetTopicIDFromTopic_TopicIdNotFoundFail(t *testing.T) {
	publisherID := 5000
	topicName := "golang"
	topicID := ""
	expectedErr := errors.New("topic not found")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromPublisher).When(mock.Anything, publisherID).Return("", false, nil)
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return(topicID, nil)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.RegisterPublisherToTopic(context.Background(), publisherID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestRegisterPublisherToTopic_InsertPublisher_Fail(t *testing.T) {
	publisherID := 5000
	notFound := true
	topicName := "golang"
	topicID := "12345"
	expectedErr := errors.New("failed to insert publisher")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromPublisher).When(mock.Anything, publisherID).Return("", notFound, nil)
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return(topicID, nil)
	mockDb.Given(storage.DatabaseIF.InsertPublisher).When(mock.Anything, publisherID, topicID).Return(expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.RegisterPublisherToTopic(context.Background(), publisherID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestRegisterPublisherToTopic_UpdateTopicIDIntoPublisher_Fail(t *testing.T) {
	publisherID := 5000
	notFound := false
	topicName := "golang"
	topicID := "12345"
	expectedErr := errors.New("failed to update publisher")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromPublisher).When(mock.Anything, publisherID).Return("", notFound, nil)
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return(topicID, nil)
	mockDb.Given(storage.DatabaseIF.UpdateTopicIDIntoPublisher).When(mock.Anything, publisherID, topicID).Return(expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.RegisterPublisherToTopic(context.Background(), publisherID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestRegisterPublisherToTopic_Pass(t *testing.T) {
	publisherID := 5000
	notFound := false
	topicName := "golang"
	topicID := "12345"

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromPublisher).When(mock.Anything, publisherID).Return("", notFound, nil)
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return(topicID, nil)
	mockDb.Given(storage.DatabaseIF.UpdateTopicIDIntoPublisher).When(mock.Anything, publisherID, topicID).Return(nil)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.RegisterPublisherToTopic(context.Background(), publisherID, topicName)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

func TestDeregisterPublisherFromTopic_GetTopicIDFromPublisher_Fail(t *testing.T) {
	publisherID := 5000
	notFound := false

	expectedErr := errors.New("failed to update publisher")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromPublisher).When(mock.Anything, publisherID).Return("", notFound, expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.DeregisterPublisherFromTopic(context.Background(), publisherID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestDeregisterPublisherFromTopic_GetTopicIDFromPublisher_NotFoundFail(t *testing.T) {
	publisherID := 5000
	notFound := true

	expectedErr := errors.New("you are not registered with any topic")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromPublisher).When(mock.Anything, publisherID).Return("", notFound, nil)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.DeregisterPublisherFromTopic(context.Background(), publisherID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestDeregisterPublisherFromTopic_RemoveTopicIDFromPublisher_Fail(t *testing.T) {
	publisherID := 5000
	notFound := false

	expectedErr := errors.New("failed to remove topicID")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromPublisher).When(mock.Anything, publisherID).Return("", notFound, nil)
	mockDb.Given(storage.DatabaseIF.RemoveTopicIDFromPublisher).When(mock.Anything, publisherID).Return(expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.DeregisterPublisherFromTopic(context.Background(), publisherID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestDeregisterPublisherFromTopic_Pass(t *testing.T) {
	publisherID := 5000
	notFound := false

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromPublisher).When(mock.Anything, publisherID).Return("", notFound, nil)
	mockDb.Given(storage.DatabaseIF.RemoveTopicIDFromPublisher).When(mock.Anything, publisherID).Return(nil)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.DeregisterPublisherFromTopic(context.Background(), publisherID)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

func TestAddMessageToTopic_GetTopicIDFromPublisher_Fail(t *testing.T) {
	publisherID := 5000
	notFound := false
	msg := domain.Message{}

	expectedErr := errors.New("failed to get topicID")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromPublisher).When(mock.Anything, publisherID).Return("", notFound, expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.AddMessageToTopic(context.Background(), publisherID, msg)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestRegisterSubscriberToTopic_GetSubscribedTopics_Fail(t *testing.T) {
	subscriberID := 5000
	topicName := "test"

	expectedErr := errors.New("failed to get subscribed topics")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetSubscribedTopics).When(mock.Anything, subscriberID).Return([]string{}, expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.RegisterSubscriberToTopic(context.Background(), subscriberID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestRegisterSubscriberToTopic_GetTopicIDFromTopic_Fail(t *testing.T) {
	subscriberID := 5000
	topicName := "test"

	expectedErr := errors.New("failed to get topicId")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetSubscribedTopics).When(mock.Anything, subscriberID).Return([]string{}, nil)
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return("", expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.RegisterSubscriberToTopic(context.Background(), subscriberID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestRegisterSubscriberToTopic_InsertSubscriberIDIntoSubscriber_Fail(t *testing.T) {
	subscriberID := 5000
	topicName := "test"
	topicID := "12345"

	expectedErr := errors.New("failed to insert")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetSubscribedTopics).When(mock.Anything, subscriberID).Return([]string{}, nil)
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return(topicID, nil)
	mockDb.Given(storage.DatabaseIF.InsertSubscriberIDIntoSubscriber).When(mock.Anything, subscriberID).Return(expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.RegisterSubscriberToTopic(context.Background(), subscriberID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestRegisterSubscriberToTopic_InsertIntoSubscriberTopicMap_Fail(t *testing.T) {
	subscriberID := 5000
	topicName := "test"
	topicID := "12345"

	expectedErr := errors.New("failed to insert")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetSubscribedTopics).When(mock.Anything, subscriberID).Return([]string{}, nil)
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return(topicID, nil)
	mockDb.Given(storage.DatabaseIF.InsertSubscriberIDIntoSubscriber).When(mock.Anything, subscriberID).Return(nil)
	mockDb.Given(storage.DatabaseIF.InsertIntoSubscriberTopicMap).When(mock.Anything, subscriberID, topicID).Return(expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.RegisterSubscriberToTopic(context.Background(), subscriberID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestRegisterSubscriberToTopic_Pass(t *testing.T) {
	subscriberID := 5000
	topicName := "test"
	topicID := "12345"

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetSubscribedTopics).When(mock.Anything, subscriberID).Return([]string{}, nil)
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return(topicID, nil)
	mockDb.Given(storage.DatabaseIF.InsertSubscriberIDIntoSubscriber).When(mock.Anything, subscriberID).Return(nil)
	mockDb.Given(storage.DatabaseIF.InsertIntoSubscriberTopicMap).When(mock.Anything, subscriberID, topicID).Return(nil)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.RegisterSubscriberToTopic(context.Background(), subscriberID, topicName)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

func TestDeregisterSubscriberFromTopic_GetTopicIDFromTopic_Fail(t *testing.T) {
	subscriberID := 5000
	topicName := "test"

	expectedErr := errors.New("failed to get topicId")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return("", expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.DeregisterSubscriberFromTopic(context.Background(), subscriberID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestDeregisterSubscriberFromTopic_GetTopicIDFromTopic_NotFoundFail(t *testing.T) {
	subscriberID := 5000
	topicName := "test"

	expectedErr := errors.New("topic not found")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return("", nil)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.DeregisterSubscriberFromTopic(context.Background(), subscriberID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestDeregisterSubscriberFromTopic_RemoveTopicIDFromSubscriberTopicMap_Fail(t *testing.T) {
	subscriberID := 5000
	topicName := "test"
	topicID := "12345"

	expectedErr := errors.New("failed to remove subscriber")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return(topicID, nil)
	mockDb.Given(storage.DatabaseIF.RemoveTopicIDFromSubscriberTopicMap).When(mock.Anything, subscriberID, topicID).Return(expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.DeregisterSubscriberFromTopic(context.Background(), subscriberID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestDeregisterSubscriberFromTopic_Pass(t *testing.T) {
	subscriberID := 5000
	topicName := "test"
	topicID := "12345"

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return(topicID, nil)
	mockDb.Given(storage.DatabaseIF.RemoveTopicIDFromSubscriberTopicMap).When(mock.Anything, subscriberID, topicID).Return(nil)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	err := topic.DeregisterSubscriberFromTopic(context.Background(), subscriberID, topicName)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

func TestGetRegisteredTopic_GetSubscribedTopics_Fail(t *testing.T) {
	subscriberID := 5000

	expectedErr := errors.New("failed to get subscribed topics")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetSubscribedTopics).When(mock.Anything, subscriberID).Return([]string{}, expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	_, err := topic.GetRegisteredTopic(context.Background(), subscriberID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestGetRegisteredTopic_GetSubscribedTopics_Pass(t *testing.T) {
	subscriberID := 5000

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetSubscribedTopics).When(mock.Anything, subscriberID).Return([]string{}, nil)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	_, err := topic.GetRegisteredTopic(context.Background(), subscriberID)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

func TestGetMessage(t *testing.T) {
	subscriberID := 5000
	topicName := "test"

	expectedErr := errors.New("failed to get topicId")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.GetTopicIDFromTopic).When(mock.Anything, topicName).Return("", expectedErr)

	mockQueue := &test.MockQueueIF{}

	topic := domain.NewTopic(&logrus.Logger{}, mockDb, mockQueue)

	_, err := topic.GetMessage(context.Background(), subscriberID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

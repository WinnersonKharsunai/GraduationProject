package queue_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/queue"
	"github.com/WinnersonKharsunai/GraduationProject/server/internal/storage"
	"github.com/WinnersonKharsunai/GraduationProject/server/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestSendMessage_EmptyTopicId_Fail(t *testing.T) {
	queueData := getQueue()
	msg := queue.SendMessageRequest{
		TopicID: "",
		Message: queue.Message{
			MessageID: "message1",
			Data:      "test data 1",
			CretedAt:  "2021-02-27 20:03:09",
			ExpiresAt: "2021-02-27 20:04:09",
		},
	}

<<<<<<< HEAD
	expectedErr := errors.New("you are not register to any topics")
=======
	expectedErr := errors.New("topicId cannot be empty")
>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.FetchQueues).When(mock.Anything).Return(&queueData, nil)
	mockDb.Given(storage.DatabaseIF.RemoveMessagesFromQueue).When(mock.Anything).Return(nil)

	q, err := queue.NewQueue(&logrus.Logger{}, mockDb)
	if err != nil {
		t.Fatalf("\nexpected: nil \n\t got: %v", err)
	}

	err = q.SendMessage(context.Background(), msg)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("\nexpected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestSendMessage_EmptyMessage_Fail(t *testing.T) {
	queueData := getQueue()
	msg := queue.SendMessageRequest{
		TopicID: "12345",
	}

	expectedErr := errors.New("message cannot be empty")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.FetchQueues).When(mock.Anything).Return(&queueData, nil)
	mockDb.Given(storage.DatabaseIF.RemoveMessagesFromQueue).When(mock.Anything).Return(nil)

	q, err := queue.NewQueue(&logrus.Logger{}, mockDb)
	if err != nil {
		t.Fatalf("\nexpected: nil \n\t got: %v", err)
	}

	err = q.SendMessage(context.Background(), msg)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("\nexpected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestSendMessage_Pass(t *testing.T) {
	queueData := getQueue()
	msg := queue.SendMessageRequest{
		TopicID: "12345",
		Message: queue.Message{
			MessageID: "message1",
			Data:      "test data 1",
			CretedAt:  "2021-02-27 20:03:09",
			ExpiresAt: "2021-02-27 20:04:09",
		},
	}

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.FetchQueues).When(mock.Anything).Return(&queueData, nil)
	mockDb.Given(storage.DatabaseIF.RemoveMessagesFromQueue).When(mock.Anything).Return(nil)

	q, err := queue.NewQueue(&logrus.Logger{}, mockDb)
	if err != nil {
		t.Fatalf("\nexpected: nil \n\t got: %v", err)
	}

	err = q.SendMessage(context.Background(), msg)
	if err != nil {
		t.Fatalf("\nexpected: nil \n\t got: %v", err)
	}
}

func TestRetrieveMessage_Fail(t *testing.T) {
	queueData := getQueue()
	topicID := "golang"

	expectedErr := errors.New("no message present in queue")

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.FetchQueues).When(mock.Anything).Return(&queueData, nil)
	mockDb.Given(storage.DatabaseIF.RemoveMessagesFromQueue).When(mock.Anything).Return(nil)

	q, err := queue.NewQueue(&logrus.Logger{}, mockDb)
	if err != nil {
		t.Fatalf("\nexpected: nil \n\t got: %v", err)
	}

	_, err = q.RetrieveMessage(context.Background(), topicID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("\nexpected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestRetrieveMessage_Pass(t *testing.T) {
	queueData := getQueue()
	topicID := "golang123"
	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.FetchQueues).When(mock.Anything).Return(&queueData, nil)
	mockDb.Given(storage.DatabaseIF.RemoveMessagesFromQueue).When(mock.Anything).Return(nil)

	q, err := queue.NewQueue(&logrus.Logger{}, mockDb)
	if err != nil {
		t.Fatalf("\nexpected: nil \n\t got: %v", err)
	}

	_, err = q.RetrieveMessage(context.Background(), topicID)
	if err != nil {
		t.Fatalf("\nexpected: nil \n\t got: %v", err)
	}
}

func TestBackUpQueue_Pass(t *testing.T) {
	queueData := getQueue()

	mockDb := &test.MockDatabaseIF{}
	mockDb.Given(storage.DatabaseIF.FetchQueues).When(mock.Anything).Return(&queueData, nil)
	mockDb.Given(storage.DatabaseIF.RemoveMessagesFromQueue).When(mock.Anything).Return(nil)
	mockDb.Given(storage.DatabaseIF.SaveQueues).When(mock.Anything, mock.Anything, true).Return(nil)
	mockDb.Given(storage.DatabaseIF.SaveQueues).When(mock.Anything, mock.Anything, false).Return(nil)

	q, err := queue.NewQueue(&logrus.Logger{}, mockDb)
	if err != nil {
		t.Fatalf("\nexpected: nil \n\t got: %v", err)
	}

	err = q.BackUpQueue(context.Background())
	if err != nil {
		t.Fatalf("\nexpected: nil \n\t got: %v", err)
	}
}

func getQueue() storage.Queue {
	now := time.Now().UTC()
	return storage.Queue{
		Topic: map[string][]storage.Message{
			"golang123": {
				{
					MessageID: "message1",
					Data:      "test data 1",
					CretedAt:  now.Format("2006-01-02 15:04:05"),
					ExpiresAt: now.Add(time.Duration(time.Second * 60)).Format("2006-01-02 15:04:05"),
				},
				{
					MessageID: "message1",
					Data:      "test data 1",
					CretedAt:  "2021-02-27 20:03:09",
					ExpiresAt: "2021-02-27 20:04:09",
				},
			},
			"java123": {
				{
					MessageID: "message2",
					Data:      "test data 2",
					CretedAt:  now.Format("2006-01-02 15:04:05"),
					ExpiresAt: now.Add(time.Duration(time.Second * 60)).Format("2006-01-02 15:04:05"),
				},
			},
		},
	}
}

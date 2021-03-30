package publisher_test

import (
	"context"
	"errors"
	"testing"

	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/publisher"
	"github.com/WinnersonKharsunai/GraduationProject/server/internal/domain"
	"github.com/WinnersonKharsunai/GraduationProject/server/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestShowTopics_Fail(t *testing.T) {
	req := &publisher.ShowTopicRequest{
		PublisherID: 5000,
	}

	expectedErr := errors.New("failed to get topics")

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.GetTopics).When(mock.Anything, req.PublisherID).Return(&[]string{}, expectedErr)

	pub := publisher.NewPublisher(&logrus.Logger{}, mockTopicSvc)
	_, err := pub.ShowTopics(context.Background(), req)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestShowTopics_Pass(t *testing.T) {
	req := &publisher.ShowTopicRequest{
		PublisherID: 5000,
	}

	resp := &[]string{"golang"}

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.GetTopics).When(mock.Anything, req.PublisherID).Return(resp, nil)

	pub := publisher.NewPublisher(&logrus.Logger{}, mockTopicSvc)
	_, err := pub.ShowTopics(context.Background(), req)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

func TestConnectToTopic_Fail(t *testing.T) {
	req := &publisher.ConnectToTopicRequest{
		PublisherID: 5000,
		TopicName:   "golang",
	}

	expectedErr := errors.New("failed to get register to topic")

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.RegisterPublisherToTopic).When(mock.Anything, req.PublisherID).Return(expectedErr)

	pub := publisher.NewPublisher(&logrus.Logger{}, mockTopicSvc)
	_, err := pub.ConnectToTopic(context.Background(), req)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestConnectToTopic_Pass(t *testing.T) {
	req := &publisher.ConnectToTopicRequest{
		PublisherID: 5000,
		TopicName:   "golang",
	}

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.RegisterPublisherToTopic).When(mock.Anything, req.PublisherID).Return(nil)

	pub := publisher.NewPublisher(&logrus.Logger{}, mockTopicSvc)
	_, err := pub.ConnectToTopic(context.Background(), req)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

func TestDisconnectFromTopic_Fail(t *testing.T) {
	req := &publisher.DisconnectFromTopicRequest{
		PublisherID: 500,
	}

	expectedErr := errors.New("failed to get deregister from topic")

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.DeregisterPublisherFromTopic).When(mock.Anything, req.PublisherID).Return(expectedErr)

	pub := publisher.NewPublisher(&logrus.Logger{}, mockTopicSvc)
	_, err := pub.DisconnectFromTopic(context.Background(), req)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestDisconnectFromTopic_Pass(t *testing.T) {
	req := &publisher.DisconnectFromTopicRequest{
		PublisherID: 500,
	}

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.DeregisterPublisherFromTopic).When(mock.Anything, req.PublisherID).Return(nil)

	pub := publisher.NewPublisher(&logrus.Logger{}, mockTopicSvc)
	_, err := pub.DisconnectFromTopic(context.Background(), req)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

func TestPublishMessage_Fail(t *testing.T) {
	req := &publisher.PublishMessageRequest{
		PublisherID: 5000,
		Message: publisher.Message{
			Data:      "test data",
			CretedAt:  "2021-02-27 20:03:09",
			ExpiresAt: "2021-02-27 20:04:09",
		},
	}

	expectedErr := errors.New("failed to add message to topic")

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.AddMessageToTopic).When(mock.Anything, req.PublisherID, mock.Anything).Return(expectedErr)

	pub := publisher.NewPublisher(&logrus.Logger{}, mockTopicSvc)
	_, err := pub.PublishMessage(context.Background(), req)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestPublishMessage_Pass(t *testing.T) {
	req := &publisher.PublishMessageRequest{
		PublisherID: 5000,
		Message: publisher.Message{
			Data:      "test data",
			CretedAt:  "2021-02-27 20:03:09",
			ExpiresAt: "2021-02-27 20:04:09",
		},
	}

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.AddMessageToTopic).When(mock.Anything, req.PublisherID, mock.Anything).Return(nil)

	pub := publisher.NewPublisher(&logrus.Logger{}, mockTopicSvc)
	_, err := pub.PublishMessage(context.Background(), req)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

package subscriber_test

import (
	"context"
	"errors"
	"testing"

	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/subscriber"
	"github.com/WinnersonKharsunai/GraduationProject/server/internal/domain"
	"github.com/WinnersonKharsunai/GraduationProject/server/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestShowTopics_Fail(t *testing.T) {
	req := &subscriber.ShowTopicRequest{
		SubscriberID: 6000,
	}

	expectedErr := errors.New("failed to get topics")

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.GetTopics).When(mock.Anything, req.SubscriberID).Return(&[]string{}, expectedErr)

	sub := subscriber.NewSubscriber(&logrus.Logger{}, mockTopicSvc)
	_, err := sub.ShowTopics(context.Background(), req)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestShowTopics_Pass(t *testing.T) {
	req := &subscriber.ShowTopicRequest{
		SubscriberID: 6000,
	}

	resp := &[]string{"golang"}

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.GetTopics).When(mock.Anything, req.SubscriberID).Return(resp, nil)

	sub := subscriber.NewSubscriber(&logrus.Logger{}, mockTopicSvc)
	_, err := sub.ShowTopics(context.Background(), req)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

func TestSubscribeToTopic_Fail(t *testing.T) {
	req := &subscriber.SubscribeToTopicRequest{
		SubscriberID: 6000,
		TopicName:    "golang",
	}

	expectedErr := errors.New("failed to register to topic")

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.RegisterSubscriberToTopic).When(mock.Anything, req.SubscriberID, req.TopicName).Return(expectedErr)

	sub := subscriber.NewSubscriber(&logrus.Logger{}, mockTopicSvc)
	_, err := sub.SubscribeToTopic(context.Background(), req)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestSubscribeToTopic_Pass(t *testing.T) {
	req := &subscriber.SubscribeToTopicRequest{
		SubscriberID: 6000,
		TopicName:    "golang",
	}

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.RegisterSubscriberToTopic).When(mock.Anything, req.SubscriberID, req.TopicName).Return(nil)

	sub := subscriber.NewSubscriber(&logrus.Logger{}, mockTopicSvc)
	_, err := sub.SubscribeToTopic(context.Background(), req)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

func TestUnsubscribeFromTopic_Fail(t *testing.T) {
	req := &subscriber.UnsubscribeFromTopicRequest{
		SubscriberID: 6000,
		TopicName:    "golang",
	}

	expectedErr := errors.New("failed to deregister to topic")

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.DeregisterSubscriberFromTopic).When(mock.Anything, req.SubscriberID, req.TopicName).Return(expectedErr)

	sub := subscriber.NewSubscriber(&logrus.Logger{}, mockTopicSvc)
	_, err := sub.UnsubscribeFromTopic(context.Background(), req)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestUnsubscribeFromTopic_Pass(t *testing.T) {
	req := &subscriber.UnsubscribeFromTopicRequest{
		SubscriberID: 6000,
		TopicName:    "golang",
	}

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.DeregisterSubscriberFromTopic).When(mock.Anything, req.SubscriberID, req.TopicName).Return(nil)

	sub := subscriber.NewSubscriber(&logrus.Logger{}, mockTopicSvc)
	_, err := sub.UnsubscribeFromTopic(context.Background(), req)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

func TestGetSubscribedTopics_Fail(t *testing.T) {
	req := &subscriber.GetSubscribedTopicsRequest{
		SubscriberID: 6000,
	}
	resp := &[]string{"golang"}

	expectedErr := errors.New("failed to get subscribed topics")

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.GetRegisteredTopic).When(mock.Anything, req.SubscriberID).Return(resp, expectedErr)

	sub := subscriber.NewSubscriber(&logrus.Logger{}, mockTopicSvc)
	_, err := sub.GetSubscribedTopics(context.Background(), req)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestGetSubscribedTopics_Pass(t *testing.T) {
	req := &subscriber.GetSubscribedTopicsRequest{
		SubscriberID: 6000,
	}
	resp := &[]string{"golang"}

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.GetRegisteredTopic).When(mock.Anything, req.SubscriberID).Return(resp, nil)

	sub := subscriber.NewSubscriber(&logrus.Logger{}, mockTopicSvc)
	_, err := sub.GetSubscribedTopics(context.Background(), req)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

func TestGetMessageFromTopic_Fail(t *testing.T) {
	req := &subscriber.GetMessageFromTopicRequest{
		SubscriberID: 6000,
		TopicName:    "golang",
	}

	expectedErr := errors.New("failed to get message from topic")

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.GetMessage).When(mock.Anything, req.SubscriberID, req.TopicName).Return(&domain.Message{}, expectedErr)

	sub := subscriber.NewSubscriber(&logrus.Logger{}, mockTopicSvc)
	_, err := sub.GetMessageFromTopic(context.Background(), req)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v \n\t got: %v", expectedErr, err)
	}
}

func TestGetMessageFromTopic_Pass(t *testing.T) {
	req := &subscriber.GetMessageFromTopicRequest{
		SubscriberID: 6000,
		TopicName:    "golang",
	}

	resp := &domain.Message{
		MessageID: "123",
		Data:      "test data",
		CretedAt:  "2021-02-27 20:03:09",
		ExpiresAt: "2021-02-27 20:04:09",
	}

	mockTopicSvc := &test.MockTopicServiceIF{}
	mockTopicSvc.Given(domain.TopicServicesIF.GetMessage).When(mock.Anything, req.SubscriberID, req.TopicName).Return(resp, nil)

	sub := subscriber.NewSubscriber(&logrus.Logger{}, mockTopicSvc)
	_, err := sub.GetMessageFromTopic(context.Background(), req)
	if err != nil {
		t.Fatalf("expected: nil \n\t got: %v", err)
	}
}

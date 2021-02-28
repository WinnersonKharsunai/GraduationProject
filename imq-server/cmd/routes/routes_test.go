package routes_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/routes"
	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/publisher"
	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/subscriber"
	"github.com/WinnersonKharsunai/GraduationProject/server/pkg/protocol"
	"github.com/WinnersonKharsunai/GraduationProject/server/test"
	"github.com/stretchr/testify/mock"
)

var hdr protocol.Header

func init() {
	hdr = protocol.Header{
		Version:     "1.0",
		ContentType: "json",
	}
}

func TestRequestRouter_ShowTopic(t *testing.T) {
	showTopicRequest := &publisher.ShowTopicRequest{
		PublisherID: 5000,
	}

	showTopicResponse := &publisher.ShowTopicResponse{
		Topics: []string{"golang"},
	}

	hdr.Method = "showTopicRequest"

	request := getRequestString(hdr, showTopicRequest)

	mockPsvc := &test.MockPublisherIF{}
	mockPsvc.Given(publisher.PublisherIF.ShowTopics).When(mock.Anything, showTopicRequest).Return(showTopicResponse, nil)
	mockSsvc := &test.MockSubscriberIF{}

	route := routes.NewHandler(mockPsvc, mockSsvc)

	resp := route.RequestRouter(context.Background(), request)
	if resp.Error != "" {
		t.Fatalf("\necpected: %v \n\t got: %v", nil, resp.Error)
	}
}

func TestRequestRouter_ConnectToTopicRequest(t *testing.T) {
	connectToTopicRequest := &publisher.ConnectToTopicRequest{
		PublisherID: 5000,
	}

	connectToTopicResponse := &publisher.ConnectToTopicResponse{
		Status: "connected",
	}

	hdr.Method = "connectToTopicRequest"

	request := getRequestString(hdr, connectToTopicRequest)

	mockPsvc := &test.MockPublisherIF{}
	mockPsvc.Given(publisher.PublisherIF.ConnectToTopic).When(mock.Anything, connectToTopicRequest).Return(connectToTopicResponse, nil)
	mockSsvc := &test.MockSubscriberIF{}

	route := routes.NewHandler(mockPsvc, mockSsvc)

	resp := route.RequestRouter(context.Background(), request)
	if resp.Error != "" {
		t.Fatalf("\necpected: %v \n\t got: %v", nil, resp.Error)
	}
}

func TestRequestRouter_DisconnectFromTopicRequest(t *testing.T) {
	disconnectFromTopicRequest := &publisher.DisconnectFromTopicRequest{
		PublisherID: 5000,
	}

	disconnectFromTopicResponse := &publisher.DisconnectFromTopicResponse{
		Status: "disconnected",
	}

	hdr.Method = "disconnectFromTopicRequest"

	request := getRequestString(hdr, disconnectFromTopicRequest)

	mockPsvc := &test.MockPublisherIF{}
	mockPsvc.Given(publisher.PublisherIF.DisconnectFromTopic).When(mock.Anything, disconnectFromTopicRequest).Return(disconnectFromTopicResponse, nil)
	mockSsvc := &test.MockSubscriberIF{}

	route := routes.NewHandler(mockPsvc, mockSsvc)

	resp := route.RequestRouter(context.Background(), request)
	if resp.Error != "" {
		t.Fatalf("\necpected: %v \n\t got: %v", nil, resp.Error)
	}
}
func TestRequestRouter_PublishMessageRequest(t *testing.T) {
	publishMessageRequest := &publisher.PublishMessageRequest{
		PublisherID: 600,
		Message: publisher.Message{
			Data:      "test data",
			CretedAt:  "2021-02-27 20:03:09",
			ExpiresAt: "2021-02-27 20:04:09",
		},
	}

	publishMessageResponse := &publisher.PublishMessageResponse{
		Status: "successful",
	}

	hdr.Method = "publishMessageRequest"

	request := getRequestString(hdr, publishMessageRequest)

	mockPsvc := &test.MockPublisherIF{}
	mockPsvc.Given(publisher.PublisherIF.PublishMessage).When(mock.Anything, publishMessageRequest).Return(publishMessageResponse, nil)
	mockSsvc := &test.MockSubscriberIF{}

	route := routes.NewHandler(mockPsvc, mockSsvc)

	resp := route.RequestRouter(context.Background(), request)
	if resp.Error != "" {
		t.Fatalf("\necpected: %v \n\t got: %v", nil, resp.Error)
	}
}
func TestRequestRouter_SubscribeToTopicRequest(t *testing.T) {
	subscribeToTopicRequest := &subscriber.SubscribeToTopicRequest{
		SubscriberID: 6000,
		TopicName:    "java",
	}

	subscribeToTopicResponse := &subscriber.SubscribeToTopicResponse{
		Status: "subscribed",
	}

	hdr.Method = "subscribeToTopicRequest"

	request := getRequestString(hdr, subscribeToTopicRequest)

	mockPsvc := &test.MockPublisherIF{}

	mockSsvc := &test.MockSubscriberIF{}
	mockSsvc.Given(subscriber.SubscriberIF.SubscribeToTopic).When(mock.Anything, subscribeToTopicRequest).Return(subscribeToTopicResponse, nil)

	route := routes.NewHandler(mockPsvc, mockSsvc)

	resp := route.RequestRouter(context.Background(), request)
	if resp.Error != "" {
		t.Fatalf("\necpected: %v \n\t got: %v", nil, resp.Error)
	}
}
func TestRequestRouter_UnsubscribeFromTopicRequest(t *testing.T) {
	unsubscribeFromTopicRequest := &subscriber.UnsubscribeFromTopicRequest{
		SubscriberID: 6000,
		TopicName:    "java",
	}

	unsubscribeFromTopicResponse := &subscriber.UnsubscribeFromTopicResponse{
		Status: "unsubscribed",
	}

	hdr.Method = "unsubscribeFromTopicRequest"

	request := getRequestString(hdr, unsubscribeFromTopicRequest)

	mockPsvc := &test.MockPublisherIF{}
	mockSsvc := &test.MockSubscriberIF{}
	mockSsvc.Given(subscriber.SubscriberIF.UnsubscribeFromTopic).When(mock.Anything, unsubscribeFromTopicRequest).Return(unsubscribeFromTopicResponse, nil)

	route := routes.NewHandler(mockPsvc, mockSsvc)

	resp := route.RequestRouter(context.Background(), request)
	if resp.Error != "" {
		t.Fatalf("\necpected: %v \n\t got: %v", nil, resp.Error)
	}
}
func TestRequestRouter_GetSubscribedTopicsRequest(t *testing.T) {
	getSubscribedTopicsRequest := &subscriber.GetSubscribedTopicsRequest{
		SubscriberID: 6000,
	}

	getSubscribedTopicsResponse := &subscriber.GetSubscribedTopicsResponse{
		Topics: []string{"java"},
	}

	hdr.Method = "getSubscribedTopicsRequest"

	request := getRequestString(hdr, getSubscribedTopicsRequest)

	mockPsvc := &test.MockPublisherIF{}
	mockSsvc := &test.MockSubscriberIF{}
	mockSsvc.Given(subscriber.SubscriberIF.GetSubscribedTopics).When(mock.Anything, getSubscribedTopicsRequest).Return(getSubscribedTopicsResponse, nil)

	route := routes.NewHandler(mockPsvc, mockSsvc)

	resp := route.RequestRouter(context.Background(), request)
	if resp.Error != "" {
		t.Fatalf("\necpected: %v \n\t got: %v", nil, resp.Error)
	}
}
func TestRequestRouter_GetMessageFromTopicRequest(t *testing.T) {
	getMessageFromTopicRequest := &subscriber.GetMessageFromTopicRequest{
		SubscriberID: 6000,
		TopicName:    "java",
	}

	getMessageFromTopicResponse := &subscriber.GetMessageFromTopicResponse{
		Message: subscriber.Message{
			Data: "test data",
		},
	}

	hdr.Method = "getMessageFromTopicRequest"

	request := getRequestString(hdr, getMessageFromTopicRequest)

	mockPsvc := &test.MockPublisherIF{}
	mockSsvc := &test.MockSubscriberIF{}
	mockSsvc.Given(subscriber.SubscriberIF.GetMessageFromTopic).When(mock.Anything, getMessageFromTopicRequest).Return(getMessageFromTopicResponse, nil)

	route := routes.NewHandler(mockPsvc, mockSsvc)

	resp := route.RequestRouter(context.Background(), request)
	if resp.Error != "" {
		t.Fatalf("\necpected: %v \n\t got: %v", nil, resp.Error)
	}
}

func getRequestString(hdr protocol.Header, v interface{}) string {
	body, _ := json.Marshal(v)
	request := protocol.Request{
		Header: hdr,
		Body:   string(body),
	}
	requestBytes, _ := json.Marshal(request)
	return string(requestBytes)
}

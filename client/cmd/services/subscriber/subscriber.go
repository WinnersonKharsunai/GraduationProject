package subscriber

import (
	"context"

	messagefactory "github.com/WinnersonKharsunai/GraduationProject/client/message-factory"
	"github.com/WinnersonKharsunai/GraduationProject/client/pkg/client"
	"github.com/WinnersonKharsunai/GraduationProject/client/pkg/protocol"
)

// Subscriber is the concrete implementation for the Subscriber
type Subscriber struct {
	client   *client.Client
	protocol protocol.Header
	factory  messagefactory.MessagefactoryIF
}

// SubscriberService is the interface to be for the Subscriber service
type SubscriberService interface {
	ShowTopics(ctx context.Context, in *ShowTopicRequest) (*ShowTopicResponse, error)
	SubscribeToTopic(ctx context.Context, in *SubscribeToTopicRequest) (*SubscribeToTopicResponse, error)
	UnsubscribeFromTopic(ctx context.Context, in *UnsubscribeFromTopicRequest) (*UnsubscribeFromTopicResponse, error)
	GetSubscribedTopics(ctx context.Context, in *GetSubscribedTopicsRequest) (*GetSubscribedTopicsResponse, error)
	GetMessageFromTopic(ctx context.Context, in *GetMessageFromTopicRequest) (*GetMessageFromTopicResponse, error)
}

// NewSubscriber is the factory function for the Subscriber
func NewSubscriber(client *client.Client, factory messagefactory.MessagefactoryIF) SubscriberService {
	return &Subscriber{
		client:  client,
		factory: factory,
	}
}

// ShowTopics fetch all the topics that are available
func (s *Subscriber) ShowTopics(ctx context.Context, in *ShowTopicRequest) (*ShowTopicResponse, error) {

	showTopicResponse := &ShowTopicResponse{}

	hdr := protocol.SetHeader(version, contentType, showTopic, s.client.Addr)

	bodyBytes, err := s.factory.MarshalRequestBody(in, contentType)
	if err != nil {
		return nil, err
	}

	request := protocol.Request{
		Header: hdr,
		Body:   string(bodyBytes),
	}

	responseBytes, err := s.client.SendRequest(ctx, &request)
	if err != nil {
		return nil, err
	}

	err = s.factory.UnmarshalRequestBody(responseBytes, &showTopicResponse, contentType)
	if err != nil {
		return nil, err
	}

	return showTopicResponse, nil
}

// GetSubscribedTopics fetches all the topics subscribed by given client
func (s *Subscriber) GetSubscribedTopics(ctx context.Context, in *GetSubscribedTopicsRequest) (*GetSubscribedTopicsResponse, error) {

	var getSubscribedTopicsResponse *GetSubscribedTopicsResponse

	hdr := protocol.SetHeader(version, contentType, getSubscribedTopics, s.client.Addr)

	bodyBytes, err := s.factory.MarshalRequestBody(in, contentType)
	if err != nil {
		return nil, err
	}

	request := protocol.Request{
		Header: hdr,
		Body:   string(bodyBytes),
	}

	responseBytes, err := s.client.SendRequest(ctx, &request)
	if err != nil {
		return nil, err
	}

	err = s.factory.UnmarshalRequestBody(responseBytes, &getSubscribedTopicsResponse, contentType)
	if err != nil {
		return nil, err
	}

	return getSubscribedTopicsResponse, nil
}

// SubscribeToTopic subscribes given subscriber to topic
func (s *Subscriber) SubscribeToTopic(ctx context.Context, in *SubscribeToTopicRequest) (*SubscribeToTopicResponse, error) {

	var subscribeToTopicResponse *SubscribeToTopicResponse

	hdr := protocol.SetHeader(version, contentType, subscribeToTopic, s.client.Addr)

	bodyBytes, err := s.factory.MarshalRequestBody(in, contentType)
	if err != nil {
		return nil, err
	}

	request := protocol.Request{
		Header: hdr,
		Body:   string(bodyBytes),
	}

	responseBytes, err := s.client.SendRequest(ctx, &request)
	if err != nil {
		return nil, err
	}

	err = s.factory.UnmarshalRequestBody(responseBytes, &subscribeToTopicResponse, contentType)
	if err != nil {
		return nil, err
	}

	return subscribeToTopicResponse, nil
}

// UnsubscribeFromTopic ubsubscribes given client from topic
func (s *Subscriber) UnsubscribeFromTopic(ctx context.Context, in *UnsubscribeFromTopicRequest) (*UnsubscribeFromTopicResponse, error) {

	var unsubscribeFromTopicResponse *UnsubscribeFromTopicResponse

	hdr := protocol.SetHeader(version, contentType, unsubscribeFromTopic, s.client.Addr)

	bodyBytes, err := s.factory.MarshalRequestBody(in, contentType)
	if err != nil {
		return nil, err
	}

	request := protocol.Request{
		Header: hdr,
		Body:   string(bodyBytes),
	}

	responseBytes, err := s.client.SendRequest(ctx, &request)
	if err != nil {
		return nil, err
	}

	err = s.factory.UnmarshalRequestBody(responseBytes, &unsubscribeFromTopicResponse, contentType)
	if err != nil {
		return nil, err
	}

	return unsubscribeFromTopicResponse, nil
}

// GetMessageFromTopic fetches message publised to a given topic
func (s *Subscriber) GetMessageFromTopic(ctx context.Context, in *GetMessageFromTopicRequest) (*GetMessageFromTopicResponse, error) {

	var getMessageFromTopicResponse *GetMessageFromTopicResponse

	hdr := protocol.SetHeader(version, contentType, getMessageFromTopic, s.client.Addr)

	bodyBytes, err := s.factory.MarshalRequestBody(in, contentType)
	if err != nil {
		return nil, err
	}

	request := protocol.Request{
		Header: hdr,
		Body:   string(bodyBytes),
	}

	responseBytes, err := s.client.SendRequest(ctx, &request)
	if err != nil {
		return nil, err
	}

	err = s.factory.UnmarshalRequestBody(responseBytes, &getMessageFromTopicResponse, contentType)
	if err != nil {
		return nil, err
	}

	return getMessageFromTopicResponse, nil
}

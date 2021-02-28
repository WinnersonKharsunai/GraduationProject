package publisher

import (
	"context"

	messagefactory "github.com/WinnersonKharsunai/GraduationProject/client/message-factory"
	"github.com/WinnersonKharsunai/GraduationProject/client/pkg/client"
	"github.com/WinnersonKharsunai/GraduationProject/client/pkg/protocol"
)

// Publisher is the concrete implementation for the Publisher
type Publisher struct {
	client   client.Service
	protocol protocol.Header
	factory  messagefactory.MessagefactoryIF
}

// Service is the interface for the Publisher service
type Service interface {
	ShowTopics(ctx context.Context, in *ShowTopicRequest) (*ShowTopicResponse, error)
	ConnectToTopic(ctx context.Context, in *ConnectToTopicRequest) (*ConnectToTopicResponse, error)
	DisconnectFromTopic(ctx context.Context, in *DisconnectFromTopicRequest) (*DisconnectFromTopicResponse, error)
	PublishMessage(ctx context.Context, in *PublishMessageRequest) (*PublishMessageResponse, error)
}

// NewPublisher is the factory function for the Publisher type
func NewPublisher(client client.Service, factory messagefactory.MessagefactoryIF) Service {
	return &Publisher{
		client:  client,
		factory: factory,
	}
}

// ShowTopics fetch all the topics that are available
func (p *Publisher) ShowTopics(ctx context.Context, in *ShowTopicRequest) (*ShowTopicResponse, error) {

	var showTopicResponse ShowTopicResponse

	hdr := protocol.SetHeader(version, contentType, showTopic, p.client.GetAddress())

	bodyBytes, err := p.factory.MarshalRequestBody(in, contentType)
	if err != nil {
		return nil, err
	}

	request := protocol.Request{
		Header: hdr,
		Body:   string(bodyBytes),
	}

	responseBytes, err := p.client.SendRequest(ctx, &request)
	if err != nil {
		return nil, err
	}

	err = p.factory.UnmarshalRequestBody(responseBytes, &showTopicResponse, contentType)
	if err != nil {
		return nil, err
	}

	return &showTopicResponse, nil
}

// ConnectToTopic register publisher to topic
func (p *Publisher) ConnectToTopic(ctx context.Context, in *ConnectToTopicRequest) (*ConnectToTopicResponse, error) {

	var connectToTopicResponse *ConnectToTopicResponse

	hdr := protocol.SetHeader(version, contentType, connectToTopic, p.client.GetAddress())

	bodyBytes, err := p.factory.MarshalRequestBody(in, contentType)
	if err != nil {
		return nil, err
	}

	request := protocol.Request{
		Header: hdr,
		Body:   string(bodyBytes),
	}

	responseBytes, err := p.client.SendRequest(ctx, &request)
	if err != nil {
		return nil, err
	}

	err = p.factory.UnmarshalRequestBody(responseBytes, &connectToTopicResponse, contentType)
	if err != nil {
		return nil, err
	}

	return connectToTopicResponse, nil
}

// DisconnectFromTopic deregister publisher from topic
func (p *Publisher) DisconnectFromTopic(ctx context.Context, in *DisconnectFromTopicRequest) (*DisconnectFromTopicResponse, error) {

	var disconnectFromTopicResponse *DisconnectFromTopicResponse

	hdr := protocol.SetHeader(version, contentType, disconnectFromTopic, p.client.GetAddress())

	bodyBytes, err := p.factory.MarshalRequestBody(in, contentType)
	if err != nil {
		return nil, err
	}

	request := protocol.Request{
		Header: hdr,
		Body:   string(bodyBytes),
	}

	responseBytes, err := p.client.SendRequest(ctx, &request)
	if err != nil {
		return nil, err
	}

	err = p.factory.UnmarshalRequestBody(responseBytes, &disconnectFromTopicResponse, contentType)
	if err != nil {
		return nil, err
	}

	return disconnectFromTopicResponse, nil
}

// PublishMessage publishes new message to topic
func (p *Publisher) PublishMessage(ctx context.Context, in *PublishMessageRequest) (*PublishMessageResponse, error) {

	var publishMessageResponse *PublishMessageResponse

	hdr := protocol.SetHeader(version, contentType, publishMessage, p.client.GetAddress())

	bodyBytes, err := p.factory.MarshalRequestBody(in, contentType)
	if err != nil {
		return nil, err
	}

	request := protocol.Request{
		Header: hdr,
		Body:   string(bodyBytes),
	}

	responseBytes, err := p.client.SendRequest(ctx, &request)
	if err != nil {
		return nil, err
	}

	err = p.factory.UnmarshalRequestBody(responseBytes, &publishMessageResponse, contentType)
	if err != nil {
		return nil, err
	}

	return publishMessageResponse, nil
}

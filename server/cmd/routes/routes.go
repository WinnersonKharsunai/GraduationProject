package routes

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"

	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/publisher"
	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/subscriber"
	"github.com/WinnersonKharsunai/GraduationProject/server/pkg/protocol"
)

// Handler is the concrete implementation for Router
type Handler struct {
	pSvc publisher.PublisherIF
	sSvc subscriber.SubscriberIF
}

// Router is the interface for the Handler type
type Router interface {
	RequestRouter(ctx context.Context, r string) *protocol.Response
}

// NewHandler is the factory function for the Handler type
func NewHandler(pSvc publisher.PublisherIF, sSvc subscriber.SubscriberIF) Router {
	return &Handler{
		pSvc: pSvc,
		sSvc: sSvc,
	}
}

// RequestRouter handles all the request and response
func (h Handler) RequestRouter(ctx context.Context, r string) *protocol.Response {

	request := protocol.Request{}
	if err := json.Unmarshal([]byte(r), &request); err != nil {
		return &protocol.Response{Error: err.Error()}
	}

	if err := request.Header.ValidateRequestHeader(ctx); err != nil {
		return &protocol.Response{Error: err.Error()}
	}

	resp, err := processRequest(ctx, h.pSvc, h.sSvc, request)
	if err != nil {
		return &protocol.Response{Error: err.Error()}
	}

	body, err := marshal(resp, request.Header.ContentType)
	if err != nil {
		return &protocol.Response{Error: err.Error()}
	}

	return &protocol.Response{Body: body}
}

func processRequest(ctx context.Context, p publisher.PublisherIF, s subscriber.SubscriberIF, request protocol.Request) (interface{}, error) {
	switch request.Header.Method {
	case showTopic:
		showTopicRequest := &publisher.ShowTopicRequest{}
		if err := unmarshal([]byte(request.Body), showTopicRequest, request.Header.ContentType); err != nil {
			return nil, err
		}
		return p.ShowTopics(ctx, showTopicRequest)

	case connectToTopic:
		connectToTopicRequest := &publisher.ConnectToTopicRequest{}
		if err := unmarshal([]byte(request.Body), connectToTopicRequest, request.Header.ContentType); err != nil {
			return nil, err
		}
		return p.ConnectToTopic(ctx, connectToTopicRequest)

	case disconnectFromTopic:
		disconnectFromTopicRequest := &publisher.DisconnectFromTopicRequest{}
		if err := unmarshal([]byte(request.Body), disconnectFromTopicRequest, request.Header.ContentType); err != nil {
			return nil, err
		}
		return p.DisconnectFromTopic(ctx, disconnectFromTopicRequest)

	case publishMessage:
		publishMessageRequest := &publisher.PublishMessageRequest{}
		if err := unmarshal([]byte(request.Body), publishMessageRequest, request.Header.ContentType); err != nil {
			return nil, err
		}
		return p.PublishMessage(ctx, publishMessageRequest)
	case checkMessageStatus:
		checkMessageStatusRequest := &publisher.CheckMessageStatusRequest{}
		if err := unmarshal([]byte(request.Body), checkMessageStatusRequest, request.Header.ContentType); err != nil {
			return nil, err
		}
		return p.CheckMessageStatus(ctx, checkMessageStatusRequest)

	case subscribeToTopic:
		subscribeToTopicRequest := &subscriber.SubscribeToTopicRequest{}
		if err := unmarshal([]byte(request.Body), subscribeToTopicRequest, request.Header.ContentType); err != nil {
			return nil, err
		}
		return s.SubscribeToTopic(ctx, subscribeToTopicRequest)

	case unsubscribeFromTopic:
		unsubscribeFromTopicRequest := &subscriber.UnsubscribeFromTopicRequest{}
		if err := unmarshal([]byte(request.Body), unsubscribeFromTopic, request.Header.ContentType); err != nil {
			return nil, err
		}
		return s.UnsubscribeFromTopic(ctx, unsubscribeFromTopicRequest)

	case getSubscribedTopics:
		getSubscribedTopicsRequest := &subscriber.GetSubscribedTopicsRequest{}
		if err := unmarshal([]byte(request.Body), getSubscribedTopicsRequest, request.Header.ContentType); err != nil {
			return nil, err
		}
		return s.GetSubscribedTopics(ctx, getSubscribedTopicsRequest)

	case getMessageFromTopic:
		getMessageFromTopicRequest := &subscriber.GetMessageFromTopicRequest{}
		if err := unmarshal([]byte(request.Body), getMessageFromTopicRequest, request.Header.ContentType); err != nil {
			return nil, err
		}
		return s.GetMessageFromTopic(ctx, getMessageFromTopicRequest)

	default:
		return nil, errors.New("method unimplemented")
	}
}

func unmarshal(data []byte, v interface{}, contentType string) error {
	switch contentType {
	case "json":
		return json.Unmarshal(data, v)
	case "xml":
		return xml.Unmarshal(data, v)
	default:
		return errors.New("unknown content-type")
	}
}

func marshal(v interface{}, contentType string) ([]byte, error) {
	switch contentType {
	case "json":
		return json.Marshal(v)
	case "xml":
		return xml.Marshal(v)
	default:
		return nil, errors.New("unknown content-type")
	}
}

package routes

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"

	pub "github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/publisher"
	sub "github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/subscriber"
	"github.com/WinnersonKharsunai/GraduationProject/server/pkg/protocol"
)

// Handler is the concrete implementation for Router
type Handler struct {
	pSvc pub.PublisherIF
	sSvc sub.SubscriberIF
}

// Router is the interface for the Handler type
type Router interface {
	RequestRouter(ctx context.Context, r string) *protocol.Response
}

// NewHandler is the factory function for the Handler type
func NewHandler(pSvc pub.PublisherIF, sSvc sub.SubscriberIF) Router {
	return &Handler{
		pSvc: pSvc,
		sSvc: sSvc,
	}
}

// RequestRouter handles all the request and response
func (h Handler) RequestRouter(ctx context.Context, r string) *protocol.Response {

	var request protocol.Request

	if err := json.Unmarshal([]byte(r), &request); err != nil {
		return &protocol.Response{Error: err.Error()}
	}

	if err := request.Header.ValidateRequestHeader(ctx); err != nil {
		return &protocol.Response{Error: err.Error()}
	}

	fmt.Println(request.Body)

	return &protocol.Response{}
}

func getRequestType(body string) (string, error) {

	return "", nil
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

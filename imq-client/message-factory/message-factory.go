package messagefactory

import (
	"encoding/json"
	"encoding/xml"
	"errors"
)

const (
	contentTypeJSON    = "json"
	contentTypeXML     = "xml"
	unknownContentType = "unknown content-type"
)

type Messagefactory struct{}

type MessagefactoryIF interface {
	MarshalRequestBody(v interface{}, contentType string) ([]byte, error)
	UnmarshalRequestBody(data []byte, v interface{}, contentType string) error
}

func NewMessageFactory() MessagefactoryIF {
	return &Messagefactory{}
}

func (m *Messagefactory) MarshalRequestBody(v interface{}, contentType string) ([]byte, error) {
	switch contentType {
	case contentTypeJSON:
		return json.Marshal(v)
	case contentTypeXML:
		return xml.Marshal(v)
	default:
		return nil, errors.New(unknownContentType)
	}
}

func (m *Messagefactory) UnmarshalRequestBody(data []byte, v interface{}, contentType string) error {
	switch contentType {
	case contentTypeJSON:
		return json.Unmarshal(data, v)
	case contentTypeXML:
		return xml.Unmarshal(data, v)
	default:
		return errors.New(unknownContentType)
	}
}

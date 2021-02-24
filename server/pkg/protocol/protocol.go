package protocol

import (
	"context"
	"errors"
)

// Protocol ..
type Protocol interface {
	ValidateRequestHeader(ctx context.Context) error
}

// ValidateRequestHeader ...
func (hdr Header) ValidateRequestHeader(ctx context.Context) error {

	if hdr.Version != "1.0" {
		return errors.New("invalid header version")
	}

	if hdr.ContentType != "json" && hdr.ContentType != "xml" {
		return errors.New("unsupported content-type")
	}

	return nil
}

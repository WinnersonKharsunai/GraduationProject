package protocol

import (
	"context"
	"errors"
)

// ProtocolIF is the interface for the protocol
type ProtocolIF interface {
	ValidateRequestHeader(ctx context.Context) error
}

// ValidateRequestHeader validates the request header
func (hdr Header) ValidateRequestHeader(ctx context.Context) error {
	if hdr.Version != "1.0" {
		return errors.New("invalid header version")
	}

	if hdr.ContentType != "json" && hdr.ContentType != "xml" {
		return errors.New("unsupported content-type")
	}

	return nil
}

package protocol

import (
	"context"
	"errors"
)

// Protocol ...
type Protocol struct{}

// ValidateRequestHeader ...
func (p *Protocol) ValidateRequestHeader(ctx context.Context, hdr Header) error {

	if hdr.Version != "1.0" {
		return errors.New("invalid header version")
	}

	if hdr.ContentType != "json" && hdr.ContentType != "xml" {
		return errors.New("unsupported content-type")
	}

	return nil
}

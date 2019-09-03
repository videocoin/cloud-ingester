package hookd

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
)

// Common ingest errors
var (
	ErrEmptyStream            = fmt.Errorf("stream is empty")
	ErrInvalidStream          = fmt.Errorf("invalid stream name")
	ErrInvalidWalletAddress   = fmt.Errorf("invalid user id")
	ErrInvalidContractAddress = fmt.Errorf("invalid contract address")
)

// StreamInfo used to parsing incoming rtmp stream
type StreamInfo struct {
	JobID string
}

// ParseStreamName parses stream info from rtmp url
func ParseStreamName(ctx context.Context, id string) (*StreamInfo, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ParseStreamName")
	defer span.Finish()
	if id == "" {
		return nil, ErrEmptyStream
	}

	return &StreamInfo{JobID: id}, nil
}

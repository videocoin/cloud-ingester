package hookd

import (
	"fmt"
	"strconv"
	"strings"
)

// Common ingest errors
var (
	ErrEmptyStream     = fmt.Errorf("stream is empty")
	ErrInvalidStream   = fmt.Errorf("invalid stream name")
	ErrInvalidUserID   = fmt.Errorf("invalid user id")
	ErrInvalidCameraID = fmt.Errorf("invalid camera id")
)

// StreamInfo used to parsing incoming rtmp stream
type StreamInfo struct {
	UserID   uint32
	CameraID string
}

// ParseStreamName parses stream info from rtmp url
func ParseStreamName(name string) (*StreamInfo, error) {
	if name == "" {
		return nil, ErrEmptyStream
	}

	parts := strings.Split(name, "-")
	if len(parts) != 2 {
		return nil, ErrInvalidStream
	}

	streamInfo := new(StreamInfo)
	userID, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	streamInfo.UserID = uint32(userID)
	streamInfo.CameraID = parts[1]

	fmt.Printf("%+v", parts)

	if streamInfo.UserID == 0 {
		return nil, ErrInvalidUserID
	}

	if len(streamInfo.CameraID) != 12 {
		return nil, ErrInvalidCameraID
	}

	return streamInfo, nil
}

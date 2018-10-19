package hookd

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrEmptyStream     = errors.New("stream is empty")
	ErrInvalidStream   = errors.New("invalid stream name")
	ErrInvalidUserID   = errors.New("invalid user id")
	ErrInvalidCameraID = errors.New("invalid camera id")
)

type StreamInfo struct {
	UserID   uint32
	CameraID string
}

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

	if streamInfo.UserID == 0 {
		return nil, ErrInvalidUserID
	}

	if len(streamInfo.CameraID) != 12 {
		return nil, ErrInvalidCameraID
	}

	return streamInfo, nil
}

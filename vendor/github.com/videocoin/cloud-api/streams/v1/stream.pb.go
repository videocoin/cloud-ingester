// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: streams/v1/stream.proto

package v1

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	golang_proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// The status of a stream is affected by API calls made on a stream object.
type StreamStatus int32

const (
	// Status not set
	StreamStatusNone StreamStatus = 0
	// Initially created and no actions have been taken
	StreamStatusNew StreamStatus = 1
	// Running and preparing input and output destinations
	StreamStatusPreparing StreamStatus = 2
	// Preparation is finished and is ready to consume input data
	StreamStatusPrepared StreamStatus = 3
	// Receiving data and pending on miner to be assigned to stream
	StreamStatusPending StreamStatus = 4
	// Miner has started work on stream, but output is not ready for use
	StreamStatusProcessing StreamStatus = 5
	// Output destination is ready to be consumed
	StreamStatusReady StreamStatus = 6
	// Stream has successfully transcoded video and is now complete
	StreamStatusCompleted StreamStatus = 7
	// Stream has not yet received any input data and has been cancelled
	StreamStatusCancelled StreamStatus = 8
	// Stream has attempted to transcode video received, but problems with the transcoder or account caused it to fail
	StreamStatusFailed  StreamStatus = 9
	StreamStatusDeleted StreamStatus = 10
)

var StreamStatus_name = map[int32]string{
	0:  "STREAM_STATUS_NONE",
	1:  "STREAM_STATUS_NEW",
	2:  "STREAM_STATUS_PREPARING",
	3:  "STREAM_STATUS_PREPARED",
	4:  "STREAM_STATUS_PENDING",
	5:  "STREAM_STATUS_PROCESSING",
	6:  "STREAM_STATUS_READY",
	7:  "STREAM_STATUS_COMPLETED",
	8:  "STREAM_STATUS_CANCELLED",
	9:  "STREAM_STATUS_FAILED",
	10: "STREAM_STATUS_DELETED",
}

var StreamStatus_value = map[string]int32{
	"STREAM_STATUS_NONE":       0,
	"STREAM_STATUS_NEW":        1,
	"STREAM_STATUS_PREPARING":  2,
	"STREAM_STATUS_PREPARED":   3,
	"STREAM_STATUS_PENDING":    4,
	"STREAM_STATUS_PROCESSING": 5,
	"STREAM_STATUS_READY":      6,
	"STREAM_STATUS_COMPLETED":  7,
	"STREAM_STATUS_CANCELLED":  8,
	"STREAM_STATUS_FAILED":     9,
	"STREAM_STATUS_DELETED":    10,
}

func (x StreamStatus) String() string {
	return proto.EnumName(StreamStatus_name, int32(x))
}

func (StreamStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_75b77e6ca64568e5, []int{0}
}

// The status of a stream's ingest is affected by the state of the encoder that's sending video data to the stream.
type InputStatus int32

const (
	// The stream has been created or has ended and is not receiving any input
	InputStatusNone InputStatus = 0
	// Ingest is awaiting for incoming data
	InputStatusPending InputStatus = 1
	// Ingest is receiving data
	InputStatusActive InputStatus = 2
	// Ingest has been failed to process incoming data
	InputStatusError InputStatus = 3
)

var InputStatus_name = map[int32]string{
	0: "INPUT_STATUS_NONE",
	1: "INPUT_STATUS_PENDING",
	2: "INPUT_STATUS_ACTIVE",
	3: "INPUT_STATUS_ERROR",
}

var InputStatus_value = map[string]int32{
	"INPUT_STATUS_NONE":    0,
	"INPUT_STATUS_PENDING": 1,
	"INPUT_STATUS_ACTIVE":  2,
	"INPUT_STATUS_ERROR":   3,
}

func (x InputStatus) String() string {
	return proto.EnumName(InputStatus_name, int32(x))
}

func (InputStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_75b77e6ca64568e5, []int{1}
}

type InputType int32

const (
	InputTypeRTMP   InputType = 0
	InputTypeWebRTC InputType = 1
	InputTypeFile   InputType = 2
)

var InputType_name = map[int32]string{
	0: "INPUT_TYPE_RTMP",
	1: "INPUT_TYPE_WEBRTC",
	2: "INPUT_TYPE_FILE",
}

var InputType_value = map[string]int32{
	"INPUT_TYPE_RTMP":   0,
	"INPUT_TYPE_WEBRTC": 1,
	"INPUT_TYPE_FILE":   2,
}

func (x InputType) String() string {
	return proto.EnumName(InputType_name, int32(x))
}

func (InputType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_75b77e6ca64568e5, []int{2}
}

type OutputType int32

const (
	OutputTypeHLS OutputType = 0
)

var OutputType_name = map[int32]string{
	0: "OUTPUT_TYPE_HLS",
}

var OutputType_value = map[string]int32{
	"OUTPUT_TYPE_HLS": 0,
}

func (x OutputType) String() string {
	return proto.EnumName(OutputType_name, int32(x))
}

func (OutputType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_75b77e6ca64568e5, []int{3}
}

func init() {
	proto.RegisterEnum("cloud.api.streams.v1.StreamStatus", StreamStatus_name, StreamStatus_value)
	golang_proto.RegisterEnum("cloud.api.streams.v1.StreamStatus", StreamStatus_name, StreamStatus_value)
	proto.RegisterEnum("cloud.api.streams.v1.InputStatus", InputStatus_name, InputStatus_value)
	golang_proto.RegisterEnum("cloud.api.streams.v1.InputStatus", InputStatus_name, InputStatus_value)
	proto.RegisterEnum("cloud.api.streams.v1.InputType", InputType_name, InputType_value)
	golang_proto.RegisterEnum("cloud.api.streams.v1.InputType", InputType_name, InputType_value)
	proto.RegisterEnum("cloud.api.streams.v1.OutputType", OutputType_name, OutputType_value)
	golang_proto.RegisterEnum("cloud.api.streams.v1.OutputType", OutputType_name, OutputType_value)
}

func init() { proto.RegisterFile("streams/v1/stream.proto", fileDescriptor_75b77e6ca64568e5) }
func init() { golang_proto.RegisterFile("streams/v1/stream.proto", fileDescriptor_75b77e6ca64568e5) }

var fileDescriptor_75b77e6ca64568e5 = []byte{
	// 645 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x94, 0x4f, 0x6e, 0x9b, 0x40,
	0x14, 0x87, 0x43, 0x92, 0xa6, 0xcd, 0xb4, 0x55, 0x30, 0x76, 0xfe, 0x74, 0x16, 0x94, 0x55, 0x17,
	0x56, 0x82, 0x93, 0xb6, 0xaa, 0xba, 0x25, 0x30, 0x6e, 0x2c, 0x39, 0x18, 0x01, 0x69, 0x94, 0x6e,
	0x22, 0x6c, 0x4f, 0x29, 0x12, 0x66, 0x10, 0x8c, 0x1d, 0xe5, 0x06, 0x15, 0x47, 0xa8, 0xc4, 0xaa,
	0x3d, 0x45, 0x57, 0x5d, 0x66, 0xd9, 0x13, 0x54, 0x51, 0x72, 0x91, 0x8a, 0x01, 0xdb, 0x60, 0x7b,
	0xe5, 0x19, 0xeb, 0xfb, 0xde, 0xfc, 0xe6, 0xcd, 0xb3, 0xc1, 0x7e, 0x4c, 0x23, 0xec, 0x8c, 0xe2,
	0xd6, 0xe4, 0xa4, 0x95, 0x2f, 0xe5, 0x30, 0x22, 0x94, 0x08, 0x8d, 0x81, 0x4f, 0xc6, 0x43, 0xd9,
	0x09, 0x3d, 0xb9, 0x40, 0xe4, 0xc9, 0x09, 0x3c, 0x72, 0x3d, 0xfa, 0x6d, 0xdc, 0x97, 0x07, 0x64,
	0xd4, 0x72, 0x89, 0x4b, 0x5a, 0x0c, 0xee, 0x8f, 0xbf, 0xb2, 0x1d, 0xdb, 0xb0, 0x55, 0x5e, 0x04,
	0xbe, 0x76, 0x09, 0x71, 0x7d, 0x3c, 0xa7, 0xa8, 0x37, 0xc2, 0x31, 0x75, 0x46, 0x61, 0x01, 0x1c,
	0xb2, 0x8f, 0xc1, 0x91, 0x8b, 0x83, 0xa3, 0xf8, 0xc6, 0x71, 0x5d, 0x1c, 0xb5, 0x48, 0x48, 0x3d,
	0x12, 0xc4, 0x2d, 0x27, 0x08, 0x08, 0x75, 0xd8, 0x3a, 0xa7, 0x9b, 0xf7, 0x9b, 0xe0, 0x85, 0xc5,
	0xc2, 0x58, 0xd4, 0xa1, 0xe3, 0x58, 0x38, 0x04, 0x82, 0x65, 0x9b, 0x48, 0x39, 0xbf, 0xb6, 0x6c,
	0xc5, 0xbe, 0xb0, 0xae, 0xf5, 0x9e, 0x8e, 0xf8, 0x35, 0xd8, 0x48, 0x52, 0x89, 0x2f, 0x93, 0x3a,
	0x09, 0xb0, 0xd0, 0x04, 0xb5, 0x05, 0x1a, 0x5d, 0xf2, 0x1c, 0xac, 0x27, 0xa9, 0xb4, 0x53, 0x81,
	0xf1, 0x8d, 0xf0, 0x01, 0xec, 0x57, 0x59, 0xc3, 0x44, 0x86, 0x62, 0x76, 0xf4, 0x4f, 0xfc, 0x3a,
	0x7c, 0x95, 0xa4, 0xd2, 0x6e, 0xd9, 0x30, 0x22, 0x1c, 0x3a, 0x91, 0x17, 0xb8, 0xc2, 0x7b, 0xb0,
	0xb7, 0xca, 0x43, 0x1a, 0xbf, 0x01, 0x0f, 0x92, 0x54, 0x6a, 0x2c, 0x6b, 0x78, 0x28, 0xbc, 0x05,
	0xbb, 0x0b, 0x16, 0xd2, 0xb5, 0xec, 0xac, 0x4d, 0xb8, 0x9f, 0xa4, 0x52, 0xbd, 0x22, 0xe1, 0x60,
	0x98, 0x9d, 0xf4, 0x11, 0x1c, 0x2c, 0x9e, 0xd4, 0x53, 0x91, 0x65, 0x65, 0xda, 0x13, 0x08, 0x93,
	0x54, 0xda, 0xab, 0x9e, 0x45, 0x06, 0x38, 0x8e, 0x33, 0x53, 0x06, 0xf5, 0xaa, 0x69, 0x22, 0x45,
	0xbb, 0xe2, 0xb7, 0xe0, 0x6e, 0x92, 0x4a, 0xb5, 0xb2, 0x64, 0x62, 0x67, 0x78, 0xbb, 0xdc, 0x0b,
	0xb5, 0x77, 0x6e, 0x74, 0x91, 0x8d, 0x34, 0xfe, 0xe9, 0x72, 0x2f, 0x54, 0x32, 0x0a, 0x7d, 0x4c,
	0xf1, 0x70, 0x85, 0xa7, 0xe8, 0x2a, 0xea, 0x76, 0x91, 0xc6, 0x3f, 0x5b, 0xe1, 0x39, 0xc1, 0x00,
	0xfb, 0x3e, 0x1e, 0x0a, 0xc7, 0xa0, 0x51, 0xf5, 0xda, 0x4a, 0x27, 0x93, 0xb6, 0xe1, 0x5e, 0x92,
	0x4a, 0x42, 0x59, 0x6a, 0x3b, 0x9e, 0xbf, 0xaa, 0x7f, 0x1a, 0xca, 0xf3, 0x81, 0xe5, 0xfe, 0x69,
	0x98, 0xa5, 0x83, 0x8d, 0xef, 0x3f, 0xc5, 0xb5, 0xdf, 0xbf, 0xc4, 0xca, 0x44, 0x35, 0xff, 0x71,
	0xe0, 0x79, 0x27, 0x08, 0xc7, 0xb4, 0x98, 0xb0, 0x26, 0xa8, 0x75, 0x74, 0xe3, 0xc2, 0x5e, 0x18,
	0x30, 0x36, 0x33, 0x25, 0x8e, 0xcd, 0xd7, 0x31, 0x68, 0x54, 0xd8, 0xe9, 0x23, 0x72, 0x79, 0xee,
	0x12, 0x3e, 0x7d, 0x43, 0x19, 0xd4, 0x2b, 0x86, 0xa2, 0xda, 0x9d, 0xcf, 0x88, 0x5f, 0xcf, 0x5f,
	0xa2, 0x24, 0x28, 0x03, 0xea, 0x4d, 0x70, 0x36, 0xef, 0x15, 0x1e, 0x99, 0x66, 0xcf, 0xe4, 0x37,
	0xf2, 0x79, 0x2f, 0xe1, 0x28, 0x8a, 0x48, 0x04, 0xeb, 0xc5, 0x0d, 0xcb, 0x17, 0x6a, 0xfe, 0xe0,
	0xc0, 0x36, 0xdb, 0xdb, 0xb7, 0x21, 0x16, 0xde, 0x80, 0x9d, 0xbc, 0xa0, 0x7d, 0x65, 0xa0, 0x6b,
	0xd3, 0x3e, 0x37, 0xf8, 0x35, 0x58, 0x4b, 0x52, 0xe9, 0xe5, 0x8c, 0xc9, 0xbe, 0x9c, 0xb7, 0x81,
	0x71, 0x97, 0xe8, 0xd4, 0xb4, 0xd5, 0xe9, 0x4f, 0x67, 0x46, 0x5e, 0xe2, 0xbe, 0x69, 0xab, 0x0b,
	0x35, 0xdb, 0x9d, 0x6e, 0x76, 0xa1, 0x6a, 0xcd, 0xb6, 0xe7, 0x63, 0x58, 0x2b, 0xe2, 0xcd, 0xe3,
	0x34, 0xcf, 0x00, 0xe8, 0x8d, 0x69, 0x29, 0x5c, 0xef, 0xc2, 0x9e, 0x55, 0x3a, 0xeb, 0x5a, 0xd3,
	0x70, 0x73, 0xe8, 0xac, 0x6b, 0x41, 0xa1, 0x28, 0x54, 0x72, 0x4f, 0x0f, 0xee, 0x1e, 0x44, 0xee,
	0xef, 0x83, 0xc8, 0xdd, 0x3f, 0x88, 0xdc, 0x9f, 0x47, 0x91, 0xbb, 0x7b, 0x14, 0xb9, 0x2f, 0xeb,
	0x93, 0x93, 0xfe, 0x16, 0xfb, 0x2f, 0x79, 0xf7, 0x3f, 0x00, 0x00, 0xff, 0xff, 0xf8, 0xea, 0x5e,
	0x7e, 0xfa, 0x04, 0x00, 0x00,
}

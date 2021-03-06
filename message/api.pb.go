//
//此文件用于定义通信协议号，通常情况下 req/ack 为一组api
//服务端主动发起的协议以NOTIFY结尾

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.9.2
// source: api.proto

//
//用于c/s协议解析

package message

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ChatServiceAPI int32

const (
	ChatServiceAPI_SERVICE_API_UNKNOWN ChatServiceAPI = 0
	//------ login
	ChatServiceAPI_SERVICE_API_LOGIN_REQ ChatServiceAPI = 1
	ChatServiceAPI_SERVICE_API_LOGIN_ACK ChatServiceAPI = 2
	//------ send
	ChatServiceAPI_SERVICE_API_SEND_REQ ChatServiceAPI = 3
	ChatServiceAPI_SERVICE_API_SEND_ACK ChatServiceAPI = 4
	//------ query
	ChatServiceAPI_SERVICE_API_STATS_REQ        ChatServiceAPI = 5
	ChatServiceAPI_SERVICE_API_STATS_ACK        ChatServiceAPI = 6
	ChatServiceAPI_SERVICE_API_POPULAR_WORD_REQ ChatServiceAPI = 7
	ChatServiceAPI_SERVICE_API_POPULAR_WORD_ACK ChatServiceAPI = 8
	//------ notify
	ChatServiceAPI_SERVICE_API_MSG_NOTIFY ChatServiceAPI = 9
)

// Enum value maps for ChatServiceAPI.
var (
	ChatServiceAPI_name = map[int32]string{
		0: "SERVICE_API_UNKNOWN",
		1: "SERVICE_API_LOGIN_REQ",
		2: "SERVICE_API_LOGIN_ACK",
		3: "SERVICE_API_SEND_REQ",
		4: "SERVICE_API_SEND_ACK",
		5: "SERVICE_API_STATS_REQ",
		6: "SERVICE_API_STATS_ACK",
		7: "SERVICE_API_POPULAR_WORD_REQ",
		8: "SERVICE_API_POPULAR_WORD_ACK",
		9: "SERVICE_API_MSG_NOTIFY",
	}
	ChatServiceAPI_value = map[string]int32{
		"SERVICE_API_UNKNOWN":          0,
		"SERVICE_API_LOGIN_REQ":        1,
		"SERVICE_API_LOGIN_ACK":        2,
		"SERVICE_API_SEND_REQ":         3,
		"SERVICE_API_SEND_ACK":         4,
		"SERVICE_API_STATS_REQ":        5,
		"SERVICE_API_STATS_ACK":        6,
		"SERVICE_API_POPULAR_WORD_REQ": 7,
		"SERVICE_API_POPULAR_WORD_ACK": 8,
		"SERVICE_API_MSG_NOTIFY":       9,
	}
)

func (x ChatServiceAPI) Enum() *ChatServiceAPI {
	p := new(ChatServiceAPI)
	*p = x
	return p
}

func (x ChatServiceAPI) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ChatServiceAPI) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_enumTypes[0].Descriptor()
}

func (ChatServiceAPI) Type() protoreflect.EnumType {
	return &file_api_proto_enumTypes[0]
}

func (x ChatServiceAPI) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ChatServiceAPI.Descriptor instead.
func (ChatServiceAPI) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{0}
}

var File_api_proto protoreflect.FileDescriptor

var file_api_proto_rawDesc = []byte{
	0x0a, 0x09, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2a, 0xa9, 0x02, 0x0a, 0x0e, 0x43, 0x68, 0x61, 0x74, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x41, 0x50, 0x49, 0x12, 0x17, 0x0a, 0x13, 0x53, 0x45, 0x52, 0x56, 0x49,
	0x43, 0x45, 0x5f, 0x41, 0x50, 0x49, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00,
	0x12, 0x19, 0x0a, 0x15, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x41, 0x50, 0x49, 0x5f,
	0x4c, 0x4f, 0x47, 0x49, 0x4e, 0x5f, 0x52, 0x45, 0x51, 0x10, 0x01, 0x12, 0x19, 0x0a, 0x15, 0x53,
	0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x41, 0x50, 0x49, 0x5f, 0x4c, 0x4f, 0x47, 0x49, 0x4e,
	0x5f, 0x41, 0x43, 0x4b, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43,
	0x45, 0x5f, 0x41, 0x50, 0x49, 0x5f, 0x53, 0x45, 0x4e, 0x44, 0x5f, 0x52, 0x45, 0x51, 0x10, 0x03,
	0x12, 0x18, 0x0a, 0x14, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x41, 0x50, 0x49, 0x5f,
	0x53, 0x45, 0x4e, 0x44, 0x5f, 0x41, 0x43, 0x4b, 0x10, 0x04, 0x12, 0x19, 0x0a, 0x15, 0x53, 0x45,
	0x52, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x41, 0x50, 0x49, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x53, 0x5f,
	0x52, 0x45, 0x51, 0x10, 0x05, 0x12, 0x19, 0x0a, 0x15, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45,
	0x5f, 0x41, 0x50, 0x49, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x53, 0x5f, 0x41, 0x43, 0x4b, 0x10, 0x06,
	0x12, 0x20, 0x0a, 0x1c, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x41, 0x50, 0x49, 0x5f,
	0x50, 0x4f, 0x50, 0x55, 0x4c, 0x41, 0x52, 0x5f, 0x57, 0x4f, 0x52, 0x44, 0x5f, 0x52, 0x45, 0x51,
	0x10, 0x07, 0x12, 0x20, 0x0a, 0x1c, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x41, 0x50,
	0x49, 0x5f, 0x50, 0x4f, 0x50, 0x55, 0x4c, 0x41, 0x52, 0x5f, 0x57, 0x4f, 0x52, 0x44, 0x5f, 0x41,
	0x43, 0x4b, 0x10, 0x08, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x5f,
	0x41, 0x50, 0x49, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x4e, 0x4f, 0x54, 0x49, 0x46, 0x59, 0x10, 0x09,
	0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e,
	0x65, 0x72, 0x6f, 0x72, 0x65, 0x64, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2d, 0x74, 0x65, 0x73, 0x74,
	0x2d, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_rawDescOnce sync.Once
	file_api_proto_rawDescData = file_api_proto_rawDesc
)

func file_api_proto_rawDescGZIP() []byte {
	file_api_proto_rawDescOnce.Do(func() {
		file_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_rawDescData)
	})
	return file_api_proto_rawDescData
}

var file_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_proto_goTypes = []interface{}{
	(ChatServiceAPI)(0), // 0: message.ChatServiceAPI
}
var file_api_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_init() }
func file_api_proto_init() {
	if File_api_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_goTypes,
		DependencyIndexes: file_api_proto_depIdxs,
		EnumInfos:         file_api_proto_enumTypes,
	}.Build()
	File_api_proto = out.File
	file_api_proto_rawDesc = nil
	file_api_proto_goTypes = nil
	file_api_proto_depIdxs = nil
}

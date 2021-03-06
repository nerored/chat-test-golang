//
//此文件用于定义服务端主动推送

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.9.2
// source: ntf.proto

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

type MESSAGE_TYPE int32

const (
	MESSAGE_TYPE_MESSAGE_TYPE_UNKNOWN   MESSAGE_TYPE = 0
	MESSAGE_TYPE_MESSAGE_TYPE_PRIVATE   MESSAGE_TYPE = 1
	MESSAGE_TYPE_MESSAGE_TYPE_BROADCAST MESSAGE_TYPE = 2
)

// Enum value maps for MESSAGE_TYPE.
var (
	MESSAGE_TYPE_name = map[int32]string{
		0: "MESSAGE_TYPE_UNKNOWN",
		1: "MESSAGE_TYPE_PRIVATE",
		2: "MESSAGE_TYPE_BROADCAST",
	}
	MESSAGE_TYPE_value = map[string]int32{
		"MESSAGE_TYPE_UNKNOWN":   0,
		"MESSAGE_TYPE_PRIVATE":   1,
		"MESSAGE_TYPE_BROADCAST": 2,
	}
)

func (x MESSAGE_TYPE) Enum() *MESSAGE_TYPE {
	p := new(MESSAGE_TYPE)
	*p = x
	return p
}

func (x MESSAGE_TYPE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MESSAGE_TYPE) Descriptor() protoreflect.EnumDescriptor {
	return file_ntf_proto_enumTypes[0].Descriptor()
}

func (MESSAGE_TYPE) Type() protoreflect.EnumType {
	return &file_ntf_proto_enumTypes[0]
}

func (x MESSAGE_TYPE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MESSAGE_TYPE.Descriptor instead.
func (MESSAGE_TYPE) EnumDescriptor() ([]byte, []int) {
	return file_ntf_proto_rawDescGZIP(), []int{0}
}

type API_NEW_MESSAGE_NOTIFY struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From     string       `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	Msg      string       `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	SendUnix int64        `protobuf:"varint,3,opt,name=sendUnix,proto3" json:"sendUnix,omitempty"`
	Type     MESSAGE_TYPE `protobuf:"varint,4,opt,name=type,proto3,enum=message.MESSAGE_TYPE" json:"type,omitempty"`
}

func (x *API_NEW_MESSAGE_NOTIFY) Reset() {
	*x = API_NEW_MESSAGE_NOTIFY{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ntf_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *API_NEW_MESSAGE_NOTIFY) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*API_NEW_MESSAGE_NOTIFY) ProtoMessage() {}

func (x *API_NEW_MESSAGE_NOTIFY) ProtoReflect() protoreflect.Message {
	mi := &file_ntf_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use API_NEW_MESSAGE_NOTIFY.ProtoReflect.Descriptor instead.
func (*API_NEW_MESSAGE_NOTIFY) Descriptor() ([]byte, []int) {
	return file_ntf_proto_rawDescGZIP(), []int{0}
}

func (x *API_NEW_MESSAGE_NOTIFY) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *API_NEW_MESSAGE_NOTIFY) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *API_NEW_MESSAGE_NOTIFY) GetSendUnix() int64 {
	if x != nil {
		return x.SendUnix
	}
	return 0
}

func (x *API_NEW_MESSAGE_NOTIFY) GetType() MESSAGE_TYPE {
	if x != nil {
		return x.Type
	}
	return MESSAGE_TYPE_MESSAGE_TYPE_UNKNOWN
}

var File_ntf_proto protoreflect.FileDescriptor

var file_ntf_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6e, 0x74, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x85, 0x01, 0x0a, 0x16, 0x41, 0x50, 0x49, 0x5f, 0x4e, 0x45, 0x57,
	0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x4e, 0x4f, 0x54, 0x49, 0x46, 0x59, 0x12,
	0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66,
	0x72, 0x6f, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x55, 0x6e, 0x69,
	0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x55, 0x6e, 0x69,
	0x78, 0x12, 0x29, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x15, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47,
	0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x2a, 0x5e, 0x0a, 0x0c,
	0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x12, 0x18, 0x0a, 0x14,
	0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x4b,
	0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47,
	0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x50, 0x52, 0x49, 0x56, 0x41, 0x54, 0x45, 0x10, 0x01,
	0x12, 0x1a, 0x0a, 0x16, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x42, 0x52, 0x4f, 0x41, 0x44, 0x43, 0x41, 0x53, 0x54, 0x10, 0x02, 0x42, 0x2d, 0x5a, 0x2b,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x65, 0x72, 0x6f, 0x72,
	0x65, 0x64, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2d, 0x74, 0x65, 0x73, 0x74, 0x2d, 0x67, 0x6f, 0x6c,
	0x61, 0x6e, 0x67, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_ntf_proto_rawDescOnce sync.Once
	file_ntf_proto_rawDescData = file_ntf_proto_rawDesc
)

func file_ntf_proto_rawDescGZIP() []byte {
	file_ntf_proto_rawDescOnce.Do(func() {
		file_ntf_proto_rawDescData = protoimpl.X.CompressGZIP(file_ntf_proto_rawDescData)
	})
	return file_ntf_proto_rawDescData
}

var file_ntf_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_ntf_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ntf_proto_goTypes = []interface{}{
	(MESSAGE_TYPE)(0),              // 0: message.MESSAGE_TYPE
	(*API_NEW_MESSAGE_NOTIFY)(nil), // 1: message.API_NEW_MESSAGE_NOTIFY
}
var file_ntf_proto_depIdxs = []int32{
	0, // 0: message.API_NEW_MESSAGE_NOTIFY.type:type_name -> message.MESSAGE_TYPE
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ntf_proto_init() }
func file_ntf_proto_init() {
	if File_ntf_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ntf_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*API_NEW_MESSAGE_NOTIFY); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ntf_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ntf_proto_goTypes,
		DependencyIndexes: file_ntf_proto_depIdxs,
		EnumInfos:         file_ntf_proto_enumTypes,
		MessageInfos:      file_ntf_proto_msgTypes,
	}.Build()
	File_ntf_proto = out.File
	file_ntf_proto_rawDesc = nil
	file_ntf_proto_goTypes = nil
	file_ntf_proto_depIdxs = nil
}

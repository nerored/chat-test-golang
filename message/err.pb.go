//
//此文件用于定义通信错误码

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.9.2
// source: err.proto

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

type ErrCode int32

const (
	ErrCode_ERR_CODE_FAILED            ErrCode = 0 // 通常请求失败,原因不指定
	ErrCode_ERR_CODE_SUCCESS           ErrCode = 1
	ErrCode_ERR_CODE_DUPLICATE_NAME    ErrCode = 2
	ErrCode_ERR_CODE_USER_IS_NOT_EXIST ErrCode = 3
)

// Enum value maps for ErrCode.
var (
	ErrCode_name = map[int32]string{
		0: "ERR_CODE_FAILED",
		1: "ERR_CODE_SUCCESS",
		2: "ERR_CODE_DUPLICATE_NAME",
		3: "ERR_CODE_USER_IS_NOT_EXIST",
	}
	ErrCode_value = map[string]int32{
		"ERR_CODE_FAILED":            0,
		"ERR_CODE_SUCCESS":           1,
		"ERR_CODE_DUPLICATE_NAME":    2,
		"ERR_CODE_USER_IS_NOT_EXIST": 3,
	}
)

func (x ErrCode) Enum() *ErrCode {
	p := new(ErrCode)
	*p = x
	return p
}

func (x ErrCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrCode) Descriptor() protoreflect.EnumDescriptor {
	return file_err_proto_enumTypes[0].Descriptor()
}

func (ErrCode) Type() protoreflect.EnumType {
	return &file_err_proto_enumTypes[0]
}

func (x ErrCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrCode.Descriptor instead.
func (ErrCode) EnumDescriptor() ([]byte, []int) {
	return file_err_proto_rawDescGZIP(), []int{0}
}

var File_err_proto protoreflect.FileDescriptor

var file_err_proto_rawDesc = []byte{
	0x0a, 0x09, 0x65, 0x72, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2a, 0x71, 0x0a, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12,
	0x13, 0x0a, 0x0f, 0x45, 0x52, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x45, 0x52, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45,
	0x5f, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x01, 0x12, 0x1b, 0x0a, 0x17, 0x45, 0x52,
	0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x44, 0x55, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54, 0x45,
	0x5f, 0x4e, 0x41, 0x4d, 0x45, 0x10, 0x02, 0x12, 0x1e, 0x0a, 0x1a, 0x45, 0x52, 0x52, 0x5f, 0x43,
	0x4f, 0x44, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x49, 0x53, 0x5f, 0x4e, 0x4f, 0x54, 0x5f,
	0x45, 0x58, 0x49, 0x53, 0x54, 0x10, 0x03, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x65, 0x72, 0x6f, 0x72, 0x65, 0x64, 0x2f, 0x63, 0x68,
	0x61, 0x74, 0x2d, 0x74, 0x65, 0x73, 0x74, 0x2d, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_err_proto_rawDescOnce sync.Once
	file_err_proto_rawDescData = file_err_proto_rawDesc
)

func file_err_proto_rawDescGZIP() []byte {
	file_err_proto_rawDescOnce.Do(func() {
		file_err_proto_rawDescData = protoimpl.X.CompressGZIP(file_err_proto_rawDescData)
	})
	return file_err_proto_rawDescData
}

var file_err_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_err_proto_goTypes = []interface{}{
	(ErrCode)(0), // 0: message.ErrCode
}
var file_err_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_err_proto_init() }
func file_err_proto_init() {
	if File_err_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_err_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_err_proto_goTypes,
		DependencyIndexes: file_err_proto_depIdxs,
		EnumInfos:         file_err_proto_enumTypes,
	}.Build()
	File_err_proto = out.File
	file_err_proto_rawDesc = nil
	file_err_proto_goTypes = nil
	file_err_proto_depIdxs = nil
}

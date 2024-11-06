// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.3
// source: api/v1/proto/helloworld/error.proto

package helloworld

import (
	_ "github.com/go-kratos/kratos/v2/errors"
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

type HelloWorldErrorReason int32

const (
	HelloWorldErrorReason_HELLO_WORLD_ERROR_REASON_GREETER_UNSPECIFIED HelloWorldErrorReason = 0
	HelloWorldErrorReason_HELLO_WORLD_ERROR_REASON_USER_NOT_FOUND      HelloWorldErrorReason = 1
)

// Enum value maps for HelloWorldErrorReason.
var (
	HelloWorldErrorReason_name = map[int32]string{
		0: "HELLO_WORLD_ERROR_REASON_GREETER_UNSPECIFIED",
		1: "HELLO_WORLD_ERROR_REASON_USER_NOT_FOUND",
	}
	HelloWorldErrorReason_value = map[string]int32{
		"HELLO_WORLD_ERROR_REASON_GREETER_UNSPECIFIED": 0,
		"HELLO_WORLD_ERROR_REASON_USER_NOT_FOUND":      1,
	}
)

func (x HelloWorldErrorReason) Enum() *HelloWorldErrorReason {
	p := new(HelloWorldErrorReason)
	*p = x
	return p
}

func (x HelloWorldErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HelloWorldErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_api_v1_proto_helloworld_error_proto_enumTypes[0].Descriptor()
}

func (HelloWorldErrorReason) Type() protoreflect.EnumType {
	return &file_api_v1_proto_helloworld_error_proto_enumTypes[0]
}

func (x HelloWorldErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HelloWorldErrorReason.Descriptor instead.
func (HelloWorldErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_api_v1_proto_helloworld_error_proto_rawDescGZIP(), []int{0}
}

var File_api_v1_proto_helloworld_error_proto protoreflect.FileDescriptor

var file_api_v1_proto_helloworld_error_proto_rawDesc = []byte{
	0x0a, 0x23, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x68,
	0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c,
	0x64, 0x1a, 0x13, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x7c, 0x0a, 0x15, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x57,
	0x6f, 0x72, 0x6c, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12,
	0x30, 0x0a, 0x2c, 0x48, 0x45, 0x4c, 0x4c, 0x4f, 0x5f, 0x57, 0x4f, 0x52, 0x4c, 0x44, 0x5f, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x5f, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x47, 0x52, 0x45, 0x45,
	0x54, 0x45, 0x52, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x2b, 0x0a, 0x27, 0x48, 0x45, 0x4c, 0x4c, 0x4f, 0x5f, 0x57, 0x4f, 0x52, 0x4c, 0x44,
	0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x55, 0x53,
	0x45, 0x52, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x01, 0x1a, 0x04,
	0xa0, 0x45, 0xf4, 0x03, 0x42, 0x53, 0x0a, 0x28, 0x63, 0x6f, 0x6d, 0x2e, 0x6f, 0x72, 0x69, 0x67,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64,
	0x50, 0x01, 0x5a, 0x25, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x3b, 0x68,
	0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_api_v1_proto_helloworld_error_proto_rawDescOnce sync.Once
	file_api_v1_proto_helloworld_error_proto_rawDescData = file_api_v1_proto_helloworld_error_proto_rawDesc
)

func file_api_v1_proto_helloworld_error_proto_rawDescGZIP() []byte {
	file_api_v1_proto_helloworld_error_proto_rawDescOnce.Do(func() {
		file_api_v1_proto_helloworld_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_proto_helloworld_error_proto_rawDescData)
	})
	return file_api_v1_proto_helloworld_error_proto_rawDescData
}

var file_api_v1_proto_helloworld_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_v1_proto_helloworld_error_proto_goTypes = []any{
	(HelloWorldErrorReason)(0), // 0: api.v1.services.helloworld.HelloWorldErrorReason
}
var file_api_v1_proto_helloworld_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_v1_proto_helloworld_error_proto_init() }
func file_api_v1_proto_helloworld_error_proto_init() {
	if File_api_v1_proto_helloworld_error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_proto_helloworld_error_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1_proto_helloworld_error_proto_goTypes,
		DependencyIndexes: file_api_v1_proto_helloworld_error_proto_depIdxs,
		EnumInfos:         file_api_v1_proto_helloworld_error_proto_enumTypes,
	}.Build()
	File_api_v1_proto_helloworld_error_proto = out.File
	file_api_v1_proto_helloworld_error_proto_rawDesc = nil
	file_api_v1_proto_helloworld_error_proto_goTypes = nil
	file_api_v1_proto_helloworld_error_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: api/rpc/Presence.proto

package rpc

import (
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

type WhoseOnReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VoiceChannel string `protobuf:"bytes,1,opt,name=voiceChannel,proto3" json:"voiceChannel,omitempty"`
}

func (x *WhoseOnReq) Reset() {
	*x = WhoseOnReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_rpc_Presence_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WhoseOnReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WhoseOnReq) ProtoMessage() {}

func (x *WhoseOnReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_rpc_Presence_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WhoseOnReq.ProtoReflect.Descriptor instead.
func (*WhoseOnReq) Descriptor() ([]byte, []int) {
	return file_api_rpc_Presence_proto_rawDescGZIP(), []int{0}
}

func (x *WhoseOnReq) GetVoiceChannel() string {
	if x != nil {
		return x.VoiceChannel
	}
	return ""
}

type WhoseOnResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []string `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *WhoseOnResp) Reset() {
	*x = WhoseOnResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_rpc_Presence_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WhoseOnResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WhoseOnResp) ProtoMessage() {}

func (x *WhoseOnResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_rpc_Presence_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WhoseOnResp.ProtoReflect.Descriptor instead.
func (*WhoseOnResp) Descriptor() ([]byte, []int) {
	return file_api_rpc_Presence_proto_rawDescGZIP(), []int{1}
}

func (x *WhoseOnResp) GetUsers() []string {
	if x != nil {
		return x.Users
	}
	return nil
}

var File_api_rpc_Presence_proto protoreflect.FileDescriptor

var file_api_rpc_Presence_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e,
	0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x62, 0x6f, 0x74, 0x69, 0x73, 0x67,
	0x75, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x22,
	0x30, 0x0a, 0x0a, 0x57, 0x68, 0x6f, 0x73, 0x65, 0x4f, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x22, 0x0a,
	0x0c, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x22, 0x23, 0x0a, 0x0b, 0x57, 0x68, 0x6f, 0x73, 0x65, 0x4f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x32, 0x5c, 0x0a, 0x08, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e,
	0x63, 0x65, 0x12, 0x50, 0x0a, 0x07, 0x57, 0x68, 0x6f, 0x73, 0x65, 0x4f, 0x6e, 0x12, 0x21, 0x2e,
	0x62, 0x6f, 0x74, 0x69, 0x73, 0x67, 0x75, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x65,
	0x73, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x57, 0x68, 0x6f, 0x73, 0x65, 0x4f, 0x6e, 0x52, 0x65, 0x71,
	0x1a, 0x22, 0x2e, 0x62, 0x6f, 0x74, 0x69, 0x73, 0x67, 0x75, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x57, 0x68, 0x6f, 0x73, 0x65, 0x4f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x42, 0x1f, 0x5a, 0x1d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x62, 0x6f, 0x74, 0x2d, 0x69, 0x73, 0x2d, 0x67, 0x75, 0x64, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_rpc_Presence_proto_rawDescOnce sync.Once
	file_api_rpc_Presence_proto_rawDescData = file_api_rpc_Presence_proto_rawDesc
)

func file_api_rpc_Presence_proto_rawDescGZIP() []byte {
	file_api_rpc_Presence_proto_rawDescOnce.Do(func() {
		file_api_rpc_Presence_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_rpc_Presence_proto_rawDescData)
	})
	return file_api_rpc_Presence_proto_rawDescData
}

var file_api_rpc_Presence_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_rpc_Presence_proto_goTypes = []interface{}{
	(*WhoseOnReq)(nil),  // 0: botisgud.api.presence.WhoseOnReq
	(*WhoseOnResp)(nil), // 1: botisgud.api.presence.WhoseOnResp
}
var file_api_rpc_Presence_proto_depIdxs = []int32{
	0, // 0: botisgud.api.presence.Presence.WhoseOn:input_type -> botisgud.api.presence.WhoseOnReq
	1, // 1: botisgud.api.presence.Presence.WhoseOn:output_type -> botisgud.api.presence.WhoseOnResp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_rpc_Presence_proto_init() }
func file_api_rpc_Presence_proto_init() {
	if File_api_rpc_Presence_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_rpc_Presence_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WhoseOnReq); i {
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
		file_api_rpc_Presence_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WhoseOnResp); i {
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
			RawDescriptor: file_api_rpc_Presence_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_rpc_Presence_proto_goTypes,
		DependencyIndexes: file_api_rpc_Presence_proto_depIdxs,
		MessageInfos:      file_api_rpc_Presence_proto_msgTypes,
	}.Build()
	File_api_rpc_Presence_proto = out.File
	file_api_rpc_Presence_proto_rawDesc = nil
	file_api_rpc_Presence_proto_goTypes = nil
	file_api_rpc_Presence_proto_depIdxs = nil
}
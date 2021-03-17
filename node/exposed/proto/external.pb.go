// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: proto/external.proto

package exposed

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

type Command int32

const (
	Command_UPLOAD   Command = 0
	Command_DOWNLOAD Command = 1
	Command_CHECK    Command = 2
)

// Enum value maps for Command.
var (
	Command_name = map[int32]string{
		0: "UPLOAD",
		1: "DOWNLOAD",
		2: "CHECK",
	}
	Command_value = map[string]int32{
		"UPLOAD":   0,
		"DOWNLOAD": 1,
		"CHECK":    2,
	}
)

func (x Command) Enum() *Command {
	p := new(Command)
	*p = x
	return p
}

func (x Command) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Command) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_external_proto_enumTypes[0].Descriptor()
}

func (Command) Type() protoreflect.EnumType {
	return &file_proto_external_proto_enumTypes[0]
}

func (x Command) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Command.Descriptor instead.
func (Command) EnumDescriptor() ([]byte, []int) {
	return file_proto_external_proto_rawDescGZIP(), []int{0}
}

// In this proof of concept, we have the key and value sent from client -> chord node
type NewFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key     string  `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value   string  `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Command Command `protobuf:"varint,3,opt,name=command,proto3,enum=grpc.Command" json:"command,omitempty"`
}

func (x *NewFileRequest) Reset() {
	*x = NewFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewFileRequest) ProtoMessage() {}

func (x *NewFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewFileRequest.ProtoReflect.Descriptor instead.
func (*NewFileRequest) Descriptor() ([]byte, []int) {
	return file_proto_external_proto_rawDescGZIP(), []int{0}
}

func (x *NewFileRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *NewFileRequest) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *NewFileRequest) GetCommand() Command {
	if x != nil {
		return x.Command
	}
	return Command_UPLOAD
}

// can merge
type CheckFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key     string  `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Command Command `protobuf:"varint,2,opt,name=command,proto3,enum=grpc.Command" json:"command,omitempty"`
}

func (x *CheckFileRequest) Reset() {
	*x = CheckFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckFileRequest) ProtoMessage() {}

func (x *CheckFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckFileRequest.ProtoReflect.Descriptor instead.
func (*CheckFileRequest) Descriptor() ([]byte, []int) {
	return file_proto_external_proto_rawDescGZIP(), []int{1}
}

func (x *CheckFileRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *CheckFileRequest) GetCommand() Command {
	if x != nil {
		return x.Command
	}
	return Command_UPLOAD
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Redirect bool      `protobuf:"varint,1,opt,name=redirect,proto3" json:"redirect,omitempty"` // set to true if need contact other person
	NodeInfo *NodeInfo `protobuf:"bytes,2,opt,name=nodeInfo,proto3" json:"nodeInfo,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_proto_external_proto_rawDescGZIP(), []int{2}
}

func (x *Response) GetRedirect() bool {
	if x != nil {
		return x.Redirect
	}
	return false
}

func (x *Response) GetNodeInfo() *NodeInfo {
	if x != nil {
		return x.NodeInfo
	}
	return nil
}

// if boolean isFileHere = false, provide the client with the addr of the next node in chord
type NodeInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID string `protobuf:"bytes,1,opt,name=nodeID,proto3" json:"nodeID,omitempty"` // who you should ask next
}

func (x *NodeInfo) Reset() {
	*x = NodeInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeInfo) ProtoMessage() {}

func (x *NodeInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeInfo.ProtoReflect.Descriptor instead.
func (*NodeInfo) Descriptor() ([]byte, []int) {
	return file_proto_external_proto_rawDescGZIP(), []int{3}
}

func (x *NodeInfo) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

var File_proto_external_proto protoreflect.FileDescriptor

var file_proto_external_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x72, 0x70, 0x63, 0x22, 0x61, 0x0a, 0x0e,
	0x4e, 0x65, 0x77, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x27, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22,
	0x4d, 0x0a, 0x10, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x27, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0x52,
	0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x72, 0x65,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x12, 0x2a, 0x0a, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x4e, 0x6f, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x22, 0x22, 0x0a, 0x08, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16,
	0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x2a, 0x2e, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x12, 0x0a, 0x0a, 0x06, 0x55, 0x50, 0x4c, 0x4f, 0x41, 0x44, 0x10, 0x00, 0x12, 0x0c, 0x0a,
	0x08, 0x44, 0x4f, 0x57, 0x4e, 0x4c, 0x4f, 0x41, 0x44, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x43,
	0x48, 0x45, 0x43, 0x4b, 0x10, 0x02, 0x32, 0x80, 0x01, 0x0a, 0x10, 0x45, 0x78, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x12, 0x35, 0x0a, 0x0b, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x4e, 0x65, 0x77, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x35, 0x0a, 0x09, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x46, 0x69, 0x6c, 0x65, 0x12,
	0x16, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x68, 0x65, 0x69, 0x6b, 0x68, 0x73, 0x68,
	0x61, 0x63, 0x6b, 0x2f, 0x64, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x2d,
	0x63, 0x68, 0x61, 0x6f, 0x73, 0x2d, 0x35, 0x30, 0x2e, 0x30, 0x34, 0x31, 0x2f, 0x6e, 0x6f, 0x64,
	0x65, 0x2f, 0x65, 0x78, 0x70, 0x6f, 0x73, 0x65, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_external_proto_rawDescOnce sync.Once
	file_proto_external_proto_rawDescData = file_proto_external_proto_rawDesc
)

func file_proto_external_proto_rawDescGZIP() []byte {
	file_proto_external_proto_rawDescOnce.Do(func() {
		file_proto_external_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_external_proto_rawDescData)
	})
	return file_proto_external_proto_rawDescData
}

var file_proto_external_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_external_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_external_proto_goTypes = []interface{}{
	(Command)(0),             // 0: grpc.Command
	(*NewFileRequest)(nil),   // 1: grpc.NewFileRequest
	(*CheckFileRequest)(nil), // 2: grpc.CheckFileRequest
	(*Response)(nil),         // 3: grpc.Response
	(*NodeInfo)(nil),         // 4: grpc.NodeInfo
}
var file_proto_external_proto_depIdxs = []int32{
	0, // 0: grpc.NewFileRequest.command:type_name -> grpc.Command
	0, // 1: grpc.CheckFileRequest.command:type_name -> grpc.Command
	4, // 2: grpc.Response.nodeInfo:type_name -> grpc.NodeInfo
	1, // 3: grpc.ExternalListener.FileRequest:input_type -> grpc.NewFileRequest
	2, // 4: grpc.ExternalListener.CheckFile:input_type -> grpc.CheckFileRequest
	3, // 5: grpc.ExternalListener.FileRequest:output_type -> grpc.Response
	3, // 6: grpc.ExternalListener.CheckFile:output_type -> grpc.Response
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_external_proto_init() }
func file_proto_external_proto_init() {
	if File_proto_external_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_external_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewFileRequest); i {
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
		file_proto_external_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckFileRequest); i {
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
		file_proto_external_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_proto_external_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeInfo); i {
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
			RawDescriptor: file_proto_external_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_external_proto_goTypes,
		DependencyIndexes: file_proto_external_proto_depIdxs,
		EnumInfos:         file_proto_external_proto_enumTypes,
		MessageInfos:      file_proto_external_proto_msgTypes,
	}.Build()
	File_proto_external_proto = out.File
	file_proto_external_proto_rawDesc = nil
	file_proto_external_proto_goTypes = nil
	file_proto_external_proto_depIdxs = nil
}

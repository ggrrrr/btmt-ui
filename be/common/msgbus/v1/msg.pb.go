// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: be/common/msgbus/v1/msg.proto

// Message bus

package msgbusv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MessageType int32

const (
	MessageType_MESSAGE_TYPE_UNSPECIFIED MessageType = 0
	MessageType_MESSAGE_TYPE_EVENT       MessageType = 1
	MessageType_MESSAGE_TYPE_COMMAND     MessageType = 2
)

// Enum value maps for MessageType.
var (
	MessageType_name = map[int32]string{
		0: "MESSAGE_TYPE_UNSPECIFIED",
		1: "MESSAGE_TYPE_EVENT",
		2: "MESSAGE_TYPE_COMMAND",
	}
	MessageType_value = map[string]int32{
		"MESSAGE_TYPE_UNSPECIFIED": 0,
		"MESSAGE_TYPE_EVENT":       1,
		"MESSAGE_TYPE_COMMAND":     2,
	}
)

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}

func (x MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_be_common_msgbus_v1_msg_proto_enumTypes[0].Descriptor()
}

func (MessageType) Type() protoreflect.EnumType {
	return &file_be_common_msgbus_v1_msg_proto_enumTypes[0]
}

func (x MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageType.Descriptor instead.
func (MessageType) EnumDescriptor() ([]byte, []int) {
	return file_be_common_msgbus_v1_msg_proto_rawDescGZIP(), []int{0}
}

type Metadata struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          map[string]string      `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Metadata) Reset() {
	*x = Metadata{}
	mi := &file_be_common_msgbus_v1_msg_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metadata) ProtoMessage() {}

func (x *Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_be_common_msgbus_v1_msg_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metadata.ProtoReflect.Descriptor instead.
func (*Metadata) Descriptor() ([]byte, []int) {
	return file_be_common_msgbus_v1_msg_proto_rawDescGZIP(), []int{0}
}

func (x *Metadata) GetData() map[string]string {
	if x != nil {
		return x.Data
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            []byte                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	MessageType   MessageType            `protobuf:"varint,2,opt,name=message_type,json=messageType,proto3,enum=msgbus.v1.MessageType" json:"message_type,omitempty"`
	Domain        string                 `protobuf:"bytes,3,opt,name=domain,proto3" json:"domain,omitempty"`
	Name          string                 `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Md            *Metadata              `protobuf:"bytes,5,opt,name=md,proto3" json:"md,omitempty"`
	Payload       *anypb.Any             `protobuf:"bytes,6,opt,name=payload,proto3" json:"payload,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_be_common_msgbus_v1_msg_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_be_common_msgbus_v1_msg_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_be_common_msgbus_v1_msg_proto_rawDescGZIP(), []int{1}
}

func (x *Message) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Message) GetMessageType() MessageType {
	if x != nil {
		return x.MessageType
	}
	return MessageType_MESSAGE_TYPE_UNSPECIFIED
}

func (x *Message) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *Message) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Message) GetMd() *Metadata {
	if x != nil {
		return x.Md
	}
	return nil
}

func (x *Message) GetPayload() *anypb.Any {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *Message) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

var File_be_common_msgbus_v1_msg_proto protoreflect.FileDescriptor

var file_be_common_msgbus_v1_msg_proto_rawDesc = string([]byte{
	0x0a, 0x1d, 0x62, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x73, 0x67, 0x62,
	0x75, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x09, 0x6d, 0x73, 0x67, 0x62, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x76, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x31, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1d, 0x2e, 0x6d, 0x73, 0x67, 0x62, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x37, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x90,
	0x02, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x0c, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x16, 0x2e, 0x6d, 0x73, 0x67, 0x62, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x23, 0x0a, 0x02, 0x6d, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x6d, 0x73, 0x67, 0x62, 0x75, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x52, 0x02, 0x6d, 0x64, 0x12, 0x2e, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x2a, 0x5d, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x1c, 0x0a, 0x18, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16,
	0x0a, 0x12, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x45,
	0x56, 0x45, 0x4e, 0x54, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47,
	0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x10, 0x02,
	0x42, 0x1b, 0x5a, 0x19, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x73, 0x67, 0x62, 0x75,
	0x73, 0x2f, 0x76, 0x31, 0x3b, 0x6d, 0x73, 0x67, 0x62, 0x75, 0x73, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_be_common_msgbus_v1_msg_proto_rawDescOnce sync.Once
	file_be_common_msgbus_v1_msg_proto_rawDescData []byte
)

func file_be_common_msgbus_v1_msg_proto_rawDescGZIP() []byte {
	file_be_common_msgbus_v1_msg_proto_rawDescOnce.Do(func() {
		file_be_common_msgbus_v1_msg_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_be_common_msgbus_v1_msg_proto_rawDesc), len(file_be_common_msgbus_v1_msg_proto_rawDesc)))
	})
	return file_be_common_msgbus_v1_msg_proto_rawDescData
}

var file_be_common_msgbus_v1_msg_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_be_common_msgbus_v1_msg_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_be_common_msgbus_v1_msg_proto_goTypes = []any{
	(MessageType)(0),              // 0: msgbus.v1.MessageType
	(*Metadata)(nil),              // 1: msgbus.v1.Metadata
	(*Message)(nil),               // 2: msgbus.v1.Message
	nil,                           // 3: msgbus.v1.Metadata.DataEntry
	(*anypb.Any)(nil),             // 4: google.protobuf.Any
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_be_common_msgbus_v1_msg_proto_depIdxs = []int32{
	3, // 0: msgbus.v1.Metadata.data:type_name -> msgbus.v1.Metadata.DataEntry
	0, // 1: msgbus.v1.Message.message_type:type_name -> msgbus.v1.MessageType
	1, // 2: msgbus.v1.Message.md:type_name -> msgbus.v1.Metadata
	4, // 3: msgbus.v1.Message.payload:type_name -> google.protobuf.Any
	5, // 4: msgbus.v1.Message.created_at:type_name -> google.protobuf.Timestamp
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_be_common_msgbus_v1_msg_proto_init() }
func file_be_common_msgbus_v1_msg_proto_init() {
	if File_be_common_msgbus_v1_msg_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_be_common_msgbus_v1_msg_proto_rawDesc), len(file_be_common_msgbus_v1_msg_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_be_common_msgbus_v1_msg_proto_goTypes,
		DependencyIndexes: file_be_common_msgbus_v1_msg_proto_depIdxs,
		EnumInfos:         file_be_common_msgbus_v1_msg_proto_enumTypes,
		MessageInfos:      file_be_common_msgbus_v1_msg_proto_msgTypes,
	}.Build()
	File_be_common_msgbus_v1_msg_proto = out.File
	file_be_common_msgbus_v1_msg_proto_goTypes = nil
	file_be_common_msgbus_v1_msg_proto_depIdxs = nil
}

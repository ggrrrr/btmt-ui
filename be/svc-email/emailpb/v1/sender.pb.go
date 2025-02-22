// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: be/svc-email/emailpb/v1/sender.proto

// Email sender service

package emailpbv1

import (
	v1 "github.com/ggrrrr/btmt-ui/be/common/templ/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type RawBody struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ContentType   string                 `protobuf:"bytes,1,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
	Subject       string                 `protobuf:"bytes,2,opt,name=subject,proto3" json:"subject,omitempty"`
	Body          string                 `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawBody) Reset() {
	*x = RawBody{}
	mi := &file_be_svc_email_emailpb_v1_sender_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawBody) ProtoMessage() {}

func (x *RawBody) ProtoReflect() protoreflect.Message {
	mi := &file_be_svc_email_emailpb_v1_sender_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawBody.ProtoReflect.Descriptor instead.
func (*RawBody) Descriptor() ([]byte, []int) {
	return file_be_svc_email_emailpb_v1_sender_proto_rawDescGZIP(), []int{0}
}

func (x *RawBody) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

func (x *RawBody) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *RawBody) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type EmailAddr struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email         string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmailAddr) Reset() {
	*x = EmailAddr{}
	mi := &file_be_svc_email_emailpb_v1_sender_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmailAddr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailAddr) ProtoMessage() {}

func (x *EmailAddr) ProtoReflect() protoreflect.Message {
	mi := &file_be_svc_email_emailpb_v1_sender_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailAddr.ProtoReflect.Descriptor instead.
func (*EmailAddr) Descriptor() ([]byte, []int) {
	return file_be_svc_email_emailpb_v1_sender_proto_rawDescGZIP(), []int{1}
}

func (x *EmailAddr) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EmailAddr) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type SenderAccount struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Realm         string                 `protobuf:"bytes,1,opt,name=realm,proto3" json:"realm,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email         string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SenderAccount) Reset() {
	*x = SenderAccount{}
	mi := &file_be_svc_email_emailpb_v1_sender_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SenderAccount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SenderAccount) ProtoMessage() {}

func (x *SenderAccount) ProtoReflect() protoreflect.Message {
	mi := &file_be_svc_email_emailpb_v1_sender_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SenderAccount.ProtoReflect.Descriptor instead.
func (*SenderAccount) Descriptor() ([]byte, []int) {
	return file_be_svc_email_emailpb_v1_sender_proto_rawDescGZIP(), []int{2}
}

func (x *SenderAccount) GetRealm() string {
	if x != nil {
		return x.Realm
	}
	return ""
}

func (x *SenderAccount) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SenderAccount) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type ToAddresses struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ToEmail       []*EmailAddr           `protobuf:"bytes,1,rep,name=to_email,json=toEmail,proto3" json:"to_email,omitempty"`
	CcEmail       []*EmailAddr           `protobuf:"bytes,2,rep,name=cc_email,json=ccEmail,proto3" json:"cc_email,omitempty"`
	BccEmail      []*EmailAddr           `protobuf:"bytes,3,rep,name=bcc_email,json=bccEmail,proto3" json:"bcc_email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ToAddresses) Reset() {
	*x = ToAddresses{}
	mi := &file_be_svc_email_emailpb_v1_sender_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ToAddresses) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToAddresses) ProtoMessage() {}

func (x *ToAddresses) ProtoReflect() protoreflect.Message {
	mi := &file_be_svc_email_emailpb_v1_sender_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToAddresses.ProtoReflect.Descriptor instead.
func (*ToAddresses) Descriptor() ([]byte, []int) {
	return file_be_svc_email_emailpb_v1_sender_proto_rawDescGZIP(), []int{3}
}

func (x *ToAddresses) GetToEmail() []*EmailAddr {
	if x != nil {
		return x.ToEmail
	}
	return nil
}

func (x *ToAddresses) GetCcEmail() []*EmailAddr {
	if x != nil {
		return x.CcEmail
	}
	return nil
}

func (x *ToAddresses) GetBccEmail() []*EmailAddr {
	if x != nil {
		return x.BccEmail
	}
	return nil
}

type EmailMessage struct {
	state       protoimpl.MessageState `protogen:"open.v1"`
	FromAccount *SenderAccount         `protobuf:"bytes,1,opt,name=from_account,proto3" json:"from_account,omitempty"`
	ToAddresses *ToAddresses           `protobuf:"bytes,2,opt,name=to_addresses,proto3" json:"to_addresses,omitempty"`
	// Types that are valid to be assigned to Body:
	//
	//	*EmailMessage_RawBody
	//	*EmailMessage_TemplateId
	Body          isEmailMessage_Body `protobuf_oneof:"body"`
	Data          *v1.Data            `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmailMessage) Reset() {
	*x = EmailMessage{}
	mi := &file_be_svc_email_emailpb_v1_sender_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmailMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailMessage) ProtoMessage() {}

func (x *EmailMessage) ProtoReflect() protoreflect.Message {
	mi := &file_be_svc_email_emailpb_v1_sender_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailMessage.ProtoReflect.Descriptor instead.
func (*EmailMessage) Descriptor() ([]byte, []int) {
	return file_be_svc_email_emailpb_v1_sender_proto_rawDescGZIP(), []int{4}
}

func (x *EmailMessage) GetFromAccount() *SenderAccount {
	if x != nil {
		return x.FromAccount
	}
	return nil
}

func (x *EmailMessage) GetToAddresses() *ToAddresses {
	if x != nil {
		return x.ToAddresses
	}
	return nil
}

func (x *EmailMessage) GetBody() isEmailMessage_Body {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *EmailMessage) GetRawBody() *RawBody {
	if x != nil {
		if x, ok := x.Body.(*EmailMessage_RawBody); ok {
			return x.RawBody
		}
	}
	return nil
}

func (x *EmailMessage) GetTemplateId() string {
	if x != nil {
		if x, ok := x.Body.(*EmailMessage_TemplateId); ok {
			return x.TemplateId
		}
	}
	return ""
}

func (x *EmailMessage) GetData() *v1.Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type isEmailMessage_Body interface {
	isEmailMessage_Body()
}

type EmailMessage_RawBody struct {
	RawBody *RawBody `protobuf:"bytes,3,opt,name=raw_body,proto3,oneof"`
}

type EmailMessage_TemplateId struct {
	TemplateId string `protobuf:"bytes,4,opt,name=template_id,proto3,oneof"`
}

func (*EmailMessage_RawBody) isEmailMessage_Body() {}

func (*EmailMessage_TemplateId) isEmailMessage_Body() {}

type SendEmail struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Message       *EmailMessage          `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendEmail) Reset() {
	*x = SendEmail{}
	mi := &file_be_svc_email_emailpb_v1_sender_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendEmail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendEmail) ProtoMessage() {}

func (x *SendEmail) ProtoReflect() protoreflect.Message {
	mi := &file_be_svc_email_emailpb_v1_sender_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendEmail.ProtoReflect.Descriptor instead.
func (*SendEmail) Descriptor() ([]byte, []int) {
	return file_be_svc_email_emailpb_v1_sender_proto_rawDescGZIP(), []int{5}
}

func (x *SendEmail) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SendEmail) GetMessage() *EmailMessage {
	if x != nil {
		return x.Message
	}
	return nil
}

var File_be_svc_email_emailpb_v1_sender_proto protoreflect.FileDescriptor

var file_be_svc_email_emailpb_v1_sender_proto_rawDesc = string([]byte{
	0x0a, 0x24, 0x62, 0x65, 0x2f, 0x73, 0x76, 0x63, 0x2d, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x2f, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x70, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x70, 0x62, 0x2e,
	0x76, 0x31, 0x1a, 0x1e, 0x62, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x74, 0x65,
	0x6d, 0x70, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x5a, 0x0a, 0x07, 0x52, 0x61, 0x77, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x21, 0x0a,
	0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f,
	0x64, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x35,
	0x0a, 0x09, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x41, 0x64, 0x64, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x4f, 0x0a, 0x0d, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0xa5, 0x01, 0x0a, 0x0b, 0x54, 0x6f, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x12, 0x30, 0x0a, 0x08, 0x74, 0x6f, 0x5f, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x41, 0x64, 0x64, 0x72, 0x52,
	0x07, 0x74, 0x6f, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x30, 0x0a, 0x08, 0x63, 0x63, 0x5f, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x41, 0x64, 0x64,
	0x72, 0x52, 0x07, 0x63, 0x63, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x32, 0x0a, 0x09, 0x62, 0x63,
	0x63, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x41, 0x64, 0x64, 0x72, 0x52, 0x08, 0x62, 0x63, 0x63, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x8d,
	0x02, 0x0a, 0x0c, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x3d, 0x0a, 0x0c, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x70, 0x62, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x52, 0x0c, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x3b,
	0x0a, 0x0c, 0x74, 0x6f, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x70, 0x62, 0x2e, 0x76,
	0x31, 0x2e, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x52, 0x0c, 0x74,
	0x6f, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x12, 0x31, 0x0a, 0x08, 0x72,
	0x61, 0x77, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x77, 0x42, 0x6f,
	0x64, 0x79, 0x48, 0x00, 0x52, 0x08, 0x72, 0x61, 0x77, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x22,
	0x0a, 0x0b, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f,
	0x69, 0x64, 0x12, 0x22, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x42, 0x06, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x4f,
	0x0a, 0x09, 0x53, 0x65, 0x6e, 0x64, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x32, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42,
	0x20, 0x5a, 0x1e, 0x73, 0x76, 0x63, 0x2d, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x2f, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x70, 0x62, 0x2f, 0x76, 0x31, 0x3b, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x70, 0x62, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_be_svc_email_emailpb_v1_sender_proto_rawDescOnce sync.Once
	file_be_svc_email_emailpb_v1_sender_proto_rawDescData []byte
)

func file_be_svc_email_emailpb_v1_sender_proto_rawDescGZIP() []byte {
	file_be_svc_email_emailpb_v1_sender_proto_rawDescOnce.Do(func() {
		file_be_svc_email_emailpb_v1_sender_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_be_svc_email_emailpb_v1_sender_proto_rawDesc), len(file_be_svc_email_emailpb_v1_sender_proto_rawDesc)))
	})
	return file_be_svc_email_emailpb_v1_sender_proto_rawDescData
}

var file_be_svc_email_emailpb_v1_sender_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_be_svc_email_emailpb_v1_sender_proto_goTypes = []any{
	(*RawBody)(nil),       // 0: emailpb.v1.RawBody
	(*EmailAddr)(nil),     // 1: emailpb.v1.EmailAddr
	(*SenderAccount)(nil), // 2: emailpb.v1.SenderAccount
	(*ToAddresses)(nil),   // 3: emailpb.v1.ToAddresses
	(*EmailMessage)(nil),  // 4: emailpb.v1.EmailMessage
	(*SendEmail)(nil),     // 5: emailpb.v1.SendEmail
	(*v1.Data)(nil),       // 6: templ.v1.Data
}
var file_be_svc_email_emailpb_v1_sender_proto_depIdxs = []int32{
	1, // 0: emailpb.v1.ToAddresses.to_email:type_name -> emailpb.v1.EmailAddr
	1, // 1: emailpb.v1.ToAddresses.cc_email:type_name -> emailpb.v1.EmailAddr
	1, // 2: emailpb.v1.ToAddresses.bcc_email:type_name -> emailpb.v1.EmailAddr
	2, // 3: emailpb.v1.EmailMessage.from_account:type_name -> emailpb.v1.SenderAccount
	3, // 4: emailpb.v1.EmailMessage.to_addresses:type_name -> emailpb.v1.ToAddresses
	0, // 5: emailpb.v1.EmailMessage.raw_body:type_name -> emailpb.v1.RawBody
	6, // 6: emailpb.v1.EmailMessage.data:type_name -> templ.v1.Data
	4, // 7: emailpb.v1.SendEmail.message:type_name -> emailpb.v1.EmailMessage
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_be_svc_email_emailpb_v1_sender_proto_init() }
func file_be_svc_email_emailpb_v1_sender_proto_init() {
	if File_be_svc_email_emailpb_v1_sender_proto != nil {
		return
	}
	file_be_svc_email_emailpb_v1_sender_proto_msgTypes[4].OneofWrappers = []any{
		(*EmailMessage_RawBody)(nil),
		(*EmailMessage_TemplateId)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_be_svc_email_emailpb_v1_sender_proto_rawDesc), len(file_be_svc_email_emailpb_v1_sender_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_be_svc_email_emailpb_v1_sender_proto_goTypes,
		DependencyIndexes: file_be_svc_email_emailpb_v1_sender_proto_depIdxs,
		MessageInfos:      file_be_svc_email_emailpb_v1_sender_proto_msgTypes,
	}.Build()
	File_be_svc_email_emailpb_v1_sender_proto = out.File
	file_be_svc_email_emailpb_v1_sender_proto_goTypes = nil
	file_be_svc_email_emailpb_v1_sender_proto_depIdxs = nil
}

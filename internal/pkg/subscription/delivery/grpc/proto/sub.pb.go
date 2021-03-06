// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.0
// source: sub.proto

// protoc --go_out=plugins=grpc:. *.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type UserId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *UserId) Reset() {
	*x = UserId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sub_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserId) ProtoMessage() {}

func (x *UserId) ProtoReflect() protoreflect.Message {
	mi := &file_sub_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserId.ProtoReflect.Descriptor instead.
func (*UserId) Descriptor() ([]byte, []int) {
	return file_sub_proto_rawDescGZIP(), []int{0}
}

func (x *UserId) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

type Payment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type        string `protobuf:"bytes,1,opt,name=Type,proto3" json:"Type,omitempty"`
	OperationID string `protobuf:"bytes,2,opt,name=OperationID,proto3" json:"OperationID,omitempty"`
	Amount      string `protobuf:"bytes,3,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Currency    string `protobuf:"bytes,4,opt,name=Currency,proto3" json:"Currency,omitempty"`
	DateTime    string `protobuf:"bytes,5,opt,name=DateTime,proto3" json:"DateTime,omitempty"`
	Sender      string `protobuf:"bytes,6,opt,name=Sender,proto3" json:"Sender,omitempty"`
	CodePro     bool   `protobuf:"varint,7,opt,name=CodePro,proto3" json:"CodePro,omitempty"`
	Label       string `protobuf:"bytes,8,opt,name=Label,proto3" json:"Label,omitempty"`
	Unaccepted  bool   `protobuf:"varint,9,opt,name=Unaccepted,proto3" json:"Unaccepted,omitempty"`
	Hash        string `protobuf:"bytes,10,opt,name=Hash,proto3" json:"Hash,omitempty"`
}

func (x *Payment) Reset() {
	*x = Payment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sub_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Payment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Payment) ProtoMessage() {}

func (x *Payment) ProtoReflect() protoreflect.Message {
	mi := &file_sub_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Payment.ProtoReflect.Descriptor instead.
func (*Payment) Descriptor() ([]byte, []int) {
	return file_sub_proto_rawDescGZIP(), []int{1}
}

func (x *Payment) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Payment) GetOperationID() string {
	if x != nil {
		return x.OperationID
	}
	return ""
}

func (x *Payment) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *Payment) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *Payment) GetDateTime() string {
	if x != nil {
		return x.DateTime
	}
	return ""
}

func (x *Payment) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *Payment) GetCodePro() bool {
	if x != nil {
		return x.CodePro
	}
	return false
}

func (x *Payment) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *Payment) GetUnaccepted() bool {
	if x != nil {
		return x.Unaccepted
	}
	return false
}

func (x *Payment) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

type ErrorMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error int32 `protobuf:"varint,1,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *ErrorMessage) Reset() {
	*x = ErrorMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sub_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrorMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorMessage) ProtoMessage() {}

func (x *ErrorMessage) ProtoReflect() protoreflect.Message {
	mi := &file_sub_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorMessage.ProtoReflect.Descriptor instead.
func (*ErrorMessage) Descriptor() ([]byte, []int) {
	return file_sub_proto_rawDescGZIP(), []int{2}
}

func (x *ErrorMessage) GetError() int32 {
	if x != nil {
		return x.Error
	}
	return 0
}

var File_sub_proto protoreflect.FileDescriptor

var file_sub_proto_rawDesc = []byte{
	0x0a, 0x09, 0x73, 0x75, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x18, 0x0a, 0x06, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x49, 0x44, 0x22, 0x8b, 0x02, 0x0a, 0x07, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1a, 0x0a,
	0x08, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x44, 0x61, 0x74,
	0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x44, 0x61, 0x74,
	0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x18, 0x0a,
	0x07, 0x43, 0x6f, 0x64, 0x65, 0x50, 0x72, 0x6f, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x43, 0x6f, 0x64, 0x65, 0x50, 0x72, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x1e, 0x0a,
	0x0a, 0x55, 0x6e, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0a, 0x55, 0x6e, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x48, 0x61, 0x73,
	0x68, 0x22, 0x24, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x8b, 0x01, 0x0a, 0x0c, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3d, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x12, 0x15, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x1a, 0x2e, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x12, 0x14, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x1a, 0x1a, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sub_proto_rawDescOnce sync.Once
	file_sub_proto_rawDescData = file_sub_proto_rawDesc
)

func file_sub_proto_rawDescGZIP() []byte {
	file_sub_proto_rawDescOnce.Do(func() {
		file_sub_proto_rawDescData = protoimpl.X.CompressGZIP(file_sub_proto_rawDescData)
	})
	return file_sub_proto_rawDescData
}

var file_sub_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_sub_proto_goTypes = []interface{}{
	(*UserId)(nil),       // 0: Subscription.UserId
	(*Payment)(nil),      // 1: Subscription.Payment
	(*ErrorMessage)(nil), // 2: Subscription.ErrorMessage
}
var file_sub_proto_depIdxs = []int32{
	1, // 0: Subscription.Subscription.Create:input_type -> Subscription.Payment
	0, // 1: Subscription.Subscription.Delete:input_type -> Subscription.UserId
	2, // 2: Subscription.Subscription.Create:output_type -> Subscription.ErrorMessage
	2, // 3: Subscription.Subscription.Delete:output_type -> Subscription.ErrorMessage
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_sub_proto_init() }
func file_sub_proto_init() {
	if File_sub_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sub_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserId); i {
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
		file_sub_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Payment); i {
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
		file_sub_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrorMessage); i {
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
			RawDescriptor: file_sub_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sub_proto_goTypes,
		DependencyIndexes: file_sub_proto_depIdxs,
		MessageInfos:      file_sub_proto_msgTypes,
	}.Build()
	File_sub_proto = out.File
	file_sub_proto_rawDesc = nil
	file_sub_proto_goTypes = nil
	file_sub_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SubscriptionClient is the client API for Subscription service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SubscriptionClient interface {
	Create(ctx context.Context, in *Payment, opts ...grpc.CallOption) (*ErrorMessage, error)
	Delete(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*ErrorMessage, error)
}

type subscriptionClient struct {
	cc grpc.ClientConnInterface
}

func NewSubscriptionClient(cc grpc.ClientConnInterface) SubscriptionClient {
	return &subscriptionClient{cc}
}

func (c *subscriptionClient) Create(ctx context.Context, in *Payment, opts ...grpc.CallOption) (*ErrorMessage, error) {
	out := new(ErrorMessage)
	err := c.cc.Invoke(ctx, "/Subscription.Subscription/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscriptionClient) Delete(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*ErrorMessage, error) {
	out := new(ErrorMessage)
	err := c.cc.Invoke(ctx, "/Subscription.Subscription/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SubscriptionServer is the server API for Subscription service.
type SubscriptionServer interface {
	Create(context.Context, *Payment) (*ErrorMessage, error)
	Delete(context.Context, *UserId) (*ErrorMessage, error)
}

// UnimplementedSubscriptionServer can be embedded to have forward compatible implementations.
type UnimplementedSubscriptionServer struct {
}

func (*UnimplementedSubscriptionServer) Create(context.Context, *Payment) (*ErrorMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedSubscriptionServer) Delete(context.Context, *UserId) (*ErrorMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func RegisterSubscriptionServer(s *grpc.Server, srv SubscriptionServer) {
	s.RegisterService(&_Subscription_serviceDesc, srv)
}

func _Subscription_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Payment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscriptionServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Subscription.Subscription/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscriptionServer).Create(ctx, req.(*Payment))
	}
	return interceptor(ctx, in, info, handler)
}

func _Subscription_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscriptionServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Subscription.Subscription/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscriptionServer).Delete(ctx, req.(*UserId))
	}
	return interceptor(ctx, in, info, handler)
}

var _Subscription_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Subscription.Subscription",
	HandlerType: (*SubscriptionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Subscription_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Subscription_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sub.proto",
}

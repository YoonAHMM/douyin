// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.19.4
// source: jwt.proto

package Jwt

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

type CreateTokenReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID       string `protobuf:"bytes,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	AccessExpire int64  `protobuf:"varint,2,opt,name=AccessExpire,proto3" json:"AccessExpire,omitempty"`
}

func (x *CreateTokenReq) Reset() {
	*x = CreateTokenReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jwt_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTokenReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTokenReq) ProtoMessage() {}

func (x *CreateTokenReq) ProtoReflect() protoreflect.Message {
	mi := &file_jwt_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTokenReq.ProtoReflect.Descriptor instead.
func (*CreateTokenReq) Descriptor() ([]byte, []int) {
	return file_jwt_proto_rawDescGZIP(), []int{0}
}

func (x *CreateTokenReq) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *CreateTokenReq) GetAccessExpire() int64 {
	if x != nil {
		return x.AccessExpire
	}
	return 0
}

type CreateTokenResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=Token,proto3" json:"Token,omitempty"`
}

func (x *CreateTokenResp) Reset() {
	*x = CreateTokenResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jwt_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTokenResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTokenResp) ProtoMessage() {}

func (x *CreateTokenResp) ProtoReflect() protoreflect.Message {
	mi := &file_jwt_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTokenResp.ProtoReflect.Descriptor instead.
func (*CreateTokenResp) Descriptor() ([]byte, []int) {
	return file_jwt_proto_rawDescGZIP(), []int{1}
}

func (x *CreateTokenResp) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type ParseTokenReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *ParseTokenReq) Reset() {
	*x = ParseTokenReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jwt_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ParseTokenReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParseTokenReq) ProtoMessage() {}

func (x *ParseTokenReq) ProtoReflect() protoreflect.Message {
	mi := &file_jwt_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParseTokenReq.ProtoReflect.Descriptor instead.
func (*ParseTokenReq) Descriptor() ([]byte, []int) {
	return file_jwt_proto_rawDescGZIP(), []int{2}
}

func (x *ParseTokenReq) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type ParseTokenResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID       string `protobuf:"bytes,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	AccessExpire int64  `protobuf:"varint,2,opt,name=AccessExpire,proto3" json:"AccessExpire,omitempty"`
}

func (x *ParseTokenResp) Reset() {
	*x = ParseTokenResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jwt_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ParseTokenResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParseTokenResp) ProtoMessage() {}

func (x *ParseTokenResp) ProtoReflect() protoreflect.Message {
	mi := &file_jwt_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParseTokenResp.ProtoReflect.Descriptor instead.
func (*ParseTokenResp) Descriptor() ([]byte, []int) {
	return file_jwt_proto_rawDescGZIP(), []int{3}
}

func (x *ParseTokenResp) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *ParseTokenResp) GetAccessExpire() int64 {
	if x != nil {
		return x.AccessExpire
	}
	return 0
}

type IsValidTokenReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *IsValidTokenReq) Reset() {
	*x = IsValidTokenReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jwt_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsValidTokenReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsValidTokenReq) ProtoMessage() {}

func (x *IsValidTokenReq) ProtoReflect() protoreflect.Message {
	mi := &file_jwt_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsValidTokenReq.ProtoReflect.Descriptor instead.
func (*IsValidTokenReq) Descriptor() ([]byte, []int) {
	return file_jwt_proto_rawDescGZIP(), []int{4}
}

func (x *IsValidTokenReq) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type IsValidTokenResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Isvaild bool `protobuf:"varint,1,opt,name=isvaild,proto3" json:"isvaild,omitempty"`
}

func (x *IsValidTokenResp) Reset() {
	*x = IsValidTokenResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_jwt_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsValidTokenResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsValidTokenResp) ProtoMessage() {}

func (x *IsValidTokenResp) ProtoReflect() protoreflect.Message {
	mi := &file_jwt_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsValidTokenResp.ProtoReflect.Descriptor instead.
func (*IsValidTokenResp) Descriptor() ([]byte, []int) {
	return file_jwt_proto_rawDescGZIP(), []int{5}
}

func (x *IsValidTokenResp) GetIsvaild() bool {
	if x != nil {
		return x.Isvaild
	}
	return false
}

var File_jwt_proto protoreflect.FileDescriptor

var file_jwt_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6a, 0x77, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x6a, 0x77, 0x74,
	0x22, 0x4c, 0x0a, 0x0e, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x22, 0x0a, 0x0c, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0c, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x22, 0x27,
	0x0a, 0x0f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x25, 0x0a, 0x0d, 0x70, 0x61, 0x72, 0x73, 0x65,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x4c,
	0x0a, 0x0e, 0x70, 0x61, 0x72, 0x73, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x22, 0x0a, 0x0c, 0x41, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c,
	0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x22, 0x27, 0x0a, 0x0f,
	0x69, 0x73, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x2c, 0x0a, 0x10, 0x69, 0x73, 0x56, 0x61, 0x6c, 0x69, 0x64,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x69, 0x73, 0x76,
	0x61, 0x69, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x76, 0x61,
	0x69, 0x6c, 0x64, 0x32, 0xbc, 0x01, 0x0a, 0x06, 0x4a, 0x77, 0x74, 0x52, 0x70, 0x63, 0x12, 0x3a,
	0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x13, 0x2e,
	0x6a, 0x77, 0x74, 0x2e, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x65, 0x71, 0x1a, 0x14, 0x2e, 0x6a, 0x77, 0x74, 0x2e, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x0a, 0x70, 0x61,
	0x72, 0x73, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x12, 0x2e, 0x6a, 0x77, 0x74, 0x2e, 0x70,
	0x61, 0x72, 0x73, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x13, 0x2e, 0x6a,
	0x77, 0x74, 0x2e, 0x70, 0x61, 0x72, 0x73, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x0c, 0x49, 0x73, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x14, 0x2e, 0x6a, 0x77, 0x74, 0x2e, 0x69, 0x73, 0x56, 0x61, 0x6c, 0x69,
	0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x15, 0x2e, 0x6a, 0x77, 0x74, 0x2e,
	0x69, 0x73, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x22, 0x00, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x4a, 0x77, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_jwt_proto_rawDescOnce sync.Once
	file_jwt_proto_rawDescData = file_jwt_proto_rawDesc
)

func file_jwt_proto_rawDescGZIP() []byte {
	file_jwt_proto_rawDescOnce.Do(func() {
		file_jwt_proto_rawDescData = protoimpl.X.CompressGZIP(file_jwt_proto_rawDescData)
	})
	return file_jwt_proto_rawDescData
}

var file_jwt_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_jwt_proto_goTypes = []interface{}{
	(*CreateTokenReq)(nil),   // 0: jwt.createTokenReq
	(*CreateTokenResp)(nil),  // 1: jwt.createTokenResp
	(*ParseTokenReq)(nil),    // 2: jwt.parseTokenReq
	(*ParseTokenResp)(nil),   // 3: jwt.parseTokenResp
	(*IsValidTokenReq)(nil),  // 4: jwt.isValidTokenReq
	(*IsValidTokenResp)(nil), // 5: jwt.isValidTokenResp
}
var file_jwt_proto_depIdxs = []int32{
	0, // 0: jwt.JwtRpc.createToken:input_type -> jwt.createTokenReq
	2, // 1: jwt.JwtRpc.parseToken:input_type -> jwt.parseTokenReq
	4, // 2: jwt.JwtRpc.IsValidToken:input_type -> jwt.isValidTokenReq
	1, // 3: jwt.JwtRpc.createToken:output_type -> jwt.createTokenResp
	3, // 4: jwt.JwtRpc.parseToken:output_type -> jwt.parseTokenResp
	5, // 5: jwt.JwtRpc.IsValidToken:output_type -> jwt.isValidTokenResp
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_jwt_proto_init() }
func file_jwt_proto_init() {
	if File_jwt_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_jwt_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTokenReq); i {
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
		file_jwt_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTokenResp); i {
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
		file_jwt_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ParseTokenReq); i {
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
		file_jwt_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ParseTokenResp); i {
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
		file_jwt_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsValidTokenReq); i {
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
		file_jwt_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsValidTokenResp); i {
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
			RawDescriptor: file_jwt_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_jwt_proto_goTypes,
		DependencyIndexes: file_jwt_proto_depIdxs,
		MessageInfos:      file_jwt_proto_msgTypes,
	}.Build()
	File_jwt_proto = out.File
	file_jwt_proto_rawDesc = nil
	file_jwt_proto_goTypes = nil
	file_jwt_proto_depIdxs = nil
}

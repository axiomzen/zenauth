// Code generated by protoc-gen-go. DO NOT EDIT.
// source: auth.proto

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	auth.proto
	user.proto

It has these top-level messages:
	UserID
	InvitationCode
	User
	UserPublic
*/
package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type UserID struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *UserID) Reset()                    { *m = UserID{} }
func (m *UserID) String() string            { return proto.CompactTextString(m) }
func (*UserID) ProtoMessage()               {}
func (*UserID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *UserID) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type InvitationCode struct {
	Type       string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	InviteCode string `protobuf:"bytes,2,opt,name=inviteCode" json:"inviteCode,omitempty"`
}

func (m *InvitationCode) Reset()                    { *m = InvitationCode{} }
func (m *InvitationCode) String() string            { return proto.CompactTextString(m) }
func (*InvitationCode) ProtoMessage()               {}
func (*InvitationCode) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *InvitationCode) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *InvitationCode) GetInviteCode() string {
	if m != nil {
		return m.InviteCode
	}
	return ""
}

func init() {
	proto.RegisterType((*UserID)(nil), "protobuf.UserID")
	proto.RegisterType((*InvitationCode)(nil), "protobuf.InvitationCode")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Auth service

type AuthClient interface {
	GetCurrentUser(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*User, error)
	GetUserByID(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*UserPublic, error)
	LinkUser(ctx context.Context, in *InvitationCode, opts ...grpc.CallOption) (*UserPublic, error)
}

type authClient struct {
	cc *grpc.ClientConn
}

func NewAuthClient(cc *grpc.ClientConn) AuthClient {
	return &authClient{cc}
}

func (c *authClient) GetCurrentUser(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/protobuf.Auth/GetCurrentUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) GetUserByID(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*UserPublic, error) {
	out := new(UserPublic)
	err := grpc.Invoke(ctx, "/protobuf.Auth/GetUserByID", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) LinkUser(ctx context.Context, in *InvitationCode, opts ...grpc.CallOption) (*UserPublic, error) {
	out := new(UserPublic)
	err := grpc.Invoke(ctx, "/protobuf.Auth/LinkUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Auth service

type AuthServer interface {
	GetCurrentUser(context.Context, *google_protobuf.Empty) (*User, error)
	GetUserByID(context.Context, *UserID) (*UserPublic, error)
	LinkUser(context.Context, *InvitationCode) (*UserPublic, error)
}

func RegisterAuthServer(s *grpc.Server, srv AuthServer) {
	s.RegisterService(&_Auth_serviceDesc, srv)
}

func _Auth_GetCurrentUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).GetCurrentUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Auth/GetCurrentUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).GetCurrentUser(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_GetUserByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).GetUserByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Auth/GetUserByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).GetUserByID(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_LinkUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InvitationCode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).LinkUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Auth/LinkUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).LinkUser(ctx, req.(*InvitationCode))
	}
	return interceptor(ctx, in, info, handler)
}

var _Auth_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCurrentUser",
			Handler:    _Auth_GetCurrentUser_Handler,
		},
		{
			MethodName: "GetUserByID",
			Handler:    _Auth_GetUserByID_Handler,
		},
		{
			MethodName: "LinkUser",
			Handler:    _Auth_LinkUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}

func init() { proto.RegisterFile("auth.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 239 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x8f, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0xbb, 0x4b, 0x29, 0x75, 0x84, 0x45, 0x06, 0x91, 0x65, 0x05, 0x91, 0x9c, 0x3c, 0xa5,
	0xa0, 0x07, 0x41, 0xbc, 0x68, 0x23, 0x25, 0xe0, 0x41, 0x04, 0x1f, 0xa0, 0xeb, 0x8e, 0x6d, 0xb0,
	0x26, 0x4b, 0x3a, 0x11, 0xf6, 0xd1, 0x7c, 0x3b, 0x49, 0xd6, 0x52, 0xf7, 0xd0, 0x53, 0x86, 0x2f,
	0xdf, 0xf0, 0xff, 0x03, 0xb0, 0x0c, 0xbc, 0x96, 0xad, 0x77, 0xec, 0x70, 0x9a, 0x9e, 0x3a, 0x7c,
	0x54, 0xe7, 0x2b, 0xe7, 0x56, 0x1b, 0x9a, 0xed, 0xc0, 0x8c, 0xbe, 0x5a, 0xee, 0x7a, 0xad, 0x82,
	0xb0, 0x25, 0xdf, 0xcf, 0xa2, 0x84, 0xc9, 0xdb, 0x96, 0xbc, 0x56, 0x58, 0x40, 0x6e, 0x9a, 0x32,
	0xbb, 0xcc, 0xae, 0x8e, 0x5e, 0x73, 0xd3, 0x08, 0x05, 0x85, 0xb6, 0xdf, 0x86, 0x97, 0x6c, 0x9c,
	0x9d, 0xbb, 0x86, 0x10, 0x61, 0xcc, 0x5d, 0x4b, 0x7f, 0x4e, 0x9a, 0xf1, 0x02, 0xc0, 0x44, 0x8b,
	0xa2, 0x51, 0xe6, 0xe9, 0xe7, 0x1f, 0xb9, 0xfe, 0xc9, 0x60, 0xfc, 0x10, 0x78, 0x8d, 0x77, 0x50,
	0x2c, 0x88, 0xe7, 0xc1, 0x7b, 0xb2, 0x1c, 0x23, 0xf1, 0x4c, 0xf6, 0x25, 0xe5, 0xae, 0xa4, 0x7c,
	0x8a, 0x25, 0xab, 0x62, 0x0f, 0xa2, 0x27, 0x46, 0x78, 0x0b, 0xc7, 0x0b, 0x4a, 0x4b, 0x8f, 0x9d,
	0x56, 0x78, 0x32, 0x14, 0xb4, 0xaa, 0x4e, 0x87, 0xe4, 0x25, 0xd4, 0x1b, 0xf3, 0x2e, 0x46, 0x78,
	0x0f, 0xd3, 0x67, 0x63, 0x3f, 0x53, 0x5c, 0xb9, 0x77, 0x86, 0x77, 0x1d, 0xda, 0xae, 0x27, 0x09,
	0xdf, 0xfc, 0x06, 0x00, 0x00, 0xff, 0xff, 0x84, 0x51, 0xce, 0x2d, 0x63, 0x01, 0x00, 0x00,
}

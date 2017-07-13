// Code generated by protoc-gen-go. DO NOT EDIT.
// source: user.proto

package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type UserStatus int32

const (
	UserStatus_invited UserStatus = 0
	UserStatus_created UserStatus = 1
	UserStatus_merged  UserStatus = 2
)

var UserStatus_name = map[int32]string{
	0: "invited",
	1: "created",
	2: "merged",
}
var UserStatus_value = map[string]int32{
	"invited": 0,
	"created": 1,
	"merged":  2,
}

func (x UserStatus) String() string {
	return proto.EnumName(UserStatus_name, int32(x))
}
func (UserStatus) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

type User struct {
	Id         string                      `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Email      string                      `protobuf:"bytes,2,opt,name=email" json:"email,omitempty"`
	CreatedAt  *google_protobuf1.Timestamp `protobuf:"bytes,3,opt,name=createdAt" json:"createdAt,omitempty"`
	UpdatedAt  *google_protobuf1.Timestamp `protobuf:"bytes,4,opt,name=updatedAt" json:"updatedAt,omitempty"`
	Verified   bool                        `protobuf:"varint,5,opt,name=verified" json:"verified,omitempty"`
	AuthToken  string                      `protobuf:"bytes,6,opt,name=authToken" json:"authToken,omitempty"`
	Status     UserStatus                  `protobuf:"varint,7,opt,name=status,enum=protobuf.UserStatus" json:"status,omitempty"`
	FacebookID string                      `protobuf:"bytes,8,opt,name=facebookID" json:"facebookID,omitempty"`
	UserName   string                      `protobuf:"bytes,9,opt,name=userName" json:"userName,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *User) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetCreatedAt() *google_protobuf1.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *User) GetUpdatedAt() *google_protobuf1.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

func (m *User) GetVerified() bool {
	if m != nil {
		return m.Verified
	}
	return false
}

func (m *User) GetAuthToken() string {
	if m != nil {
		return m.AuthToken
	}
	return ""
}

func (m *User) GetStatus() UserStatus {
	if m != nil {
		return m.Status
	}
	return UserStatus_invited
}

func (m *User) GetFacebookID() string {
	if m != nil {
		return m.FacebookID
	}
	return ""
}

func (m *User) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

type UserPublic struct {
	Id         string     `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Email      string     `protobuf:"bytes,2,opt,name=email" json:"email,omitempty"`
	Status     UserStatus `protobuf:"varint,3,opt,name=status,enum=protobuf.UserStatus" json:"status,omitempty"`
	FacebookID string     `protobuf:"bytes,4,opt,name=facebookID" json:"facebookID,omitempty"`
	UserName   string     `protobuf:"bytes,5,opt,name=userName" json:"userName,omitempty"`
}

func (m *UserPublic) Reset()                    { *m = UserPublic{} }
func (m *UserPublic) String() string            { return proto.CompactTextString(m) }
func (*UserPublic) ProtoMessage()               {}
func (*UserPublic) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *UserPublic) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UserPublic) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UserPublic) GetStatus() UserStatus {
	if m != nil {
		return m.Status
	}
	return UserStatus_invited
}

func (m *UserPublic) GetFacebookID() string {
	if m != nil {
		return m.FacebookID
	}
	return ""
}

func (m *UserPublic) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

type UsersPublic struct {
	Users []*UserPublic `protobuf:"bytes,1,rep,name=users" json:"users,omitempty"`
}

func (m *UsersPublic) Reset()                    { *m = UsersPublic{} }
func (m *UsersPublic) String() string            { return proto.CompactTextString(m) }
func (*UsersPublic) ProtoMessage()               {}
func (*UsersPublic) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *UsersPublic) GetUsers() []*UserPublic {
	if m != nil {
		return m.Users
	}
	return nil
}

func init() {
	proto.RegisterType((*User)(nil), "protobuf.User")
	proto.RegisterType((*UserPublic)(nil), "protobuf.UserPublic")
	proto.RegisterType((*UsersPublic)(nil), "protobuf.UsersPublic")
	proto.RegisterEnum("protobuf.UserStatus", UserStatus_name, UserStatus_value)
}

func init() { proto.RegisterFile("user.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 337 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xcd, 0x6e, 0xea, 0x30,
	0x10, 0x85, 0xaf, 0x03, 0x09, 0xc9, 0x44, 0x42, 0xc8, 0x62, 0x61, 0xa1, 0xab, 0x7b, 0x23, 0x56,
	0x11, 0xaa, 0x82, 0x44, 0x37, 0xed, 0xb2, 0x52, 0x37, 0xdd, 0x54, 0x55, 0x4a, 0x1f, 0x20, 0xc1,
	0x03, 0xb5, 0x20, 0x18, 0xc5, 0x0e, 0x4f, 0xd3, 0xc7, 0xe9, 0x83, 0x55, 0xb6, 0x03, 0xa9, 0x58,
	0xf4, 0x67, 0x15, 0xcd, 0xcc, 0x99, 0x73, 0x3e, 0x4f, 0x00, 0x1a, 0x85, 0x75, 0x76, 0xa8, 0xa5,
	0x96, 0x34, 0xb4, 0x9f, 0xb2, 0x59, 0x4f, 0xfe, 0x6f, 0xa4, 0xdc, 0xec, 0x70, 0x7e, 0x6a, 0xcc,
	0xb5, 0xa8, 0x50, 0xe9, 0xa2, 0x3a, 0x38, 0xe9, 0xf4, 0xdd, 0x83, 0xfe, 0x8b, 0xc2, 0x9a, 0x0e,
	0xc1, 0x13, 0x9c, 0x91, 0x84, 0xa4, 0x51, 0xee, 0x09, 0x4e, 0xc7, 0xe0, 0x63, 0x55, 0x88, 0x1d,
	0xf3, 0x6c, 0xcb, 0x15, 0xf4, 0x06, 0xa2, 0x55, 0x8d, 0x85, 0x46, 0x7e, 0xa7, 0x59, 0x2f, 0x21,
	0x69, 0xbc, 0x98, 0x64, 0x2e, 0x23, 0x3b, 0x65, 0x64, 0xcb, 0x53, 0x46, 0xde, 0x89, 0xcd, 0x66,
	0x73, 0xe0, 0xed, 0x66, 0xff, 0xfb, 0xcd, 0xb3, 0x98, 0x4e, 0x20, 0x3c, 0x62, 0x2d, 0xd6, 0x02,
	0x39, 0xf3, 0x13, 0x92, 0x86, 0xf9, 0xb9, 0xa6, 0x7f, 0x21, 0x2a, 0x1a, 0xfd, 0xba, 0x94, 0x5b,
	0xdc, 0xb3, 0xc0, 0x92, 0x76, 0x0d, 0x7a, 0x05, 0x81, 0xd2, 0x85, 0x6e, 0x14, 0x1b, 0x24, 0x24,
	0x1d, 0x2e, 0xc6, 0x5d, 0x92, 0x79, 0xf3, 0xb3, 0x9d, 0xe5, 0xad, 0x86, 0xfe, 0x03, 0x58, 0x17,
	0x2b, 0x2c, 0xa5, 0xdc, 0x3e, 0xdc, 0xb3, 0xd0, 0x9a, 0x7d, 0xea, 0x18, 0x0e, 0x73, 0xe3, 0xc7,
	0xa2, 0x42, 0x16, 0xd9, 0xe9, 0xb9, 0x9e, 0xbe, 0x11, 0x00, 0x63, 0xf9, 0xd4, 0x94, 0x3b, 0xb1,
	0xfa, 0xe1, 0x31, 0x3b, 0xbc, 0xde, 0xaf, 0xf1, 0xfa, 0x5f, 0xe2, 0xf9, 0x17, 0x78, 0xb7, 0x10,
	0x1b, 0x47, 0xd5, 0xe2, 0xcd, 0xc0, 0x37, 0x23, 0xc5, 0x48, 0xd2, 0x4b, 0xe3, 0xcb, 0x5c, 0x27,
	0xca, 0x9d, 0x64, 0xb6, 0x70, 0x0f, 0x73, 0x30, 0x34, 0x86, 0x81, 0xd8, 0x1f, 0x85, 0x46, 0x3e,
	0xfa, 0x63, 0x8a, 0xf6, 0xff, 0x8e, 0x08, 0x05, 0x08, 0x2a, 0xac, 0x37, 0xc8, 0x47, 0x5e, 0x19,
	0x58, 0xbf, 0xeb, 0x8f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc4, 0x86, 0xac, 0x55, 0x94, 0x02, 0x00,
	0x00,
}

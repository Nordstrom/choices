// Code generated by protoc-gen-go.
// source: storage.proto
// DO NOT EDIT!

/*
Package storage is a generated protocol buffer package.

It is generated from these files:
	storage.proto

It has these top-level messages:
	AllRequest
	AllReply
	CreateRequest
	CreateReply
	ReadRequest
	ReadReply
	UpdateRequest
	UpdateReply
	DeleteRequest
	DeleteReply
	Namespace
	Experiment
	Param
	Value
*/
package storage

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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

// Environment structure
type Environment int32

const (
	Environment_Staging    Environment = 0
	Environment_Production Environment = 1
)

var Environment_name = map[int32]string{
	0: "Staging",
	1: "Production",
}
var Environment_value = map[string]int32{
	"Staging":    0,
	"Production": 1,
}

func (x Environment) String() string {
	return proto.EnumName(Environment_name, int32(x))
}
func (Environment) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// AllRequest retuns all the experiments from the given environment.
type AllRequest struct {
	Environment Environment `protobuf:"varint,1,opt,name=environment,enum=Environment" json:"environment,omitempty"`
}

func (m *AllRequest) Reset()                    { *m = AllRequest{} }
func (m *AllRequest) String() string            { return proto.CompactTextString(m) }
func (*AllRequest) ProtoMessage()               {}
func (*AllRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// The response message containing the Namespaces
type AllReply struct {
	Namespaces []*Namespace `protobuf:"bytes,1,rep,name=namespaces" json:"namespaces,omitempty"`
}

func (m *AllReply) Reset()                    { *m = AllReply{} }
func (m *AllReply) String() string            { return proto.CompactTextString(m) }
func (*AllReply) ProtoMessage()               {}
func (*AllReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *AllReply) GetNamespaces() []*Namespace {
	if m != nil {
		return m.Namespaces
	}
	return nil
}

// CreateRequest request message to create a new namespace in an environment.
type CreateRequest struct {
	Namespace   *Namespace  `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
	Environment Environment `protobuf:"varint,2,opt,name=environment,enum=Environment" json:"environment,omitempty"`
}

func (m *CreateRequest) Reset()                    { *m = CreateRequest{} }
func (m *CreateRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()               {}
func (*CreateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CreateRequest) GetNamespace() *Namespace {
	if m != nil {
		return m.Namespace
	}
	return nil
}

// CreateReply response containing the newly created Namespace.
type CreateReply struct {
	Namespace *Namespace `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
}

func (m *CreateReply) Reset()                    { *m = CreateReply{} }
func (m *CreateReply) String() string            { return proto.CompactTextString(m) }
func (*CreateReply) ProtoMessage()               {}
func (*CreateReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CreateReply) GetNamespace() *Namespace {
	if m != nil {
		return m.Namespace
	}
	return nil
}

// ReadRequest request message to get a namespace by name.
type ReadRequest struct {
	Name        string      `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Environment Environment `protobuf:"varint,2,opt,name=environment,enum=Environment" json:"environment,omitempty"`
}

func (m *ReadRequest) Reset()                    { *m = ReadRequest{} }
func (m *ReadRequest) String() string            { return proto.CompactTextString(m) }
func (*ReadRequest) ProtoMessage()               {}
func (*ReadRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

// ReadReply response containing the namespace requested.
type ReadReply struct {
	Namespace *Namespace `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
}

func (m *ReadReply) Reset()                    { *m = ReadReply{} }
func (m *ReadReply) String() string            { return proto.CompactTextString(m) }
func (*ReadReply) ProtoMessage()               {}
func (*ReadReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ReadReply) GetNamespace() *Namespace {
	if m != nil {
		return m.Namespace
	}
	return nil
}

// UpdateRequest request message to update an existing namespace in an
// environment.
type UpdateRequest struct {
	Namespace   *Namespace  `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
	Environment Environment `protobuf:"varint,2,opt,name=environment,enum=Environment" json:"environment,omitempty"`
}

func (m *UpdateRequest) Reset()                    { *m = UpdateRequest{} }
func (m *UpdateRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateRequest) ProtoMessage()               {}
func (*UpdateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *UpdateRequest) GetNamespace() *Namespace {
	if m != nil {
		return m.Namespace
	}
	return nil
}

// UpdateReply response containing the updated namespace.
type UpdateReply struct {
	Namespace *Namespace `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
}

func (m *UpdateReply) Reset()                    { *m = UpdateReply{} }
func (m *UpdateReply) String() string            { return proto.CompactTextString(m) }
func (*UpdateReply) ProtoMessage()               {}
func (*UpdateReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *UpdateReply) GetNamespace() *Namespace {
	if m != nil {
		return m.Namespace
	}
	return nil
}

// DeleteRequest request message to delete an existing namespace from an
// environment.
type DeleteRequest struct {
	Name        string      `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Environment Environment `protobuf:"varint,2,opt,name=environment,enum=Environment" json:"environment,omitempty"`
}

func (m *DeleteRequest) Reset()                    { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()               {}
func (*DeleteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

// DeleteReply response containing the deleted namespace.
type DeleteReply struct {
	Namespace *Namespace `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
}

func (m *DeleteReply) Reset()                    { *m = DeleteReply{} }
func (m *DeleteReply) String() string            { return proto.CompactTextString(m) }
func (*DeleteReply) ProtoMessage()               {}
func (*DeleteReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *DeleteReply) GetNamespace() *Namespace {
	if m != nil {
		return m.Namespace
	}
	return nil
}

// Namespace structure
type Namespace struct {
	Name        string        `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Experiments []*Experiment `protobuf:"bytes,3,rep,name=experiments" json:"experiments,omitempty"`
}

func (m *Namespace) Reset()                    { *m = Namespace{} }
func (m *Namespace) String() string            { return proto.CompactTextString(m) }
func (*Namespace) ProtoMessage()               {}
func (*Namespace) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *Namespace) GetExperiments() []*Experiment {
	if m != nil {
		return m.Experiments
	}
	return nil
}

// Experiment structure
type Experiment struct {
	Name     string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Segments []byte   `protobuf:"bytes,2,opt,name=segments,proto3" json:"segments,omitempty"`
	Params   []*Param `protobuf:"bytes,3,rep,name=params" json:"params,omitempty"`
}

func (m *Experiment) Reset()                    { *m = Experiment{} }
func (m *Experiment) String() string            { return proto.CompactTextString(m) }
func (*Experiment) ProtoMessage()               {}
func (*Experiment) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *Experiment) GetParams() []*Param {
	if m != nil {
		return m.Params
	}
	return nil
}

// Param structure
type Param struct {
	Name  string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Value *Value `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Param) Reset()                    { *m = Param{} }
func (m *Param) String() string            { return proto.CompactTextString(m) }
func (*Param) ProtoMessage()               {}
func (*Param) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *Param) GetValue() *Value {
	if m != nil {
		return m.Value
	}
	return nil
}

// Value structure
type Value struct {
	Choices []string  `protobuf:"bytes,1,rep,name=choices" json:"choices,omitempty"`
	Weights []float64 `protobuf:"fixed64,2,rep,packed,name=weights" json:"weights,omitempty"`
}

func (m *Value) Reset()                    { *m = Value{} }
func (m *Value) String() string            { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()               {}
func (*Value) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func init() {
	proto.RegisterType((*AllRequest)(nil), "AllRequest")
	proto.RegisterType((*AllReply)(nil), "AllReply")
	proto.RegisterType((*CreateRequest)(nil), "CreateRequest")
	proto.RegisterType((*CreateReply)(nil), "CreateReply")
	proto.RegisterType((*ReadRequest)(nil), "ReadRequest")
	proto.RegisterType((*ReadReply)(nil), "ReadReply")
	proto.RegisterType((*UpdateRequest)(nil), "UpdateRequest")
	proto.RegisterType((*UpdateReply)(nil), "UpdateReply")
	proto.RegisterType((*DeleteRequest)(nil), "DeleteRequest")
	proto.RegisterType((*DeleteReply)(nil), "DeleteReply")
	proto.RegisterType((*Namespace)(nil), "Namespace")
	proto.RegisterType((*Experiment)(nil), "Experiment")
	proto.RegisterType((*Param)(nil), "Param")
	proto.RegisterType((*Value)(nil), "Value")
	proto.RegisterEnum("Environment", Environment_name, Environment_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for ElwinStorage service

type ElwinStorageClient interface {
	// All returns all the namespaces for a given environment.
	All(ctx context.Context, in *AllRequest, opts ...grpc.CallOption) (*AllReply, error)
	// Create creates a namespace in the given environment.
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateReply, error)
	// Read returns the namespace matching the supplied name from the given
	// environment.
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadReply, error)
	// Update replaces the namespace in the given environment with the namespace
	// supplied.
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateReply, error)
	// Delete deletes the namespace from the given environment.
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error)
}

type elwinStorageClient struct {
	cc *grpc.ClientConn
}

func NewElwinStorageClient(cc *grpc.ClientConn) ElwinStorageClient {
	return &elwinStorageClient{cc}
}

func (c *elwinStorageClient) All(ctx context.Context, in *AllRequest, opts ...grpc.CallOption) (*AllReply, error) {
	out := new(AllReply)
	err := grpc.Invoke(ctx, "/ElwinStorage/All", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elwinStorageClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateReply, error) {
	out := new(CreateReply)
	err := grpc.Invoke(ctx, "/ElwinStorage/Create", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elwinStorageClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadReply, error) {
	out := new(ReadReply)
	err := grpc.Invoke(ctx, "/ElwinStorage/Read", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elwinStorageClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateReply, error) {
	out := new(UpdateReply)
	err := grpc.Invoke(ctx, "/ElwinStorage/Update", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *elwinStorageClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error) {
	out := new(DeleteReply)
	err := grpc.Invoke(ctx, "/ElwinStorage/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ElwinStorage service

type ElwinStorageServer interface {
	// All returns all the namespaces for a given environment.
	All(context.Context, *AllRequest) (*AllReply, error)
	// Create creates a namespace in the given environment.
	Create(context.Context, *CreateRequest) (*CreateReply, error)
	// Read returns the namespace matching the supplied name from the given
	// environment.
	Read(context.Context, *ReadRequest) (*ReadReply, error)
	// Update replaces the namespace in the given environment with the namespace
	// supplied.
	Update(context.Context, *UpdateRequest) (*UpdateReply, error)
	// Delete deletes the namespace from the given environment.
	Delete(context.Context, *DeleteRequest) (*DeleteReply, error)
}

func RegisterElwinStorageServer(s *grpc.Server, srv ElwinStorageServer) {
	s.RegisterService(&_ElwinStorage_serviceDesc, srv)
}

func _ElwinStorage_All_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ElwinStorageServer).All(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElwinStorage/All",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ElwinStorageServer).All(ctx, req.(*AllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ElwinStorage_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ElwinStorageServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElwinStorage/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ElwinStorageServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ElwinStorage_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ElwinStorageServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElwinStorage/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ElwinStorageServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ElwinStorage_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ElwinStorageServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElwinStorage/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ElwinStorageServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ElwinStorage_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ElwinStorageServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ElwinStorage/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ElwinStorageServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ElwinStorage_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ElwinStorage",
	HandlerType: (*ElwinStorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "All",
			Handler:    _ElwinStorage_All_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _ElwinStorage_Create_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _ElwinStorage_Read_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _ElwinStorage_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ElwinStorage_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("storage.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 464 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xbc, 0x54, 0x41, 0x8f, 0xd3, 0x3c,
	0x14, 0x6c, 0xbe, 0x6e, 0xbb, 0x9b, 0x71, 0x1a, 0xad, 0x7c, 0xaa, 0xa2, 0x4f, 0xb0, 0xf2, 0x29,
	0xaa, 0x84, 0x0f, 0x45, 0xb0, 0x42, 0x70, 0x59, 0x41, 0xaf, 0xab, 0xc5, 0x15, 0x9c, 0xb8, 0x98,
	0xd6, 0xca, 0x5a, 0x4a, 0x93, 0x90, 0xa4, 0xbb, 0xf0, 0x33, 0xf9, 0x47, 0xc8, 0x36, 0x4e, 0x52,
	0x54, 0x09, 0x2a, 0x21, 0x6e, 0x7d, 0x6f, 0xf2, 0x66, 0xa6, 0xd6, 0x9b, 0x87, 0x59, 0xd3, 0x96,
	0xb5, 0xcc, 0x14, 0xaf, 0xea, 0xb2, 0x2d, 0xd9, 0x1b, 0xe0, 0x26, 0xcf, 0x85, 0xfa, 0xb2, 0x57,
	0x4d, 0x4b, 0x39, 0x88, 0x2a, 0x1e, 0x74, 0x5d, 0x16, 0x3b, 0x55, 0xb4, 0xf3, 0xe0, 0x2a, 0x48,
	0xe3, 0x65, 0xc4, 0x57, 0x7d, 0x4f, 0x0c, 0x3f, 0x60, 0x2f, 0x71, 0x61, 0xa7, 0xab, 0xfc, 0x1b,
	0x5d, 0x00, 0x85, 0xdc, 0xa9, 0xa6, 0x92, 0x1b, 0xd5, 0xcc, 0x83, 0xab, 0x71, 0x4a, 0x96, 0xe0,
	0xb7, 0xbe, 0x25, 0x06, 0x28, 0xd3, 0x98, 0xbd, 0xad, 0x95, 0x6c, 0x95, 0x17, 0x4e, 0x11, 0x76,
	0xb0, 0x95, 0x3d, 0x9c, 0xed, 0xc1, 0x5f, 0x2d, 0xfe, 0xf7, 0x3b, 0x8b, 0xd7, 0x20, 0x5e, 0xca,
	0xb8, 0xfc, 0x63, 0x21, 0xf6, 0x1e, 0x44, 0x28, 0xb9, 0xf5, 0x0e, 0x29, 0xce, 0x0c, 0x66, 0x67,
	0x42, 0x61, 0x7f, 0x9f, 0xec, 0xe5, 0x05, 0x42, 0x47, 0x79, 0x9a, 0x13, 0x8d, 0xd9, 0x87, 0x6a,
	0xfb, 0xaf, 0x5e, 0xcb, 0x4b, 0x9d, 0xe6, 0x71, 0x8d, 0xd9, 0x3b, 0x95, 0xab, 0xde, 0xe3, 0xdf,
	0x78, 0xaf, 0x6b, 0x10, 0x4f, 0x7a, 0x9a, 0x9b, 0x5b, 0x84, 0x5d, 0xff, 0xa8, 0x93, 0x67, 0x20,
	0xea, 0x6b, 0xa5, 0x6a, 0x6d, 0x74, 0x9a, 0xf9, 0xd8, 0x6e, 0x2b, 0xe1, 0xab, 0xae, 0x27, 0x86,
	0x38, 0xfb, 0x04, 0xf4, 0xd0, 0x51, 0xc2, 0x04, 0x17, 0x8d, 0xca, 0x1c, 0x9b, 0xf9, 0x5f, 0x91,
	0xe8, 0x6a, 0xfa, 0x04, 0xd3, 0x4a, 0xd6, 0x72, 0xe7, 0x75, 0xa6, 0xfc, 0xce, 0x94, 0xe2, 0x67,
	0x97, 0xbd, 0xc2, 0xc4, 0x36, 0x8e, 0x12, 0xff, 0x8f, 0xc9, 0x83, 0xcc, 0xf7, 0xca, 0xb2, 0x9a,
	0xd9, 0x8f, 0xa6, 0x12, 0xae, 0xc9, 0x5e, 0x63, 0x62, 0x6b, 0x3a, 0xc7, 0xf9, 0xe6, 0xbe, 0xd4,
	0x3e, 0x7a, 0xa1, 0xf0, 0xa5, 0x41, 0x1e, 0x95, 0xce, 0xee, 0xad, 0xb1, 0x71, 0x1a, 0x08, 0x5f,
	0x2e, 0x16, 0x20, 0x83, 0xa7, 0xa7, 0x04, 0xe7, 0xeb, 0x56, 0x66, 0xba, 0xc8, 0x2e, 0x47, 0x34,
	0x06, 0xee, 0xea, 0x72, 0xbb, 0xdf, 0xb4, 0xba, 0x2c, 0x2e, 0x83, 0xe5, 0xf7, 0x00, 0xd1, 0x2a,
	0x7f, 0xd4, 0xc5, 0xda, 0x9d, 0x0f, 0xfa, 0x14, 0xe3, 0x9b, 0x3c, 0xa7, 0x84, 0xf7, 0xe7, 0x23,
	0x09, 0xb9, 0xbf, 0x06, 0x6c, 0x44, 0x53, 0x4c, 0x5d, 0xf0, 0x68, 0xcc, 0x0f, 0xc2, 0x9e, 0x44,
	0x7c, 0x90, 0x48, 0x36, 0xa2, 0x0c, 0x67, 0x26, 0x16, 0x34, 0xe2, 0x83, 0xc0, 0x25, 0xe0, 0x5d,
	0x56, 0x1c, 0x9b, 0x5b, 0x4c, 0x1a, 0xf3, 0x83, 0x30, 0x24, 0x11, 0x1f, 0x6c, 0xac, 0xfb, 0xd2,
	0x2d, 0x0d, 0x8d, 0xf9, 0xc1, 0x4a, 0x26, 0x11, 0x1f, 0x6c, 0x13, 0x1b, 0x7d, 0x9e, 0xda, 0x13,
	0xf8, 0xfc, 0x47, 0x00, 0x00, 0x00, 0xff, 0xff, 0x43, 0x10, 0x37, 0xb7, 0x13, 0x05, 0x00, 0x00,
}

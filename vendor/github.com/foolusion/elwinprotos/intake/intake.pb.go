// Code generated by protoc-gen-gogo.
// source: intake.proto
// DO NOT EDIT!

/*
	Package intake is a generated protocol buffer package.

	It is generated from these files:
		intake.proto

	It has these top-level messages:
		ExperimentIntakeRequest
		ExperimentIntakeReply
		ExperimentMetadata
*/
package intake

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api"
import elwin_storage "github.com/foolusion/elwinprotos/storage"

import strings "strings"
import reflect "reflect"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// ExperimentIntakeRequest creates an experiment in the database and sends a notification for reviewers
type ExperimentIntakeRequest struct {
	Metadata   *ExperimentMetadata       `protobuf:"bytes,1,opt,name=metadata" json:"metadata,omitempty"`
	Experiment *elwin_storage.Experiment `protobuf:"bytes,2,opt,name=experiment" json:"experiment,omitempty"`
}

func (m *ExperimentIntakeRequest) Reset()                    { *m = ExperimentIntakeRequest{} }
func (*ExperimentIntakeRequest) ProtoMessage()               {}
func (*ExperimentIntakeRequest) Descriptor() ([]byte, []int) { return fileDescriptorIntake, []int{0} }

func (m *ExperimentIntakeRequest) GetMetadata() *ExperimentMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *ExperimentIntakeRequest) GetExperiment() *elwin_storage.Experiment {
	if m != nil {
		return m.Experiment
	}
	return nil
}

type ExperimentIntakeReply struct {
}

func (m *ExperimentIntakeReply) Reset()                    { *m = ExperimentIntakeReply{} }
func (*ExperimentIntakeReply) ProtoMessage()               {}
func (*ExperimentIntakeReply) Descriptor() ([]byte, []int) { return fileDescriptorIntake, []int{1} }

// ExperimentMetadata all the junk that elwin doesn't care about
type ExperimentMetadata struct {
	UserID             string `protobuf:"bytes,1,opt,name=userID,proto3" json:"userID,omitempty"`
	ProgramManagerID   string `protobuf:"bytes,2,opt,name=programManagerID,proto3" json:"programManagerID,omitempty"`
	ProductManagerID   string `protobuf:"bytes,3,opt,name=productManagerID,proto3" json:"productManagerID,omitempty"`
	Hypothesis         string `protobuf:"bytes,4,opt,name=hypothesis,proto3" json:"hypothesis,omitempty"`
	Kpi                string `protobuf:"bytes,5,opt,name=kpi,proto3" json:"kpi,omitempty"`
	TimeBound          bool   `protobuf:"varint,6,opt,name=timeBound,proto3" json:"timeBound,omitempty"`
	PlannedStartTime   string `protobuf:"bytes,7,opt,name=plannedStartTime,proto3" json:"plannedStartTime,omitempty"`
	PlannedEndTime     string `protobuf:"bytes,8,opt,name=plannedEndTime,proto3" json:"plannedEndTime,omitempty"`
	ActualStartTime    string `protobuf:"bytes,9,opt,name=actualStartTime,proto3" json:"actualStartTime,omitempty"`
	ActualEndTime      string `protobuf:"bytes,10,opt,name=actualEndTime,proto3" json:"actualEndTime,omitempty"`
	ActionPlanNegative string `protobuf:"bytes,11,opt,name=actionPlanNegative,proto3" json:"actionPlanNegative,omitempty"`
	ActionPlanNeutral  string `protobuf:"bytes,12,opt,name=actionPlanNeutral,proto3" json:"actionPlanNeutral,omitempty"`
	ExperimentType     string `protobuf:"bytes,13,opt,name=experimentType,proto3" json:"experimentType,omitempty"`
}

func (m *ExperimentMetadata) Reset()                    { *m = ExperimentMetadata{} }
func (*ExperimentMetadata) ProtoMessage()               {}
func (*ExperimentMetadata) Descriptor() ([]byte, []int) { return fileDescriptorIntake, []int{2} }

func (m *ExperimentMetadata) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *ExperimentMetadata) GetProgramManagerID() string {
	if m != nil {
		return m.ProgramManagerID
	}
	return ""
}

func (m *ExperimentMetadata) GetProductManagerID() string {
	if m != nil {
		return m.ProductManagerID
	}
	return ""
}

func (m *ExperimentMetadata) GetHypothesis() string {
	if m != nil {
		return m.Hypothesis
	}
	return ""
}

func (m *ExperimentMetadata) GetKpi() string {
	if m != nil {
		return m.Kpi
	}
	return ""
}

func (m *ExperimentMetadata) GetTimeBound() bool {
	if m != nil {
		return m.TimeBound
	}
	return false
}

func (m *ExperimentMetadata) GetPlannedStartTime() string {
	if m != nil {
		return m.PlannedStartTime
	}
	return ""
}

func (m *ExperimentMetadata) GetPlannedEndTime() string {
	if m != nil {
		return m.PlannedEndTime
	}
	return ""
}

func (m *ExperimentMetadata) GetActualStartTime() string {
	if m != nil {
		return m.ActualStartTime
	}
	return ""
}

func (m *ExperimentMetadata) GetActualEndTime() string {
	if m != nil {
		return m.ActualEndTime
	}
	return ""
}

func (m *ExperimentMetadata) GetActionPlanNegative() string {
	if m != nil {
		return m.ActionPlanNegative
	}
	return ""
}

func (m *ExperimentMetadata) GetActionPlanNeutral() string {
	if m != nil {
		return m.ActionPlanNeutral
	}
	return ""
}

func (m *ExperimentMetadata) GetExperimentType() string {
	if m != nil {
		return m.ExperimentType
	}
	return ""
}

func init() {
	proto.RegisterType((*ExperimentIntakeRequest)(nil), "elwin.intake.ExperimentIntakeRequest")
	proto.RegisterType((*ExperimentIntakeReply)(nil), "elwin.intake.ExperimentIntakeReply")
	proto.RegisterType((*ExperimentMetadata)(nil), "elwin.intake.ExperimentMetadata")
}
func (this *ExperimentIntakeRequest) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*ExperimentIntakeRequest)
	if !ok {
		that2, ok := that.(ExperimentIntakeRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if !this.Metadata.Equal(that1.Metadata) {
		return false
	}
	if !this.Experiment.Equal(that1.Experiment) {
		return false
	}
	return true
}
func (this *ExperimentIntakeReply) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*ExperimentIntakeReply)
	if !ok {
		that2, ok := that.(ExperimentIntakeReply)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	return true
}
func (this *ExperimentMetadata) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*ExperimentMetadata)
	if !ok {
		that2, ok := that.(ExperimentMetadata)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.UserID != that1.UserID {
		return false
	}
	if this.ProgramManagerID != that1.ProgramManagerID {
		return false
	}
	if this.ProductManagerID != that1.ProductManagerID {
		return false
	}
	if this.Hypothesis != that1.Hypothesis {
		return false
	}
	if this.Kpi != that1.Kpi {
		return false
	}
	if this.TimeBound != that1.TimeBound {
		return false
	}
	if this.PlannedStartTime != that1.PlannedStartTime {
		return false
	}
	if this.PlannedEndTime != that1.PlannedEndTime {
		return false
	}
	if this.ActualStartTime != that1.ActualStartTime {
		return false
	}
	if this.ActualEndTime != that1.ActualEndTime {
		return false
	}
	if this.ActionPlanNegative != that1.ActionPlanNegative {
		return false
	}
	if this.ActionPlanNeutral != that1.ActionPlanNeutral {
		return false
	}
	if this.ExperimentType != that1.ExperimentType {
		return false
	}
	return true
}
func (this *ExperimentIntakeRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&intake.ExperimentIntakeRequest{")
	if this.Metadata != nil {
		s = append(s, "Metadata: "+fmt.Sprintf("%#v", this.Metadata)+",\n")
	}
	if this.Experiment != nil {
		s = append(s, "Experiment: "+fmt.Sprintf("%#v", this.Experiment)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ExperimentIntakeReply) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 4)
	s = append(s, "&intake.ExperimentIntakeReply{")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ExperimentMetadata) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 17)
	s = append(s, "&intake.ExperimentMetadata{")
	s = append(s, "UserID: "+fmt.Sprintf("%#v", this.UserID)+",\n")
	s = append(s, "ProgramManagerID: "+fmt.Sprintf("%#v", this.ProgramManagerID)+",\n")
	s = append(s, "ProductManagerID: "+fmt.Sprintf("%#v", this.ProductManagerID)+",\n")
	s = append(s, "Hypothesis: "+fmt.Sprintf("%#v", this.Hypothesis)+",\n")
	s = append(s, "Kpi: "+fmt.Sprintf("%#v", this.Kpi)+",\n")
	s = append(s, "TimeBound: "+fmt.Sprintf("%#v", this.TimeBound)+",\n")
	s = append(s, "PlannedStartTime: "+fmt.Sprintf("%#v", this.PlannedStartTime)+",\n")
	s = append(s, "PlannedEndTime: "+fmt.Sprintf("%#v", this.PlannedEndTime)+",\n")
	s = append(s, "ActualStartTime: "+fmt.Sprintf("%#v", this.ActualStartTime)+",\n")
	s = append(s, "ActualEndTime: "+fmt.Sprintf("%#v", this.ActualEndTime)+",\n")
	s = append(s, "ActionPlanNegative: "+fmt.Sprintf("%#v", this.ActionPlanNegative)+",\n")
	s = append(s, "ActionPlanNeutral: "+fmt.Sprintf("%#v", this.ActionPlanNeutral)+",\n")
	s = append(s, "ExperimentType: "+fmt.Sprintf("%#v", this.ExperimentType)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringIntake(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ExperimentIntake service

type ExperimentIntakeClient interface {
	// ExperimentIntake takes a request from a web form and creates the
	// experiment in the data store.
	ExperimentIntake(ctx context.Context, in *ExperimentIntakeRequest, opts ...grpc.CallOption) (*ExperimentIntakeReply, error)
}

type experimentIntakeClient struct {
	cc *grpc.ClientConn
}

func NewExperimentIntakeClient(cc *grpc.ClientConn) ExperimentIntakeClient {
	return &experimentIntakeClient{cc}
}

func (c *experimentIntakeClient) ExperimentIntake(ctx context.Context, in *ExperimentIntakeRequest, opts ...grpc.CallOption) (*ExperimentIntakeReply, error) {
	out := new(ExperimentIntakeReply)
	err := grpc.Invoke(ctx, "/elwin.intake.ExperimentIntake/ExperimentIntake", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ExperimentIntake service

type ExperimentIntakeServer interface {
	// ExperimentIntake takes a request from a web form and creates the
	// experiment in the data store.
	ExperimentIntake(context.Context, *ExperimentIntakeRequest) (*ExperimentIntakeReply, error)
}

func RegisterExperimentIntakeServer(s *grpc.Server, srv ExperimentIntakeServer) {
	s.RegisterService(&_ExperimentIntake_serviceDesc, srv)
}

func _ExperimentIntake_ExperimentIntake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExperimentIntakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExperimentIntakeServer).ExperimentIntake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elwin.intake.ExperimentIntake/ExperimentIntake",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExperimentIntakeServer).ExperimentIntake(ctx, req.(*ExperimentIntakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ExperimentIntake_serviceDesc = grpc.ServiceDesc{
	ServiceName: "elwin.intake.ExperimentIntake",
	HandlerType: (*ExperimentIntakeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExperimentIntake",
			Handler:    _ExperimentIntake_ExperimentIntake_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "intake.proto",
}

func (m *ExperimentIntakeRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExperimentIntakeRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Metadata != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintIntake(dAtA, i, uint64(m.Metadata.Size()))
		n1, err := m.Metadata.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.Experiment != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintIntake(dAtA, i, uint64(m.Experiment.Size()))
		n2, err := m.Experiment.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func (m *ExperimentIntakeReply) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExperimentIntakeReply) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *ExperimentMetadata) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExperimentMetadata) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.UserID) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintIntake(dAtA, i, uint64(len(m.UserID)))
		i += copy(dAtA[i:], m.UserID)
	}
	if len(m.ProgramManagerID) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintIntake(dAtA, i, uint64(len(m.ProgramManagerID)))
		i += copy(dAtA[i:], m.ProgramManagerID)
	}
	if len(m.ProductManagerID) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintIntake(dAtA, i, uint64(len(m.ProductManagerID)))
		i += copy(dAtA[i:], m.ProductManagerID)
	}
	if len(m.Hypothesis) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintIntake(dAtA, i, uint64(len(m.Hypothesis)))
		i += copy(dAtA[i:], m.Hypothesis)
	}
	if len(m.Kpi) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintIntake(dAtA, i, uint64(len(m.Kpi)))
		i += copy(dAtA[i:], m.Kpi)
	}
	if m.TimeBound {
		dAtA[i] = 0x30
		i++
		if m.TimeBound {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if len(m.PlannedStartTime) > 0 {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintIntake(dAtA, i, uint64(len(m.PlannedStartTime)))
		i += copy(dAtA[i:], m.PlannedStartTime)
	}
	if len(m.PlannedEndTime) > 0 {
		dAtA[i] = 0x42
		i++
		i = encodeVarintIntake(dAtA, i, uint64(len(m.PlannedEndTime)))
		i += copy(dAtA[i:], m.PlannedEndTime)
	}
	if len(m.ActualStartTime) > 0 {
		dAtA[i] = 0x4a
		i++
		i = encodeVarintIntake(dAtA, i, uint64(len(m.ActualStartTime)))
		i += copy(dAtA[i:], m.ActualStartTime)
	}
	if len(m.ActualEndTime) > 0 {
		dAtA[i] = 0x52
		i++
		i = encodeVarintIntake(dAtA, i, uint64(len(m.ActualEndTime)))
		i += copy(dAtA[i:], m.ActualEndTime)
	}
	if len(m.ActionPlanNegative) > 0 {
		dAtA[i] = 0x5a
		i++
		i = encodeVarintIntake(dAtA, i, uint64(len(m.ActionPlanNegative)))
		i += copy(dAtA[i:], m.ActionPlanNegative)
	}
	if len(m.ActionPlanNeutral) > 0 {
		dAtA[i] = 0x62
		i++
		i = encodeVarintIntake(dAtA, i, uint64(len(m.ActionPlanNeutral)))
		i += copy(dAtA[i:], m.ActionPlanNeutral)
	}
	if len(m.ExperimentType) > 0 {
		dAtA[i] = 0x6a
		i++
		i = encodeVarintIntake(dAtA, i, uint64(len(m.ExperimentType)))
		i += copy(dAtA[i:], m.ExperimentType)
	}
	return i, nil
}

func encodeFixed64Intake(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Intake(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintIntake(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ExperimentIntakeRequest) Size() (n int) {
	var l int
	_ = l
	if m.Metadata != nil {
		l = m.Metadata.Size()
		n += 1 + l + sovIntake(uint64(l))
	}
	if m.Experiment != nil {
		l = m.Experiment.Size()
		n += 1 + l + sovIntake(uint64(l))
	}
	return n
}

func (m *ExperimentIntakeReply) Size() (n int) {
	var l int
	_ = l
	return n
}

func (m *ExperimentMetadata) Size() (n int) {
	var l int
	_ = l
	l = len(m.UserID)
	if l > 0 {
		n += 1 + l + sovIntake(uint64(l))
	}
	l = len(m.ProgramManagerID)
	if l > 0 {
		n += 1 + l + sovIntake(uint64(l))
	}
	l = len(m.ProductManagerID)
	if l > 0 {
		n += 1 + l + sovIntake(uint64(l))
	}
	l = len(m.Hypothesis)
	if l > 0 {
		n += 1 + l + sovIntake(uint64(l))
	}
	l = len(m.Kpi)
	if l > 0 {
		n += 1 + l + sovIntake(uint64(l))
	}
	if m.TimeBound {
		n += 2
	}
	l = len(m.PlannedStartTime)
	if l > 0 {
		n += 1 + l + sovIntake(uint64(l))
	}
	l = len(m.PlannedEndTime)
	if l > 0 {
		n += 1 + l + sovIntake(uint64(l))
	}
	l = len(m.ActualStartTime)
	if l > 0 {
		n += 1 + l + sovIntake(uint64(l))
	}
	l = len(m.ActualEndTime)
	if l > 0 {
		n += 1 + l + sovIntake(uint64(l))
	}
	l = len(m.ActionPlanNegative)
	if l > 0 {
		n += 1 + l + sovIntake(uint64(l))
	}
	l = len(m.ActionPlanNeutral)
	if l > 0 {
		n += 1 + l + sovIntake(uint64(l))
	}
	l = len(m.ExperimentType)
	if l > 0 {
		n += 1 + l + sovIntake(uint64(l))
	}
	return n
}

func sovIntake(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozIntake(x uint64) (n int) {
	return sovIntake(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *ExperimentIntakeRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ExperimentIntakeRequest{`,
		`Metadata:` + strings.Replace(fmt.Sprintf("%v", this.Metadata), "ExperimentMetadata", "ExperimentMetadata", 1) + `,`,
		`Experiment:` + strings.Replace(fmt.Sprintf("%v", this.Experiment), "Experiment", "elwin_storage.Experiment", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ExperimentIntakeReply) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ExperimentIntakeReply{`,
		`}`,
	}, "")
	return s
}
func (this *ExperimentMetadata) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ExperimentMetadata{`,
		`UserID:` + fmt.Sprintf("%v", this.UserID) + `,`,
		`ProgramManagerID:` + fmt.Sprintf("%v", this.ProgramManagerID) + `,`,
		`ProductManagerID:` + fmt.Sprintf("%v", this.ProductManagerID) + `,`,
		`Hypothesis:` + fmt.Sprintf("%v", this.Hypothesis) + `,`,
		`Kpi:` + fmt.Sprintf("%v", this.Kpi) + `,`,
		`TimeBound:` + fmt.Sprintf("%v", this.TimeBound) + `,`,
		`PlannedStartTime:` + fmt.Sprintf("%v", this.PlannedStartTime) + `,`,
		`PlannedEndTime:` + fmt.Sprintf("%v", this.PlannedEndTime) + `,`,
		`ActualStartTime:` + fmt.Sprintf("%v", this.ActualStartTime) + `,`,
		`ActualEndTime:` + fmt.Sprintf("%v", this.ActualEndTime) + `,`,
		`ActionPlanNegative:` + fmt.Sprintf("%v", this.ActionPlanNegative) + `,`,
		`ActionPlanNeutral:` + fmt.Sprintf("%v", this.ActionPlanNeutral) + `,`,
		`ExperimentType:` + fmt.Sprintf("%v", this.ExperimentType) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringIntake(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *ExperimentIntakeRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIntake
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ExperimentIntakeRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExperimentIntakeRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metadata", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Metadata == nil {
				m.Metadata = &ExperimentMetadata{}
			}
			if err := m.Metadata.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Experiment", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Experiment == nil {
				m.Experiment = &elwin_storage.Experiment{}
			}
			if err := m.Experiment.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIntake(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIntake
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ExperimentIntakeReply) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIntake
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ExperimentIntakeReply: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExperimentIntakeReply: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipIntake(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIntake
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ExperimentMetadata) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIntake
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ExperimentMetadata: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExperimentMetadata: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UserID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProgramManagerID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProgramManagerID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProductManagerID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProductManagerID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hypothesis", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hypothesis = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Kpi", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Kpi = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TimeBound", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.TimeBound = bool(v != 0)
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PlannedStartTime", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PlannedStartTime = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PlannedEndTime", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PlannedEndTime = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActualStartTime", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ActualStartTime = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActualEndTime", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ActualEndTime = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActionPlanNegative", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ActionPlanNegative = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActionPlanNeutral", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ActionPlanNeutral = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExperimentType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntake
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExperimentType = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIntake(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIntake
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipIntake(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIntake
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIntake
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthIntake
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowIntake
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipIntake(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthIntake = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIntake   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("intake.proto", fileDescriptorIntake) }

var fileDescriptorIntake = []byte{
	// 517 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0x3f, 0x6f, 0xd3, 0x40,
	0x14, 0xc0, 0x73, 0x0d, 0x84, 0xe4, 0x35, 0x81, 0x70, 0x12, 0xd4, 0x8d, 0xaa, 0x53, 0x14, 0x0a,
	0x44, 0x15, 0xd8, 0xa2, 0x4c, 0x20, 0xa6, 0x8a, 0x0e, 0x1d, 0x8a, 0x90, 0x29, 0x0b, 0xdb, 0x6b,
	0x7c, 0x38, 0xa7, 0xda, 0x77, 0xc6, 0x3e, 0x17, 0xb2, 0x21, 0xc4, 0x07, 0x40, 0x30, 0xf1, 0x0d,
	0xf8, 0x28, 0x8c, 0x95, 0x58, 0x18, 0x89, 0x61, 0x60, 0xec, 0x27, 0x40, 0x28, 0x67, 0x37, 0xce,
	0x9f, 0x22, 0x36, 0xdf, 0xef, 0xfd, 0xde, 0x9f, 0x93, 0xdf, 0x41, 0x53, 0x48, 0x8d, 0x47, 0xdc,
	0x8e, 0x62, 0xa5, 0x15, 0x6d, 0xf2, 0xe0, 0xb5, 0x90, 0x76, 0xce, 0x3a, 0x1b, 0xbe, 0x52, 0x7e,
	0xc0, 0x1d, 0x8c, 0x84, 0x83, 0x52, 0x2a, 0x8d, 0x5a, 0x28, 0x99, 0xe4, 0x6e, 0xa7, 0x95, 0x68,
	0x15, 0xa3, 0x5f, 0xa4, 0xf6, 0x3e, 0x12, 0x58, 0xdb, 0x7d, 0x13, 0xf1, 0x58, 0x84, 0x5c, 0xea,
	0x3d, 0x53, 0xc1, 0xe5, 0xaf, 0x52, 0x9e, 0x68, 0xfa, 0x08, 0xea, 0x21, 0xd7, 0xe8, 0xa1, 0x46,
	0x8b, 0x74, 0x49, 0x7f, 0x75, 0xbb, 0x6b, 0xcf, 0x76, 0xb2, 0xcb, 0xc4, 0xfd, 0xc2, 0x73, 0xa7,
	0x19, 0xf4, 0x01, 0x00, 0x9f, 0xc6, 0xad, 0x15, 0x93, 0xbf, 0x5e, 0xe4, 0x9f, 0xcd, 0x50, 0x16,
	0x70, 0x67, 0xe4, 0xde, 0x1a, 0x5c, 0x5b, 0x9e, 0x29, 0x0a, 0x46, 0xbd, 0x3f, 0x55, 0xa0, 0xcb,
	0x4d, 0xe9, 0x75, 0xa8, 0xa5, 0x09, 0x8f, 0xf7, 0x1e, 0x9b, 0x31, 0x1b, 0x6e, 0x71, 0xa2, 0x5b,
	0xd0, 0x8e, 0x62, 0xe5, 0xc7, 0x18, 0xee, 0xa3, 0x44, 0xdf, 0x18, 0x2b, 0xc6, 0x58, 0xe2, 0x85,
	0xeb, 0xa5, 0x03, 0x5d, 0xba, 0xd5, 0xa9, 0x3b, 0xc7, 0x29, 0x03, 0x18, 0x8e, 0x22, 0xa5, 0x87,
	0x3c, 0x11, 0x89, 0x75, 0xc1, 0x58, 0x33, 0x84, 0xb6, 0xa1, 0x7a, 0x14, 0x09, 0xeb, 0xa2, 0x09,
	0x4c, 0x3e, 0xe9, 0x06, 0x34, 0xb4, 0x08, 0xf9, 0x8e, 0x4a, 0xa5, 0x67, 0xd5, 0xba, 0xa4, 0x5f,
	0x77, 0x4b, 0x60, 0x7a, 0x07, 0x28, 0x25, 0xf7, 0x9e, 0x69, 0x8c, 0xf5, 0x81, 0x08, 0xb9, 0x75,
	0xa9, 0xe8, 0xbd, 0xc0, 0xe9, 0x2d, 0xb8, 0x5c, 0xb0, 0x5d, 0xe9, 0x19, 0xb3, 0x6e, 0xcc, 0x05,
	0x4a, 0xfb, 0x70, 0x05, 0x07, 0x3a, 0xc5, 0xa0, 0x2c, 0xd9, 0x30, 0xe2, 0x22, 0xa6, 0x9b, 0xd0,
	0xca, 0xd1, 0x59, 0x41, 0x30, 0xde, 0x3c, 0xa4, 0x36, 0x50, 0x1c, 0x4c, 0x16, 0xe9, 0x69, 0x80,
	0xf2, 0x09, 0xf7, 0x51, 0x8b, 0x63, 0x6e, 0xad, 0x1a, 0xf5, 0x9c, 0x08, 0xbd, 0x03, 0x57, 0x67,
	0x69, 0xaa, 0x63, 0x0c, 0xac, 0xa6, 0xd1, 0x97, 0x03, 0x93, 0x5b, 0x95, 0xff, 0xff, 0x60, 0x14,
	0x71, 0xab, 0x95, 0xdf, 0x6a, 0x9e, 0x6e, 0x7f, 0x26, 0xd0, 0x5e, 0x5c, 0x0d, 0xfa, 0xfe, 0x3c,
	0x78, 0xf3, 0x5f, 0xab, 0x3a, 0xb7, 0xe3, 0x9d, 0x1b, 0xff, 0xd3, 0x26, 0x6b, 0xb7, 0xf9, 0xee,
	0xdb, 0xaf, 0x4f, 0x2b, 0xac, 0xb7, 0x6e, 0xde, 0xd4, 0xf1, 0x3d, 0xa7, 0x9c, 0xea, 0x6e, 0x9e,
	0xf8, 0x90, 0x6c, 0xed, 0x3c, 0x3f, 0x19, 0xb3, 0xca, 0xf7, 0x31, 0xab, 0x9c, 0x8e, 0x19, 0x79,
	0x9b, 0x31, 0xf2, 0x25, 0x63, 0xe4, 0x6b, 0xc6, 0xc8, 0x49, 0xc6, 0xc8, 0x8f, 0x8c, 0x91, 0xdf,
	0x19, 0xab, 0x9c, 0x66, 0x8c, 0x7c, 0xf8, 0xc9, 0x2a, 0x2f, 0x6e, 0xfb, 0x42, 0x0f, 0xd3, 0x43,
	0x7b, 0xa0, 0x42, 0xe7, 0xa5, 0x52, 0x41, 0x9a, 0x08, 0x25, 0x1d, 0x33, 0x8c, 0x79, 0x99, 0x89,
	0x93, 0x57, 0x3e, 0xac, 0x99, 0xe3, 0xfd, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x42, 0x88, 0x0b,
	0x2b, 0xf3, 0x03, 0x00, 0x00,
}

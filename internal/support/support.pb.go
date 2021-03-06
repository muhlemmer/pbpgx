// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: support.proto

//
//SPDX-License-Identifier: AGPL-3.0-only
//
//Copyright (C) 2021, Tim Möhlmann
//
//This program is free software: you can redistribute it and/or modify
//it under the terms of the GNU Affero General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU Affero General Public License for more details.
//
//You should have received a copy of the GNU Affero General Public License
//along with this program.  If not, see <https://www.gnu.org/licenses/>.

package support

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SimpleColumns int32

const (
	SimpleColumns_id      SimpleColumns = 0
	SimpleColumns_title   SimpleColumns = 1
	SimpleColumns_data    SimpleColumns = 2
	SimpleColumns_created SimpleColumns = 3
)

// Enum value maps for SimpleColumns.
var (
	SimpleColumns_name = map[int32]string{
		0: "id",
		1: "title",
		2: "data",
		3: "created",
	}
	SimpleColumns_value = map[string]int32{
		"id":      0,
		"title":   1,
		"data":    2,
		"created": 3,
	}
)

func (x SimpleColumns) Enum() *SimpleColumns {
	p := new(SimpleColumns)
	*p = x
	return p
}

func (x SimpleColumns) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SimpleColumns) Descriptor() protoreflect.EnumDescriptor {
	return file_support_proto_enumTypes[0].Descriptor()
}

func (SimpleColumns) Type() protoreflect.EnumType {
	return &file_support_proto_enumTypes[0]
}

func (x SimpleColumns) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SimpleColumns.Descriptor instead.
func (SimpleColumns) EnumDescriptor() ([]byte, []int) {
	return file_support_proto_rawDescGZIP(), []int{0}
}

// Supported destination types
type Supported struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bl   bool                     `protobuf:"varint,1,opt,name=bl,proto3" json:"bl,omitempty"`
	I32  int32                    `protobuf:"varint,2,opt,name=i32,proto3" json:"i32,omitempty"`
	I64  int64                    `protobuf:"varint,3,opt,name=i64,proto3" json:"i64,omitempty"`
	F    float32                  `protobuf:"fixed32,4,opt,name=f,proto3" json:"f,omitempty"`
	D    float64                  `protobuf:"fixed64,5,opt,name=d,proto3" json:"d,omitempty"`
	S    string                   `protobuf:"bytes,6,opt,name=s,proto3" json:"s,omitempty"`
	Bt   []byte                   `protobuf:"bytes,7,opt,name=bt,proto3" json:"bt,omitempty"`
	U32  uint32                   `protobuf:"varint,8,opt,name=u32,proto3" json:"u32,omitempty"`
	U64  uint64                   `protobuf:"varint,9,opt,name=u64,proto3" json:"u64,omitempty"`
	Ts   *timestamppb.Timestamp   `protobuf:"bytes,10,opt,name=ts,proto3" json:"ts,omitempty"`
	RBl  []bool                   `protobuf:"varint,11,rep,packed,name=r_bl,json=rBl,proto3" json:"r_bl,omitempty"`
	RI32 []int32                  `protobuf:"varint,12,rep,packed,name=r_i32,json=rI32,proto3" json:"r_i32,omitempty"`
	RI64 []int64                  `protobuf:"varint,13,rep,packed,name=r_i64,json=rI64,proto3" json:"r_i64,omitempty"`
	RF   []float32                `protobuf:"fixed32,14,rep,packed,name=r_f,json=rF,proto3" json:"r_f,omitempty"`
	RD   []float64                `protobuf:"fixed64,15,rep,packed,name=r_d,json=rD,proto3" json:"r_d,omitempty"`
	RS   []string                 `protobuf:"bytes,16,rep,name=r_s,json=rS,proto3" json:"r_s,omitempty"`
	RU32 []uint32                 `protobuf:"varint,17,rep,packed,name=r_u32,json=rU32,proto3" json:"r_u32,omitempty"`
	RBt  [][]byte                 `protobuf:"bytes,18,rep,name=r_bt,json=rBt,proto3" json:"r_bt,omitempty"`
	RU64 []uint64                 `protobuf:"varint,19,rep,packed,name=r_u64,json=rU64,proto3" json:"r_u64,omitempty"`
	RTs  []*timestamppb.Timestamp `protobuf:"bytes,20,rep,name=r_ts,json=rTs,proto3" json:"r_ts,omitempty"`
	// Types that are assignable to O:
	//	*Supported_Ob
	//	*Supported_Oi
	O isSupported_O `protobuf_oneof:"o"`
}

func (x *Supported) Reset() {
	*x = Supported{}
	if protoimpl.UnsafeEnabled {
		mi := &file_support_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Supported) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Supported) ProtoMessage() {}

func (x *Supported) ProtoReflect() protoreflect.Message {
	mi := &file_support_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Supported.ProtoReflect.Descriptor instead.
func (*Supported) Descriptor() ([]byte, []int) {
	return file_support_proto_rawDescGZIP(), []int{0}
}

func (x *Supported) GetBl() bool {
	if x != nil {
		return x.Bl
	}
	return false
}

func (x *Supported) GetI32() int32 {
	if x != nil {
		return x.I32
	}
	return 0
}

func (x *Supported) GetI64() int64 {
	if x != nil {
		return x.I64
	}
	return 0
}

func (x *Supported) GetF() float32 {
	if x != nil {
		return x.F
	}
	return 0
}

func (x *Supported) GetD() float64 {
	if x != nil {
		return x.D
	}
	return 0
}

func (x *Supported) GetS() string {
	if x != nil {
		return x.S
	}
	return ""
}

func (x *Supported) GetBt() []byte {
	if x != nil {
		return x.Bt
	}
	return nil
}

func (x *Supported) GetU32() uint32 {
	if x != nil {
		return x.U32
	}
	return 0
}

func (x *Supported) GetU64() uint64 {
	if x != nil {
		return x.U64
	}
	return 0
}

func (x *Supported) GetTs() *timestamppb.Timestamp {
	if x != nil {
		return x.Ts
	}
	return nil
}

func (x *Supported) GetRBl() []bool {
	if x != nil {
		return x.RBl
	}
	return nil
}

func (x *Supported) GetRI32() []int32 {
	if x != nil {
		return x.RI32
	}
	return nil
}

func (x *Supported) GetRI64() []int64 {
	if x != nil {
		return x.RI64
	}
	return nil
}

func (x *Supported) GetRF() []float32 {
	if x != nil {
		return x.RF
	}
	return nil
}

func (x *Supported) GetRD() []float64 {
	if x != nil {
		return x.RD
	}
	return nil
}

func (x *Supported) GetRS() []string {
	if x != nil {
		return x.RS
	}
	return nil
}

func (x *Supported) GetRU32() []uint32 {
	if x != nil {
		return x.RU32
	}
	return nil
}

func (x *Supported) GetRBt() [][]byte {
	if x != nil {
		return x.RBt
	}
	return nil
}

func (x *Supported) GetRU64() []uint64 {
	if x != nil {
		return x.RU64
	}
	return nil
}

func (x *Supported) GetRTs() []*timestamppb.Timestamp {
	if x != nil {
		return x.RTs
	}
	return nil
}

func (m *Supported) GetO() isSupported_O {
	if m != nil {
		return m.O
	}
	return nil
}

func (x *Supported) GetOb() bool {
	if x, ok := x.GetO().(*Supported_Ob); ok {
		return x.Ob
	}
	return false
}

func (x *Supported) GetOi() int32 {
	if x, ok := x.GetO().(*Supported_Oi); ok {
		return x.Oi
	}
	return 0
}

type isSupported_O interface {
	isSupported_O()
}

type Supported_Ob struct {
	Ob bool `protobuf:"varint,21,opt,name=ob,proto3,oneof"`
}

type Supported_Oi struct {
	Oi int32 `protobuf:"varint,22,opt,name=oi,proto3,oneof"`
}

func (*Supported_Ob) isSupported_O() {}

func (*Supported_Oi) isSupported_O() {}

// Unsupported scan destination types (for now)
type Unsupported struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sup  *Supported                       `protobuf:"bytes,1,opt,name=sup,proto3" json:"sup,omitempty"` // Nested, unregisterd messages
	Mp   map[string]int32                 `protobuf:"bytes,2,rep,name=mp,proto3" json:"mp,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	TsMp map[int32]*timestamppb.Timestamp `protobuf:"bytes,3,rep,name=ts_mp,json=tsMp,proto3" json:"ts_mp,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	En   SimpleColumns                    `protobuf:"varint,4,opt,name=en,proto3,enum=support.SimpleColumns" json:"en,omitempty"`
	REn  []SimpleColumns                  `protobuf:"varint,5,rep,packed,name=r_en,json=rEn,proto3,enum=support.SimpleColumns" json:"r_en,omitempty"`
}

func (x *Unsupported) Reset() {
	*x = Unsupported{}
	if protoimpl.UnsafeEnabled {
		mi := &file_support_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Unsupported) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Unsupported) ProtoMessage() {}

func (x *Unsupported) ProtoReflect() protoreflect.Message {
	mi := &file_support_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Unsupported.ProtoReflect.Descriptor instead.
func (*Unsupported) Descriptor() ([]byte, []int) {
	return file_support_proto_rawDescGZIP(), []int{1}
}

func (x *Unsupported) GetSup() *Supported {
	if x != nil {
		return x.Sup
	}
	return nil
}

func (x *Unsupported) GetMp() map[string]int32 {
	if x != nil {
		return x.Mp
	}
	return nil
}

func (x *Unsupported) GetTsMp() map[int32]*timestamppb.Timestamp {
	if x != nil {
		return x.TsMp
	}
	return nil
}

func (x *Unsupported) GetEn() SimpleColumns {
	if x != nil {
		return x.En
	}
	return SimpleColumns_id
}

func (x *Unsupported) GetREn() []SimpleColumns {
	if x != nil {
		return x.REn
	}
	return nil
}

// Simple is used for unit testing
type Simple struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title   string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Data    string                 `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Created *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created,proto3" json:"created,omitempty"`
}

func (x *Simple) Reset() {
	*x = Simple{}
	if protoimpl.UnsafeEnabled {
		mi := &file_support_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Simple) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Simple) ProtoMessage() {}

func (x *Simple) ProtoReflect() protoreflect.Message {
	mi := &file_support_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Simple.ProtoReflect.Descriptor instead.
func (*Simple) Descriptor() ([]byte, []int) {
	return file_support_proto_rawDescGZIP(), []int{2}
}

func (x *Simple) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Simple) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Simple) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

func (x *Simple) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

type SimpleQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      int32           `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Columns []SimpleColumns `protobuf:"varint,2,rep,packed,name=columns,proto3,enum=support.SimpleColumns" json:"columns,omitempty"`
}

func (x *SimpleQuery) Reset() {
	*x = SimpleQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_support_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SimpleQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SimpleQuery) ProtoMessage() {}

func (x *SimpleQuery) ProtoReflect() protoreflect.Message {
	mi := &file_support_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SimpleQuery.ProtoReflect.Descriptor instead.
func (*SimpleQuery) Descriptor() ([]byte, []int) {
	return file_support_proto_rawDescGZIP(), []int{3}
}

func (x *SimpleQuery) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SimpleQuery) GetColumns() []SimpleColumns {
	if x != nil {
		return x.Columns
	}
	return nil
}

var File_support_proto protoreflect.FileDescriptor

var file_support_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xce, 0x03, 0x0a, 0x09, 0x53, 0x75,
	0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x62, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x02, 0x62, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x33, 0x32, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x69, 0x33, 0x32, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x36, 0x34,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x69, 0x36, 0x34, 0x12, 0x0c, 0x0a, 0x01, 0x66,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x66, 0x12, 0x0c, 0x0a, 0x01, 0x64, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x73, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x01, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x62, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x02, 0x62, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x33, 0x32, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x03, 0x75, 0x33, 0x32, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x36, 0x34, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x75, 0x36, 0x34, 0x12, 0x2a, 0x0a, 0x02, 0x74, 0x73, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x02, 0x74, 0x73, 0x12, 0x11, 0x0a, 0x04, 0x72, 0x5f, 0x62, 0x6c, 0x18, 0x0b, 0x20,
	0x03, 0x28, 0x08, 0x52, 0x03, 0x72, 0x42, 0x6c, 0x12, 0x13, 0x0a, 0x05, 0x72, 0x5f, 0x69, 0x33,
	0x32, 0x18, 0x0c, 0x20, 0x03, 0x28, 0x05, 0x52, 0x04, 0x72, 0x49, 0x33, 0x32, 0x12, 0x13, 0x0a,
	0x05, 0x72, 0x5f, 0x69, 0x36, 0x34, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x03, 0x52, 0x04, 0x72, 0x49,
	0x36, 0x34, 0x12, 0x0f, 0x0a, 0x03, 0x72, 0x5f, 0x66, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x02, 0x52,
	0x02, 0x72, 0x46, 0x12, 0x0f, 0x0a, 0x03, 0x72, 0x5f, 0x64, 0x18, 0x0f, 0x20, 0x03, 0x28, 0x01,
	0x52, 0x02, 0x72, 0x44, 0x12, 0x0f, 0x0a, 0x03, 0x72, 0x5f, 0x73, 0x18, 0x10, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x02, 0x72, 0x53, 0x12, 0x13, 0x0a, 0x05, 0x72, 0x5f, 0x75, 0x33, 0x32, 0x18, 0x11,
	0x20, 0x03, 0x28, 0x0d, 0x52, 0x04, 0x72, 0x55, 0x33, 0x32, 0x12, 0x11, 0x0a, 0x04, 0x72, 0x5f,
	0x62, 0x74, 0x18, 0x12, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x03, 0x72, 0x42, 0x74, 0x12, 0x13, 0x0a,
	0x05, 0x72, 0x5f, 0x75, 0x36, 0x34, 0x18, 0x13, 0x20, 0x03, 0x28, 0x04, 0x52, 0x04, 0x72, 0x55,
	0x36, 0x34, 0x12, 0x2d, 0x0a, 0x04, 0x72, 0x5f, 0x74, 0x73, 0x18, 0x14, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x03, 0x72, 0x54,
	0x73, 0x12, 0x10, 0x0a, 0x02, 0x6f, 0x62, 0x18, 0x15, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52,
	0x02, 0x6f, 0x62, 0x12, 0x10, 0x0a, 0x02, 0x6f, 0x69, 0x18, 0x16, 0x20, 0x01, 0x28, 0x05, 0x48,
	0x00, 0x52, 0x02, 0x6f, 0x69, 0x42, 0x03, 0x0a, 0x01, 0x6f, 0x22, 0xf5, 0x02, 0x0a, 0x0b, 0x55,
	0x6e, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x12, 0x24, 0x0a, 0x03, 0x73, 0x75,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72,
	0x74, 0x2e, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x52, 0x03, 0x73, 0x75, 0x70,
	0x12, 0x2c, 0x0a, 0x02, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73,
	0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x55, 0x6e, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74,
	0x65, 0x64, 0x2e, 0x4d, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x02, 0x6d, 0x70, 0x12, 0x33,
	0x0a, 0x05, 0x74, 0x73, 0x5f, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x55, 0x6e, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72,
	0x74, 0x65, 0x64, 0x2e, 0x54, 0x73, 0x4d, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x74,
	0x73, 0x4d, 0x70, 0x12, 0x26, 0x0a, 0x02, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x16, 0x2e, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65,
	0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x52, 0x02, 0x65, 0x6e, 0x12, 0x29, 0x0a, 0x04, 0x72,
	0x5f, 0x65, 0x6e, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x73, 0x75, 0x70, 0x70,
	0x6f, 0x72, 0x74, 0x2e, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e,
	0x73, 0x52, 0x03, 0x72, 0x45, 0x6e, 0x1a, 0x35, 0x0a, 0x07, 0x4d, 0x70, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x53, 0x0a,
	0x09, 0x54, 0x73, 0x4d, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x30, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x78, 0x0a, 0x06, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x22, 0x4f, 0x0a, 0x0b,
	0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x30, 0x0a, 0x07, 0x63,
	0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x73,
	0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x43, 0x6f, 0x6c,
	0x75, 0x6d, 0x6e, 0x73, 0x52, 0x07, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x2a, 0x39, 0x0a,
	0x0d, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x12, 0x06,
	0x0a, 0x02, 0x69, 0x64, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x10,
	0x01, 0x12, 0x08, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x10, 0x03, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x75, 0x68, 0x6c, 0x65, 0x6d, 0x6d, 0x65, 0x72,
	0x2f, 0x70, 0x62, 0x70, 0x67, 0x78, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_support_proto_rawDescOnce sync.Once
	file_support_proto_rawDescData = file_support_proto_rawDesc
)

func file_support_proto_rawDescGZIP() []byte {
	file_support_proto_rawDescOnce.Do(func() {
		file_support_proto_rawDescData = protoimpl.X.CompressGZIP(file_support_proto_rawDescData)
	})
	return file_support_proto_rawDescData
}

var file_support_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_support_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_support_proto_goTypes = []interface{}{
	(SimpleColumns)(0),            // 0: support.SimpleColumns
	(*Supported)(nil),             // 1: support.Supported
	(*Unsupported)(nil),           // 2: support.Unsupported
	(*Simple)(nil),                // 3: support.Simple
	(*SimpleQuery)(nil),           // 4: support.SimpleQuery
	nil,                           // 5: support.Unsupported.MpEntry
	nil,                           // 6: support.Unsupported.TsMpEntry
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
}
var file_support_proto_depIdxs = []int32{
	7,  // 0: support.Supported.ts:type_name -> google.protobuf.Timestamp
	7,  // 1: support.Supported.r_ts:type_name -> google.protobuf.Timestamp
	1,  // 2: support.Unsupported.sup:type_name -> support.Supported
	5,  // 3: support.Unsupported.mp:type_name -> support.Unsupported.MpEntry
	6,  // 4: support.Unsupported.ts_mp:type_name -> support.Unsupported.TsMpEntry
	0,  // 5: support.Unsupported.en:type_name -> support.SimpleColumns
	0,  // 6: support.Unsupported.r_en:type_name -> support.SimpleColumns
	7,  // 7: support.Simple.created:type_name -> google.protobuf.Timestamp
	0,  // 8: support.SimpleQuery.columns:type_name -> support.SimpleColumns
	7,  // 9: support.Unsupported.TsMpEntry.value:type_name -> google.protobuf.Timestamp
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_support_proto_init() }
func file_support_proto_init() {
	if File_support_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_support_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Supported); i {
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
		file_support_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Unsupported); i {
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
		file_support_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Simple); i {
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
		file_support_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SimpleQuery); i {
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
	file_support_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Supported_Ob)(nil),
		(*Supported_Oi)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_support_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_support_proto_goTypes,
		DependencyIndexes: file_support_proto_depIdxs,
		EnumInfos:         file_support_proto_enumTypes,
		MessageInfos:      file_support_proto_msgTypes,
	}.Build()
	File_support_proto = out.File
	file_support_proto_rawDesc = nil
	file_support_proto_goTypes = nil
	file_support_proto_depIdxs = nil
}

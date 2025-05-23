// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.3
// source: preset.proto

package grpc

import (
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

// Channel type enumeration
type ChannelType int32

const (
	ChannelType_CHANNEL_TYPE_UNSPECIFIED ChannelType = 0
	ChannelType_CHANNEL_TYPE_GLOBAL      ChannelType = 1
	ChannelType_CHANNEL_TYPE_SAMPLER     ChannelType = 2
	ChannelType_CHANNEL_TYPE_INSTRUMENT  ChannelType = 3
	ChannelType_CHANNEL_TYPE_MIXER       ChannelType = 4
	ChannelType_CHANNEL_TYPE_PLAYER      ChannelType = 5
)

// Enum value maps for ChannelType.
var (
	ChannelType_name = map[int32]string{
		0: "CHANNEL_TYPE_UNSPECIFIED",
		1: "CHANNEL_TYPE_GLOBAL",
		2: "CHANNEL_TYPE_SAMPLER",
		3: "CHANNEL_TYPE_INSTRUMENT",
		4: "CHANNEL_TYPE_MIXER",
		5: "CHANNEL_TYPE_PLAYER",
	}
	ChannelType_value = map[string]int32{
		"CHANNEL_TYPE_UNSPECIFIED": 0,
		"CHANNEL_TYPE_GLOBAL":      1,
		"CHANNEL_TYPE_SAMPLER":     2,
		"CHANNEL_TYPE_INSTRUMENT":  3,
		"CHANNEL_TYPE_MIXER":       4,
		"CHANNEL_TYPE_PLAYER":      5,
	}
)

func (x ChannelType) Enum() *ChannelType {
	p := new(ChannelType)
	*p = x
	return p
}

func (x ChannelType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ChannelType) Descriptor() protoreflect.EnumDescriptor {
	return file_preset_proto_enumTypes[0].Descriptor()
}

func (ChannelType) Type() protoreflect.EnumType {
	return &file_preset_proto_enumTypes[0]
}

func (x ChannelType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ChannelType.Descriptor instead.
func (ChannelType) EnumDescriptor() ([]byte, []int) {
	return file_preset_proto_rawDescGZIP(), []int{0}
}

// FX parameter type enumeration
type FXParamType int32

const (
	FXParamType_FX_PARAM_TYPE_UNSPECIFIED FXParamType = 0
	FXParamType_FX_PARAM_TYPE_RANGE       FXParamType = 1
	FXParamType_FX_PARAM_TYPE_FIXED       FXParamType = 2
	FXParamType_FX_PARAM_TYPE_BOOLEAN     FXParamType = 3
)

// Enum value maps for FXParamType.
var (
	FXParamType_name = map[int32]string{
		0: "FX_PARAM_TYPE_UNSPECIFIED",
		1: "FX_PARAM_TYPE_RANGE",
		2: "FX_PARAM_TYPE_FIXED",
		3: "FX_PARAM_TYPE_BOOLEAN",
	}
	FXParamType_value = map[string]int32{
		"FX_PARAM_TYPE_UNSPECIFIED": 0,
		"FX_PARAM_TYPE_RANGE":       1,
		"FX_PARAM_TYPE_FIXED":       2,
		"FX_PARAM_TYPE_BOOLEAN":     3,
	}
)

func (x FXParamType) Enum() *FXParamType {
	p := new(FXParamType)
	*p = x
	return p
}

func (x FXParamType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FXParamType) Descriptor() protoreflect.EnumDescriptor {
	return file_preset_proto_enumTypes[1].Descriptor()
}

func (FXParamType) Type() protoreflect.EnumType {
	return &file_preset_proto_enumTypes[1]
}

func (x FXParamType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FXParamType.Descriptor instead.
func (FXParamType) EnumDescriptor() ([]byte, []int) {
	return file_preset_proto_rawDescGZIP(), []int{1}
}

// Request message for loading a preset
type LoadPresetRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PresetId      int64                  `protobuf:"varint,1,opt,name=preset_id,json=presetId,proto3" json:"preset_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoadPresetRequest) Reset() {
	*x = LoadPresetRequest{}
	mi := &file_preset_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoadPresetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadPresetRequest) ProtoMessage() {}

func (x *LoadPresetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_preset_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadPresetRequest.ProtoReflect.Descriptor instead.
func (*LoadPresetRequest) Descriptor() ([]byte, []int) {
	return file_preset_proto_rawDescGZIP(), []int{0}
}

func (x *LoadPresetRequest) GetPresetId() int64 {
	if x != nil {
		return x.PresetId
	}
	return 0
}

// Response message for loading a preset
type LoadPresetResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Preset        *Preset                `protobuf:"bytes,1,opt,name=preset,proto3" json:"preset,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoadPresetResponse) Reset() {
	*x = LoadPresetResponse{}
	mi := &file_preset_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoadPresetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadPresetResponse) ProtoMessage() {}

func (x *LoadPresetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_preset_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadPresetResponse.ProtoReflect.Descriptor instead.
func (*LoadPresetResponse) Descriptor() ([]byte, []int) {
	return file_preset_proto_rawDescGZIP(), []int{1}
}

func (x *LoadPresetResponse) GetPreset() *Preset {
	if x != nil {
		return x.Preset
	}
	return nil
}

// Preset message
type Preset struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   *string                `protobuf:"bytes,3,opt,name=description,proto3,oneof" json:"description,omitempty"`
	Channels      []*Channel             `protobuf:"bytes,4,rep,name=channels,proto3" json:"channels,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Preset) Reset() {
	*x = Preset{}
	mi := &file_preset_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Preset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Preset) ProtoMessage() {}

func (x *Preset) ProtoReflect() protoreflect.Message {
	mi := &file_preset_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Preset.ProtoReflect.Descriptor instead.
func (*Preset) Descriptor() ([]byte, []int) {
	return file_preset_proto_rawDescGZIP(), []int{2}
}

func (x *Preset) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Preset) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Preset) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *Preset) GetChannels() []*Channel {
	if x != nil {
		return x.Channels
	}
	return nil
}

// Channel message
type Channel struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Type          ChannelType            `protobuf:"varint,3,opt,name=type,proto3,enum=kitPreset.v1.ChannelType" json:"type,omitempty"`
	Volume        float64                `protobuf:"fixed64,4,opt,name=volume,proto3" json:"volume,omitempty"`
	Pan           float64                `protobuf:"fixed64,5,opt,name=pan,proto3" json:"pan,omitempty"`
	Fxs           []*FX                  `protobuf:"bytes,6,rep,name=fxs,proto3" json:"fxs,omitempty"`
	Instruments   []*Instrument          `protobuf:"bytes,7,rep,name=instruments,proto3" json:"instruments,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Channel) Reset() {
	*x = Channel{}
	mi := &file_preset_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Channel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Channel) ProtoMessage() {}

func (x *Channel) ProtoReflect() protoreflect.Message {
	mi := &file_preset_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Channel.ProtoReflect.Descriptor instead.
func (*Channel) Descriptor() ([]byte, []int) {
	return file_preset_proto_rawDescGZIP(), []int{3}
}

func (x *Channel) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Channel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Channel) GetType() ChannelType {
	if x != nil {
		return x.Type
	}
	return ChannelType_CHANNEL_TYPE_UNSPECIFIED
}

func (x *Channel) GetVolume() float64 {
	if x != nil {
		return x.Volume
	}
	return 0
}

func (x *Channel) GetPan() float64 {
	if x != nil {
		return x.Pan
	}
	return 0
}

func (x *Channel) GetFxs() []*FX {
	if x != nil {
		return x.Fxs
	}
	return nil
}

func (x *Channel) GetInstruments() []*Instrument {
	if x != nil {
		return x.Instruments
	}
	return nil
}

// Instrument message
type Instrument struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Volume        *float64               `protobuf:"fixed64,3,opt,name=volume,proto3,oneof" json:"volume,omitempty"`
	Pan           *float64               `protobuf:"fixed64,4,opt,name=pan,proto3,oneof" json:"pan,omitempty"`
	Tunes         []*FX                  `protobuf:"bytes,5,rep,name=tunes,proto3" json:"tunes,omitempty"`
	Layers        []*Layer               `protobuf:"bytes,6,rep,name=layers,proto3" json:"layers,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Instrument) Reset() {
	*x = Instrument{}
	mi := &file_preset_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Instrument) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Instrument) ProtoMessage() {}

func (x *Instrument) ProtoReflect() protoreflect.Message {
	mi := &file_preset_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Instrument.ProtoReflect.Descriptor instead.
func (*Instrument) Descriptor() ([]byte, []int) {
	return file_preset_proto_rawDescGZIP(), []int{4}
}

func (x *Instrument) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Instrument) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Instrument) GetVolume() float64 {
	if x != nil && x.Volume != nil {
		return *x.Volume
	}
	return 0
}

func (x *Instrument) GetPan() float64 {
	if x != nil && x.Pan != nil {
		return *x.Pan
	}
	return 0
}

func (x *Instrument) GetTunes() []*FX {
	if x != nil {
		return x.Tunes
	}
	return nil
}

func (x *Instrument) GetLayers() []*Layer {
	if x != nil {
		return x.Layers
	}
	return nil
}

// Layer message
type Layer struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Volume        *float64               `protobuf:"fixed64,3,opt,name=volume,proto3,oneof" json:"volume,omitempty"`
	Pan           *float64               `protobuf:"fixed64,4,opt,name=pan,proto3,oneof" json:"pan,omitempty"`
	Fxs           []*FX                  `protobuf:"bytes,5,rep,name=fxs,proto3" json:"fxs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Layer) Reset() {
	*x = Layer{}
	mi := &file_preset_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Layer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Layer) ProtoMessage() {}

func (x *Layer) ProtoReflect() protoreflect.Message {
	mi := &file_preset_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Layer.ProtoReflect.Descriptor instead.
func (*Layer) Descriptor() ([]byte, []int) {
	return file_preset_proto_rawDescGZIP(), []int{5}
}

func (x *Layer) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Layer) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Layer) GetVolume() float64 {
	if x != nil && x.Volume != nil {
		return *x.Volume
	}
	return 0
}

func (x *Layer) GetPan() float64 {
	if x != nil && x.Pan != nil {
		return *x.Pan
	}
	return 0
}

func (x *Layer) GetFxs() []*FX {
	if x != nil {
		return x.Fxs
	}
	return nil
}

// FX message
type FX struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Order         int32                  `protobuf:"varint,3,opt,name=order,proto3" json:"order,omitempty"`
	Params        []*FXParam             `protobuf:"bytes,4,rep,name=params,proto3" json:"params,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FX) Reset() {
	*x = FX{}
	mi := &file_preset_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FX) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FX) ProtoMessage() {}

func (x *FX) ProtoReflect() protoreflect.Message {
	mi := &file_preset_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FX.ProtoReflect.Descriptor instead.
func (*FX) Descriptor() ([]byte, []int) {
	return file_preset_proto_rawDescGZIP(), []int{6}
}

func (x *FX) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *FX) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FX) GetOrder() int32 {
	if x != nil {
		return x.Order
	}
	return 0
}

func (x *FX) GetParams() []*FXParam {
	if x != nil {
		return x.Params
	}
	return nil
}

// FX Parameter message
type FXParam struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Order         int32                  `protobuf:"varint,3,opt,name=order,proto3" json:"order,omitempty"`
	Type          FXParamType            `protobuf:"varint,4,opt,name=type,proto3,enum=kitPreset.v1.FXParamType" json:"type,omitempty"`
	Min           *float64               `protobuf:"fixed64,5,opt,name=min,proto3,oneof" json:"min,omitempty"`
	Max           *float64               `protobuf:"fixed64,6,opt,name=max,proto3,oneof" json:"max,omitempty"`
	Divisions     *int32                 `protobuf:"varint,7,opt,name=divisions,proto3,oneof" json:"divisions,omitempty"`
	DiscreteVals  []*FXParamDiscreteVal  `protobuf:"bytes,8,rep,name=discrete_vals,json=discreteVals,proto3" json:"discrete_vals,omitempty"`
	Value         float64                `protobuf:"fixed64,9,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FXParam) Reset() {
	*x = FXParam{}
	mi := &file_preset_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FXParam) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FXParam) ProtoMessage() {}

func (x *FXParam) ProtoReflect() protoreflect.Message {
	mi := &file_preset_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FXParam.ProtoReflect.Descriptor instead.
func (*FXParam) Descriptor() ([]byte, []int) {
	return file_preset_proto_rawDescGZIP(), []int{7}
}

func (x *FXParam) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *FXParam) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FXParam) GetOrder() int32 {
	if x != nil {
		return x.Order
	}
	return 0
}

func (x *FXParam) GetType() FXParamType {
	if x != nil {
		return x.Type
	}
	return FXParamType_FX_PARAM_TYPE_UNSPECIFIED
}

func (x *FXParam) GetMin() float64 {
	if x != nil && x.Min != nil {
		return *x.Min
	}
	return 0
}

func (x *FXParam) GetMax() float64 {
	if x != nil && x.Max != nil {
		return *x.Max
	}
	return 0
}

func (x *FXParam) GetDivisions() int32 {
	if x != nil && x.Divisions != nil {
		return *x.Divisions
	}
	return 0
}

func (x *FXParam) GetDiscreteVals() []*FXParamDiscreteVal {
	if x != nil {
		return x.DiscreteVals
	}
	return nil
}

func (x *FXParam) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// FX Parameter Discrete Value message
type FXParamDiscreteVal struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          *string                `protobuf:"bytes,1,opt,name=name,proto3,oneof" json:"name,omitempty"`
	Val           float64                `protobuf:"fixed64,2,opt,name=val,proto3" json:"val,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FXParamDiscreteVal) Reset() {
	*x = FXParamDiscreteVal{}
	mi := &file_preset_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FXParamDiscreteVal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FXParamDiscreteVal) ProtoMessage() {}

func (x *FXParamDiscreteVal) ProtoReflect() protoreflect.Message {
	mi := &file_preset_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FXParamDiscreteVal.ProtoReflect.Descriptor instead.
func (*FXParamDiscreteVal) Descriptor() ([]byte, []int) {
	return file_preset_proto_rawDescGZIP(), []int{8}
}

func (x *FXParamDiscreteVal) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *FXParamDiscreteVal) GetVal() float64 {
	if x != nil {
		return x.Val
	}
	return 0
}

var File_preset_proto protoreflect.FileDescriptor

var file_preset_proto_rawDesc = string([]byte{
	0x0a, 0x0c, 0x70, 0x72, 0x65, 0x73, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c,
	0x6b, 0x69, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x22, 0x30, 0x0a, 0x11,
	0x4c, 0x6f, 0x61, 0x64, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x72, 0x65, 0x73, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x72, 0x65, 0x73, 0x65, 0x74, 0x49, 0x64, 0x22, 0x42,
	0x0a, 0x12, 0x4c, 0x6f, 0x61, 0x64, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x06, 0x70, 0x72, 0x65, 0x73, 0x65, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6b, 0x69, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x52, 0x06, 0x70, 0x72, 0x65, 0x73,
	0x65, 0x74, 0x22, 0x98, 0x01, 0x0a, 0x06, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x31, 0x0a, 0x08, 0x63, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6b,
	0x69, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x52, 0x08, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x42, 0x0e, 0x0a,
	0x0c, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xe8, 0x01,
	0x0a, 0x07, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x2d, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e,
	0x6b, 0x69, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06,
	0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x61, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x03, 0x70, 0x61, 0x6e, 0x12, 0x22, 0x0a, 0x03, 0x66, 0x78, 0x73, 0x18,
	0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b, 0x69, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x58, 0x52, 0x03, 0x66, 0x78, 0x73, 0x12, 0x3a, 0x0a, 0x0b,
	0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x6b, 0x69, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0b, 0x69, 0x6e, 0x73,
	0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0xce, 0x01, 0x0a, 0x0a, 0x49, 0x6e, 0x73,
	0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a,
	0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x48, 0x00, 0x52,
	0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x70, 0x61,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x48, 0x01, 0x52, 0x03, 0x70, 0x61, 0x6e, 0x88, 0x01,
	0x01, 0x12, 0x26, 0x0a, 0x05, 0x74, 0x75, 0x6e, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x6b, 0x69, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x46, 0x58, 0x52, 0x05, 0x74, 0x75, 0x6e, 0x65, 0x73, 0x12, 0x2b, 0x0a, 0x06, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6b, 0x69, 0x74, 0x50,
	0x72, 0x65, 0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x06,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x76, 0x6f, 0x6c, 0x75, 0x6d,
	0x65, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x70, 0x61, 0x6e, 0x22, 0x98, 0x01, 0x0a, 0x05, 0x4c, 0x61,
	0x79, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x06, 0x76, 0x6f, 0x6c,
	0x75, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x06, 0x76, 0x6f, 0x6c,
	0x75, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x70, 0x61, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x01, 0x48, 0x01, 0x52, 0x03, 0x70, 0x61, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a,
	0x03, 0x66, 0x78, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6b, 0x69, 0x74,
	0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x58, 0x52, 0x03, 0x66, 0x78,
	0x73, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x42, 0x06, 0x0a, 0x04,
	0x5f, 0x70, 0x61, 0x6e, 0x22, 0x6f, 0x0a, 0x02, 0x46, 0x58, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6b, 0x69, 0x74, 0x50, 0x72, 0x65, 0x73,
	0x65, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x58, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x52, 0x06, 0x70,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x22, 0xc0, 0x02, 0x0a, 0x07, 0x46, 0x58, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x2d, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x6b, 0x69,
	0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x58, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x03,
	0x6d, 0x69, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x03, 0x6d, 0x69, 0x6e,
	0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01,
	0x48, 0x01, 0x52, 0x03, 0x6d, 0x61, 0x78, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x64, 0x69,
	0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x48, 0x02, 0x52,
	0x09, 0x64, 0x69, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x88, 0x01, 0x01, 0x12, 0x45, 0x0a,
	0x0d, 0x64, 0x69, 0x73, 0x63, 0x72, 0x65, 0x74, 0x65, 0x5f, 0x76, 0x61, 0x6c, 0x73, 0x18, 0x08,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6b, 0x69, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x46, 0x58, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x44, 0x69, 0x73, 0x63, 0x72,
	0x65, 0x74, 0x65, 0x56, 0x61, 0x6c, 0x52, 0x0c, 0x64, 0x69, 0x73, 0x63, 0x72, 0x65, 0x74, 0x65,
	0x56, 0x61, 0x6c, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x6d,
	0x69, 0x6e, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x6d, 0x61, 0x78, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x64,
	0x69, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x48, 0x0a, 0x12, 0x46, 0x58, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x44, 0x69, 0x73, 0x63, 0x72, 0x65, 0x74, 0x65, 0x56, 0x61, 0x6c, 0x12, 0x17,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x2a, 0xac, 0x01, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1c, 0x0a, 0x18, 0x43, 0x48, 0x41, 0x4e, 0x4e, 0x45, 0x4c, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x17, 0x0a, 0x13, 0x43, 0x48, 0x41, 0x4e, 0x4e, 0x45, 0x4c, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x47, 0x4c, 0x4f, 0x42, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x43, 0x48, 0x41,
	0x4e, 0x4e, 0x45, 0x4c, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x41, 0x4d, 0x50, 0x4c, 0x45,
	0x52, 0x10, 0x02, 0x12, 0x1b, 0x0a, 0x17, 0x43, 0x48, 0x41, 0x4e, 0x4e, 0x45, 0x4c, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x53, 0x54, 0x52, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x03,
	0x12, 0x16, 0x0a, 0x12, 0x43, 0x48, 0x41, 0x4e, 0x4e, 0x45, 0x4c, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x4d, 0x49, 0x58, 0x45, 0x52, 0x10, 0x04, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x48, 0x41, 0x4e,
	0x4e, 0x45, 0x4c, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x50, 0x4c, 0x41, 0x59, 0x45, 0x52, 0x10,
	0x05, 0x2a, 0x79, 0x0a, 0x0b, 0x46, 0x58, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x1d, 0x0a, 0x19, 0x46, 0x58, 0x5f, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x17, 0x0a, 0x13, 0x46, 0x58, 0x5f, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x52, 0x41, 0x4e, 0x47, 0x45, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x46, 0x58, 0x5f, 0x50,
	0x41, 0x52, 0x41, 0x4d, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x46, 0x49, 0x58, 0x45, 0x44, 0x10,
	0x02, 0x12, 0x19, 0x0a, 0x15, 0x46, 0x58, 0x5f, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x42, 0x4f, 0x4f, 0x4c, 0x45, 0x41, 0x4e, 0x10, 0x03, 0x32, 0x5c, 0x0a, 0x09,
	0x4b, 0x69, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x12, 0x4f, 0x0a, 0x0a, 0x4c, 0x6f, 0x61,
	0x64, 0x50, 0x72, 0x65, 0x73, 0x65, 0x74, 0x12, 0x1f, 0x2e, 0x6b, 0x69, 0x74, 0x50, 0x72, 0x65,
	0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x50, 0x72, 0x65, 0x73, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x6b, 0x69, 0x74, 0x50, 0x72,
	0x65, 0x73, 0x65, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x50, 0x72, 0x65, 0x73,
	0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x23, 0x5a, 0x21, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x73, 0x70, 0x69, 0x64, 0x72,
	0x75, 0x6d, 0x2d, 0x73, 0x72, 0x76, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_preset_proto_rawDescOnce sync.Once
	file_preset_proto_rawDescData []byte
)

func file_preset_proto_rawDescGZIP() []byte {
	file_preset_proto_rawDescOnce.Do(func() {
		file_preset_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_preset_proto_rawDesc), len(file_preset_proto_rawDesc)))
	})
	return file_preset_proto_rawDescData
}

var file_preset_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_preset_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_preset_proto_goTypes = []any{
	(ChannelType)(0),           // 0: kitPreset.v1.ChannelType
	(FXParamType)(0),           // 1: kitPreset.v1.FXParamType
	(*LoadPresetRequest)(nil),  // 2: kitPreset.v1.LoadPresetRequest
	(*LoadPresetResponse)(nil), // 3: kitPreset.v1.LoadPresetResponse
	(*Preset)(nil),             // 4: kitPreset.v1.Preset
	(*Channel)(nil),            // 5: kitPreset.v1.Channel
	(*Instrument)(nil),         // 6: kitPreset.v1.Instrument
	(*Layer)(nil),              // 7: kitPreset.v1.Layer
	(*FX)(nil),                 // 8: kitPreset.v1.FX
	(*FXParam)(nil),            // 9: kitPreset.v1.FXParam
	(*FXParamDiscreteVal)(nil), // 10: kitPreset.v1.FXParamDiscreteVal
}
var file_preset_proto_depIdxs = []int32{
	4,  // 0: kitPreset.v1.LoadPresetResponse.preset:type_name -> kitPreset.v1.Preset
	5,  // 1: kitPreset.v1.Preset.channels:type_name -> kitPreset.v1.Channel
	0,  // 2: kitPreset.v1.Channel.type:type_name -> kitPreset.v1.ChannelType
	8,  // 3: kitPreset.v1.Channel.fxs:type_name -> kitPreset.v1.FX
	6,  // 4: kitPreset.v1.Channel.instruments:type_name -> kitPreset.v1.Instrument
	8,  // 5: kitPreset.v1.Instrument.tunes:type_name -> kitPreset.v1.FX
	7,  // 6: kitPreset.v1.Instrument.layers:type_name -> kitPreset.v1.Layer
	8,  // 7: kitPreset.v1.Layer.fxs:type_name -> kitPreset.v1.FX
	9,  // 8: kitPreset.v1.FX.params:type_name -> kitPreset.v1.FXParam
	1,  // 9: kitPreset.v1.FXParam.type:type_name -> kitPreset.v1.FXParamType
	10, // 10: kitPreset.v1.FXParam.discrete_vals:type_name -> kitPreset.v1.FXParamDiscreteVal
	2,  // 11: kitPreset.v1.KitPreset.LoadPreset:input_type -> kitPreset.v1.LoadPresetRequest
	3,  // 12: kitPreset.v1.KitPreset.LoadPreset:output_type -> kitPreset.v1.LoadPresetResponse
	12, // [12:13] is the sub-list for method output_type
	11, // [11:12] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_preset_proto_init() }
func file_preset_proto_init() {
	if File_preset_proto != nil {
		return
	}
	file_preset_proto_msgTypes[2].OneofWrappers = []any{}
	file_preset_proto_msgTypes[4].OneofWrappers = []any{}
	file_preset_proto_msgTypes[5].OneofWrappers = []any{}
	file_preset_proto_msgTypes[7].OneofWrappers = []any{}
	file_preset_proto_msgTypes[8].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_preset_proto_rawDesc), len(file_preset_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_preset_proto_goTypes,
		DependencyIndexes: file_preset_proto_depIdxs,
		EnumInfos:         file_preset_proto_enumTypes,
		MessageInfos:      file_preset_proto_msgTypes,
	}.Build()
	File_preset_proto = out.File
	file_preset_proto_goTypes = nil
	file_preset_proto_depIdxs = nil
}

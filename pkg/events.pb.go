// Code generated by protoc-gen-go. DO NOT EDIT.
// source: events.proto

package pkg

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Direction int32

const (
	Direction_left  Direction = 0
	Direction_right Direction = 1
	Direction_up    Direction = 2
	Direction_down  Direction = 3
)

var Direction_name = map[int32]string{
	0: "left",
	1: "right",
	2: "up",
	3: "down",
}

var Direction_value = map[string]int32{
	"left":  0,
	"right": 1,
	"up":    2,
	"down":  3,
}

func (x Direction) String() string {
	return proto.EnumName(Direction_name, int32(x))
}

func (Direction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8f22242cb04491f9, []int{0}
}

type Event_Type int32

const (
	Event_type_init    Event_Type = 0
	Event_type_connect Event_Type = 1
	Event_type_exit    Event_Type = 2
	Event_type_idle    Event_Type = 3
	Event_type_move    Event_Type = 4
	Event_type_empty   Event_Type = 5
)

var Event_Type_name = map[int32]string{
	0: "type_init",
	1: "type_connect",
	2: "type_exit",
	3: "type_idle",
	4: "type_move",
	5: "type_empty",
}

var Event_Type_value = map[string]int32{
	"type_init":    0,
	"type_connect": 1,
	"type_exit":    2,
	"type_idle":    3,
	"type_move":    4,
	"type_empty":   5,
}

func (x Event_Type) String() string {
	return proto.EnumName(Event_Type_name, int32(x))
}

func (Event_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8f22242cb04491f9, []int{1, 0}
}

type Unit struct {
	Id                   string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	X                    float64   `protobuf:"fixed64,2,opt,name=x,proto3" json:"x,omitempty"`
	Y                    float64   `protobuf:"fixed64,3,opt,name=y,proto3" json:"y,omitempty"`
	Frame                int32     `protobuf:"varint,4,opt,name=frame,proto3" json:"frame,omitempty"`
	Skin                 string    `protobuf:"bytes,5,opt,name=skin,proto3" json:"skin,omitempty"`
	Action               string    `protobuf:"bytes,6,opt,name=action,proto3" json:"action,omitempty"`
	Speed                float64   `protobuf:"fixed64,7,opt,name=speed,proto3" json:"speed,omitempty"`
	Direction            Direction `protobuf:"varint,8,opt,name=direction,proto3,enum=tinyrpg.Direction" json:"direction,omitempty"`
	Side                 Direction `protobuf:"varint,9,opt,name=side,proto3,enum=tinyrpg.Direction" json:"side,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Unit) Reset()         { *m = Unit{} }
func (m *Unit) String() string { return proto.CompactTextString(m) }
func (*Unit) ProtoMessage()    {}
func (*Unit) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f22242cb04491f9, []int{0}
}

func (m *Unit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Unit.Unmarshal(m, b)
}
func (m *Unit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Unit.Marshal(b, m, deterministic)
}
func (m *Unit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Unit.Merge(m, src)
}
func (m *Unit) XXX_Size() int {
	return xxx_messageInfo_Unit.Size(m)
}
func (m *Unit) XXX_DiscardUnknown() {
	xxx_messageInfo_Unit.DiscardUnknown(m)
}

var xxx_messageInfo_Unit proto.InternalMessageInfo

func (m *Unit) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Unit) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Unit) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Unit) GetFrame() int32 {
	if m != nil {
		return m.Frame
	}
	return 0
}

func (m *Unit) GetSkin() string {
	if m != nil {
		return m.Skin
	}
	return ""
}

func (m *Unit) GetAction() string {
	if m != nil {
		return m.Action
	}
	return ""
}

func (m *Unit) GetSpeed() float64 {
	if m != nil {
		return m.Speed
	}
	return 0
}

func (m *Unit) GetDirection() Direction {
	if m != nil {
		return m.Direction
	}
	return Direction_left
}

func (m *Unit) GetSide() Direction {
	if m != nil {
		return m.Side
	}
	return Direction_left
}

type Event struct {
	Type Event_Type `protobuf:"varint,1,opt,name=type,proto3,enum=tinyrpg.Event_Type" json:"type,omitempty"`
	// Types that are valid to be assigned to Data:
	//	*Event_Init
	//	*Event_Connect
	//	*Event_Exit
	//	*Event_Idle
	//	*Event_Move
	Data                 isEvent_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f22242cb04491f9, []int{1}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetType() Event_Type {
	if m != nil {
		return m.Type
	}
	return Event_type_init
}

type isEvent_Data interface {
	isEvent_Data()
}

type Event_Init struct {
	Init *EventInit `protobuf:"bytes,2,opt,name=init,proto3,oneof"`
}

type Event_Connect struct {
	Connect *EventConnect `protobuf:"bytes,3,opt,name=connect,proto3,oneof"`
}

type Event_Exit struct {
	Exit *EventExit `protobuf:"bytes,4,opt,name=exit,proto3,oneof"`
}

type Event_Idle struct {
	Idle *EventIdle `protobuf:"bytes,5,opt,name=idle,proto3,oneof"`
}

type Event_Move struct {
	Move *EventMove `protobuf:"bytes,6,opt,name=move,proto3,oneof"`
}

func (*Event_Init) isEvent_Data() {}

func (*Event_Connect) isEvent_Data() {}

func (*Event_Exit) isEvent_Data() {}

func (*Event_Idle) isEvent_Data() {}

func (*Event_Move) isEvent_Data() {}

func (m *Event) GetData() isEvent_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Event) GetInit() *EventInit {
	if x, ok := m.GetData().(*Event_Init); ok {
		return x.Init
	}
	return nil
}

func (m *Event) GetConnect() *EventConnect {
	if x, ok := m.GetData().(*Event_Connect); ok {
		return x.Connect
	}
	return nil
}

func (m *Event) GetExit() *EventExit {
	if x, ok := m.GetData().(*Event_Exit); ok {
		return x.Exit
	}
	return nil
}

func (m *Event) GetIdle() *EventIdle {
	if x, ok := m.GetData().(*Event_Idle); ok {
		return x.Idle
	}
	return nil
}

func (m *Event) GetMove() *EventMove {
	if x, ok := m.GetData().(*Event_Move); ok {
		return x.Move
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Event) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Event_Init)(nil),
		(*Event_Connect)(nil),
		(*Event_Exit)(nil),
		(*Event_Idle)(nil),
		(*Event_Move)(nil),
	}
}

type EventInit struct {
	PlayerId             string           `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	Units                map[string]*Unit `protobuf:"bytes,2,rep,name=units,proto3" json:"units,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *EventInit) Reset()         { *m = EventInit{} }
func (m *EventInit) String() string { return proto.CompactTextString(m) }
func (*EventInit) ProtoMessage()    {}
func (*EventInit) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f22242cb04491f9, []int{2}
}

func (m *EventInit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventInit.Unmarshal(m, b)
}
func (m *EventInit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventInit.Marshal(b, m, deterministic)
}
func (m *EventInit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventInit.Merge(m, src)
}
func (m *EventInit) XXX_Size() int {
	return xxx_messageInfo_EventInit.Size(m)
}
func (m *EventInit) XXX_DiscardUnknown() {
	xxx_messageInfo_EventInit.DiscardUnknown(m)
}

var xxx_messageInfo_EventInit proto.InternalMessageInfo

func (m *EventInit) GetPlayerId() string {
	if m != nil {
		return m.PlayerId
	}
	return ""
}

func (m *EventInit) GetUnits() map[string]*Unit {
	if m != nil {
		return m.Units
	}
	return nil
}

type EventConnect struct {
	Unit                 *Unit    `protobuf:"bytes,1,opt,name=unit,proto3" json:"unit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventConnect) Reset()         { *m = EventConnect{} }
func (m *EventConnect) String() string { return proto.CompactTextString(m) }
func (*EventConnect) ProtoMessage()    {}
func (*EventConnect) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f22242cb04491f9, []int{3}
}

func (m *EventConnect) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventConnect.Unmarshal(m, b)
}
func (m *EventConnect) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventConnect.Marshal(b, m, deterministic)
}
func (m *EventConnect) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventConnect.Merge(m, src)
}
func (m *EventConnect) XXX_Size() int {
	return xxx_messageInfo_EventConnect.Size(m)
}
func (m *EventConnect) XXX_DiscardUnknown() {
	xxx_messageInfo_EventConnect.DiscardUnknown(m)
}

var xxx_messageInfo_EventConnect proto.InternalMessageInfo

func (m *EventConnect) GetUnit() *Unit {
	if m != nil {
		return m.Unit
	}
	return nil
}

type EventExit struct {
	PlayerId             string   `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventExit) Reset()         { *m = EventExit{} }
func (m *EventExit) String() string { return proto.CompactTextString(m) }
func (*EventExit) ProtoMessage()    {}
func (*EventExit) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f22242cb04491f9, []int{4}
}

func (m *EventExit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventExit.Unmarshal(m, b)
}
func (m *EventExit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventExit.Marshal(b, m, deterministic)
}
func (m *EventExit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventExit.Merge(m, src)
}
func (m *EventExit) XXX_Size() int {
	return xxx_messageInfo_EventExit.Size(m)
}
func (m *EventExit) XXX_DiscardUnknown() {
	xxx_messageInfo_EventExit.DiscardUnknown(m)
}

var xxx_messageInfo_EventExit proto.InternalMessageInfo

func (m *EventExit) GetPlayerId() string {
	if m != nil {
		return m.PlayerId
	}
	return ""
}

type EventIdle struct {
	PlayerId             string   `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventIdle) Reset()         { *m = EventIdle{} }
func (m *EventIdle) String() string { return proto.CompactTextString(m) }
func (*EventIdle) ProtoMessage()    {}
func (*EventIdle) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f22242cb04491f9, []int{5}
}

func (m *EventIdle) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventIdle.Unmarshal(m, b)
}
func (m *EventIdle) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventIdle.Marshal(b, m, deterministic)
}
func (m *EventIdle) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventIdle.Merge(m, src)
}
func (m *EventIdle) XXX_Size() int {
	return xxx_messageInfo_EventIdle.Size(m)
}
func (m *EventIdle) XXX_DiscardUnknown() {
	xxx_messageInfo_EventIdle.DiscardUnknown(m)
}

var xxx_messageInfo_EventIdle proto.InternalMessageInfo

func (m *EventIdle) GetPlayerId() string {
	if m != nil {
		return m.PlayerId
	}
	return ""
}

type EventMove struct {
	PlayerId             string    `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	Direction            Direction `protobuf:"varint,2,opt,name=direction,proto3,enum=tinyrpg.Direction" json:"direction,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *EventMove) Reset()         { *m = EventMove{} }
func (m *EventMove) String() string { return proto.CompactTextString(m) }
func (*EventMove) ProtoMessage()    {}
func (*EventMove) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f22242cb04491f9, []int{6}
}

func (m *EventMove) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventMove.Unmarshal(m, b)
}
func (m *EventMove) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventMove.Marshal(b, m, deterministic)
}
func (m *EventMove) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventMove.Merge(m, src)
}
func (m *EventMove) XXX_Size() int {
	return xxx_messageInfo_EventMove.Size(m)
}
func (m *EventMove) XXX_DiscardUnknown() {
	xxx_messageInfo_EventMove.DiscardUnknown(m)
}

var xxx_messageInfo_EventMove proto.InternalMessageInfo

func (m *EventMove) GetPlayerId() string {
	if m != nil {
		return m.PlayerId
	}
	return ""
}

func (m *EventMove) GetDirection() Direction {
	if m != nil {
		return m.Direction
	}
	return Direction_left
}

func init() {
	proto.RegisterEnum("tinyrpg.Direction", Direction_name, Direction_value)
	proto.RegisterEnum("tinyrpg.Event_Type", Event_Type_name, Event_Type_value)
	proto.RegisterType((*Unit)(nil), "tinyrpg.Unit")
	proto.RegisterType((*Event)(nil), "tinyrpg.Event")
	proto.RegisterType((*EventInit)(nil), "tinyrpg.EventInit")
	proto.RegisterMapType((map[string]*Unit)(nil), "tinyrpg.EventInit.UnitsEntry")
	proto.RegisterType((*EventConnect)(nil), "tinyrpg.EventConnect")
	proto.RegisterType((*EventExit)(nil), "tinyrpg.EventExit")
	proto.RegisterType((*EventIdle)(nil), "tinyrpg.EventIdle")
	proto.RegisterType((*EventMove)(nil), "tinyrpg.EventMove")
}

func init() { proto.RegisterFile("events.proto", fileDescriptor_8f22242cb04491f9) }

var fileDescriptor_8f22242cb04491f9 = []byte{
	// 530 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x94, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xb3, 0x8e, 0x9d, 0xc4, 0x93, 0x34, 0xb2, 0x86, 0x3f, 0xb2, 0x40, 0x48, 0xc1, 0x48,
	0x60, 0x71, 0x88, 0x68, 0x7a, 0x41, 0x1c, 0x81, 0x88, 0xf6, 0xc0, 0xc5, 0x82, 0x0b, 0x97, 0xca,
	0x64, 0xb7, 0x65, 0x55, 0x67, 0x6d, 0xd9, 0x9b, 0x10, 0xbf, 0x10, 0x4f, 0xc6, 0x4b, 0x70, 0x43,
	0x33, 0x76, 0x9d, 0x56, 0x11, 0x11, 0x37, 0xcf, 0xcc, 0x6f, 0xbe, 0xcc, 0x7c, 0x63, 0x07, 0x26,
	0x6a, 0xab, 0x8c, 0xad, 0xe6, 0x45, 0x99, 0xdb, 0x1c, 0x87, 0x56, 0x9b, 0xba, 0x2c, 0xae, 0xa3,
	0xdf, 0x02, 0xdc, 0xaf, 0x46, 0x5b, 0x9c, 0x82, 0xa3, 0x65, 0x28, 0x66, 0x22, 0xf6, 0x13, 0x47,
	0x4b, 0x9c, 0x80, 0xd8, 0x85, 0xce, 0x4c, 0xc4, 0x22, 0x11, 0x3b, 0x8a, 0xea, 0xb0, 0xdf, 0x44,
	0x35, 0x3e, 0x04, 0xef, 0xaa, 0x4c, 0xd7, 0x2a, 0x74, 0x67, 0x22, 0xf6, 0x92, 0x26, 0x40, 0x04,
	0xb7, 0xba, 0xd1, 0x26, 0xf4, 0x58, 0x83, 0x9f, 0xf1, 0x31, 0x0c, 0xd2, 0x95, 0xd5, 0xb9, 0x09,
	0x07, 0x9c, 0x6d, 0x23, 0x52, 0xa8, 0x0a, 0xa5, 0x64, 0x38, 0x64, 0xcd, 0x26, 0xc0, 0x37, 0xe0,
	0x4b, 0x5d, 0xaa, 0xa6, 0x61, 0x34, 0x13, 0xf1, 0x74, 0x81, 0xf3, 0x76, 0xd2, 0xf9, 0xc7, 0xdb,
	0x4a, 0xb2, 0x87, 0xf0, 0x25, 0xb8, 0x95, 0x96, 0x2a, 0xf4, 0xff, 0x09, 0x73, 0x3d, 0xfa, 0xe3,
	0x80, 0xb7, 0x24, 0x03, 0xf0, 0x15, 0xb8, 0xb6, 0x2e, 0x14, 0x6f, 0x3a, 0x5d, 0x3c, 0xe8, 0x3a,
	0xb8, 0x3a, 0xff, 0x52, 0x17, 0x2a, 0x61, 0x00, 0x63, 0x70, 0xb5, 0xd1, 0x96, 0x3d, 0x18, 0xdf,
	0x91, 0x66, 0xf0, 0xc2, 0x68, 0x7b, 0xde, 0x4b, 0x98, 0xc0, 0x53, 0x18, 0xae, 0x72, 0x63, 0xd4,
	0xca, 0xb2, 0x45, 0xe3, 0xc5, 0xa3, 0xfb, 0xf0, 0x87, 0xa6, 0x78, 0xde, 0x4b, 0x6e, 0x39, 0x12,
	0x57, 0x3b, 0x6d, 0xd9, 0xc0, 0x03, 0xf1, 0xe5, 0xae, 0x11, 0x27, 0x82, 0xc7, 0x90, 0x99, 0x62,
	0x57, 0x0f, 0xc7, 0x90, 0x99, 0xe2, 0x31, 0x64, 0xc6, 0x03, 0xaf, 0xf3, 0xad, 0x62, 0xa7, 0x0f,
	0xc8, 0xcf, 0xf9, 0x96, 0x49, 0x22, 0x22, 0x09, 0x2e, 0x2d, 0x8a, 0x27, 0xe0, 0xd3, 0xaa, 0x97,
	0xb4, 0x45, 0xd0, 0xc3, 0x00, 0x26, 0x1c, 0xb6, 0x43, 0x06, 0xa2, 0x03, 0x68, 0x92, 0xc0, 0xd9,
	0xf3, 0x32, 0x53, 0x41, 0xbf, 0x0b, 0x49, 0x33, 0x70, 0x71, 0x0a, 0xd0, 0xc0, 0xeb, 0xc2, 0xd6,
	0x81, 0xf7, 0x7e, 0x00, 0xae, 0x4c, 0x6d, 0x1a, 0xfd, 0x12, 0xe0, 0x77, 0xa6, 0xe1, 0x53, 0xf0,
	0x8b, 0x2c, 0xad, 0x55, 0x79, 0xd9, 0xbd, 0x6e, 0xa3, 0x26, 0x71, 0x21, 0xf1, 0x0c, 0xbc, 0x8d,
	0xd1, 0xb6, 0x0a, 0x9d, 0x59, 0x3f, 0x1e, 0x2f, 0x9e, 0x1d, 0x9a, 0x3e, 0xa7, 0x97, 0xb5, 0x5a,
	0x1a, 0x5b, 0xd6, 0x49, 0xc3, 0x3e, 0xf9, 0x04, 0xb0, 0x4f, 0x62, 0x00, 0xfd, 0x1b, 0x55, 0xb7,
	0xca, 0xf4, 0x88, 0x2f, 0xc0, 0xdb, 0xa6, 0xd9, 0x46, 0xb5, 0x97, 0x3c, 0xe9, 0x44, 0xa9, 0x2b,
	0x69, 0x6a, 0xef, 0x9c, 0xb7, 0x22, 0x3a, 0x85, 0xc9, 0xdd, 0x7b, 0xe1, 0x73, 0x70, 0xe9, 0x17,
	0x58, 0xeb, 0xa0, 0x8f, 0x4b, 0x51, 0xdc, 0xae, 0x46, 0x27, 0x3b, 0xba, 0x5a, 0x47, 0xd2, 0xc9,
	0x8e, 0x93, 0xdf, 0x5a, 0x92, 0x4e, 0x76, 0xdc, 0xae, 0x7b, 0xdf, 0x8b, 0xf3, 0x1f, 0xdf, 0xcb,
	0xeb, 0x05, 0xf8, 0x5d, 0x1e, 0x47, 0xe0, 0x66, 0xea, 0x8a, 0x2e, 0xef, 0x83, 0x57, 0xea, 0xeb,
	0x1f, 0x74, 0xf2, 0x01, 0x38, 0x9b, 0x22, 0x70, 0xa8, 0x28, 0xf3, 0x9f, 0x26, 0xe8, 0x7f, 0x1f,
	0xf0, 0x5f, 0xc6, 0xd9, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xaf, 0x29, 0x0a, 0x73, 0x42, 0x04,
	0x00, 0x00,
}

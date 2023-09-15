// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: kademlia.proto

package protobuf

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

type KademliaMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageType int32  `protobuf:"varint,1,opt,name=message_type,json=messageType,proto3" json:"message_type,omitempty"`
	Data        []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *KademliaMessage) Reset() {
	*x = KademliaMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kademlia_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KademliaMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KademliaMessage) ProtoMessage() {}

func (x *KademliaMessage) ProtoReflect() protoreflect.Message {
	mi := &file_kademlia_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KademliaMessage.ProtoReflect.Descriptor instead.
func (*KademliaMessage) Descriptor() ([]byte, []int) {
	return file_kademlia_proto_rawDescGZIP(), []int{0}
}

func (x *KademliaMessage) GetMessageType() int32 {
	if x != nil {
		return x.MessageType
	}
	return 0
}

func (x *KademliaMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Ping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender *Node  `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	RpcID  string `protobuf:"bytes,2,opt,name=rpcID,proto3" json:"rpcID,omitempty"`
}

func (x *Ping) Reset() {
	*x = Ping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kademlia_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ping) ProtoMessage() {}

func (x *Ping) ProtoReflect() protoreflect.Message {
	mi := &file_kademlia_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ping.ProtoReflect.Descriptor instead.
func (*Ping) Descriptor() ([]byte, []int) {
	return file_kademlia_proto_rawDescGZIP(), []int{1}
}

func (x *Ping) GetSender() *Node {
	if x != nil {
		return x.Sender
	}
	return nil
}

func (x *Ping) GetRpcID() string {
	if x != nil {
		return x.RpcID
	}
	return ""
}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kademlia_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_kademlia_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_kademlia_proto_rawDescGZIP(), []int{2}
}

func (x *Node) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Node) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

var File_kademlia_proto protoreflect.FileDescriptor

var file_kademlia_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6b, 0x61, 0x64, 0x65, 0x6d, 0x6c, 0x69, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x48, 0x0a, 0x0f, 0x4b, 0x61,
	0x64, 0x65, 0x6d, 0x6c, 0x69, 0x61, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x21, 0x0a,
	0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x44, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x26, 0x0a, 0x06,
	0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x06, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x70, 0x63, 0x49, 0x44, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x70, 0x63, 0x49, 0x44, 0x22, 0x30, 0x0a, 0x04, 0x4e, 0x6f,
	0x64, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x42, 0x0d, 0x5a, 0x0b,
	0x2e, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_kademlia_proto_rawDescOnce sync.Once
	file_kademlia_proto_rawDescData = file_kademlia_proto_rawDesc
)

func file_kademlia_proto_rawDescGZIP() []byte {
	file_kademlia_proto_rawDescOnce.Do(func() {
		file_kademlia_proto_rawDescData = protoimpl.X.CompressGZIP(file_kademlia_proto_rawDescData)
	})
	return file_kademlia_proto_rawDescData
}

var file_kademlia_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_kademlia_proto_goTypes = []interface{}{
	(*KademliaMessage)(nil), // 0: protobuf.KademliaMessage
	(*Ping)(nil),            // 1: protobuf.Ping
	(*Node)(nil),            // 2: protobuf.Node
}
var file_kademlia_proto_depIdxs = []int32{
	2, // 0: protobuf.Ping.sender:type_name -> protobuf.Node
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_kademlia_proto_init() }
func file_kademlia_proto_init() {
	if File_kademlia_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kademlia_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KademliaMessage); i {
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
		file_kademlia_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ping); i {
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
		file_kademlia_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
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
			RawDescriptor: file_kademlia_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_kademlia_proto_goTypes,
		DependencyIndexes: file_kademlia_proto_depIdxs,
		MessageInfos:      file_kademlia_proto_msgTypes,
	}.Build()
	File_kademlia_proto = out.File
	file_kademlia_proto_rawDesc = nil
	file_kademlia_proto_goTypes = nil
	file_kademlia_proto_depIdxs = nil
}

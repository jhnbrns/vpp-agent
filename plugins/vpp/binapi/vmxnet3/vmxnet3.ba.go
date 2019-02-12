// Code generated by GoVPP binapi-generator. DO NOT EDIT.
//  source: /usr/share/vpp/api/vmxnet3.api.json

/*
 Package vmxnet3 is a generated from VPP binary API module 'vmxnet3'.

 It contains following objects:
	  3 services
	  6 messages
*/
package vmxnet3

import api "git.fd.io/govpp.git/api"
import struc "github.com/lunixbochs/struc"
import bytes "bytes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = api.RegisterMessage
var _ = struc.Pack
var _ = bytes.NewBuffer

// Services represents VPP binary API services:
type Services interface {
	DumpVmxnet3(*Vmxnet3Dump) ([]*Vmxnet3Details, error)
	Vmxnet3Create(*Vmxnet3Create) (*Vmxnet3CreateReply, error)
	Vmxnet3Delete(*Vmxnet3Delete) (*Vmxnet3DeleteReply, error)
}

/* Messages */

// Vmxnet3Create represents VPP binary API message 'vmxnet3_create':
type Vmxnet3Create struct {
	PciAddr    uint32
	EnableElog int32
	RxqSize    uint16
	TxqSize    uint16
}

func (*Vmxnet3Create) GetMessageName() string {
	return "vmxnet3_create"
}
func (*Vmxnet3Create) GetCrcString() string {
	return "7318251d"
}
func (*Vmxnet3Create) GetMessageType() api.MessageType {
	return api.RequestMessage
}

// Vmxnet3CreateReply represents VPP binary API message 'vmxnet3_create_reply':
type Vmxnet3CreateReply struct {
	Retval    int32
	SwIfIndex uint32
}

func (*Vmxnet3CreateReply) GetMessageName() string {
	return "vmxnet3_create_reply"
}
func (*Vmxnet3CreateReply) GetCrcString() string {
	return "fda5941f"
}
func (*Vmxnet3CreateReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}

// Vmxnet3Delete represents VPP binary API message 'vmxnet3_delete':
type Vmxnet3Delete struct {
	SwIfIndex uint32
}

func (*Vmxnet3Delete) GetMessageName() string {
	return "vmxnet3_delete"
}
func (*Vmxnet3Delete) GetCrcString() string {
	return "529cb13f"
}
func (*Vmxnet3Delete) GetMessageType() api.MessageType {
	return api.RequestMessage
}

// Vmxnet3DeleteReply represents VPP binary API message 'vmxnet3_delete_reply':
type Vmxnet3DeleteReply struct {
	Retval int32
}

func (*Vmxnet3DeleteReply) GetMessageName() string {
	return "vmxnet3_delete_reply"
}
func (*Vmxnet3DeleteReply) GetCrcString() string {
	return "e8d4e804"
}
func (*Vmxnet3DeleteReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}

// Vmxnet3Details represents VPP binary API message 'vmxnet3_details':
type Vmxnet3Details struct {
	SwIfIndex   uint32
	IfName      []byte `struc:"[64]byte"`
	HwAddr      []byte `struc:"[6]byte"`
	PciAddr     uint32
	Version     uint8
	RxQid       uint16
	RxQsize     uint16
	RxFill      []uint16 `struc:"[2]uint16"`
	RxNext      uint16
	RxProduce   []uint16 `struc:"[2]uint16"`
	RxConsume   []uint16 `struc:"[2]uint16"`
	TxQid       uint16
	TxQsize     uint16
	TxNext      uint16
	TxProduce   uint16
	TxConsume   uint16
	AdminUpDown uint8
}

func (*Vmxnet3Details) GetMessageName() string {
	return "vmxnet3_details"
}
func (*Vmxnet3Details) GetCrcString() string {
	return "2374ddc9"
}
func (*Vmxnet3Details) GetMessageType() api.MessageType {
	return api.ReplyMessage
}

// Vmxnet3Dump represents VPP binary API message 'vmxnet3_dump':
type Vmxnet3Dump struct{}

func (*Vmxnet3Dump) GetMessageName() string {
	return "vmxnet3_dump"
}
func (*Vmxnet3Dump) GetCrcString() string {
	return "51077d14"
}
func (*Vmxnet3Dump) GetMessageType() api.MessageType {
	return api.RequestMessage
}

func init() {
	api.RegisterMessage((*Vmxnet3Create)(nil), "vmxnet3.Vmxnet3Create")
	api.RegisterMessage((*Vmxnet3CreateReply)(nil), "vmxnet3.Vmxnet3CreateReply")
	api.RegisterMessage((*Vmxnet3Delete)(nil), "vmxnet3.Vmxnet3Delete")
	api.RegisterMessage((*Vmxnet3DeleteReply)(nil), "vmxnet3.Vmxnet3DeleteReply")
	api.RegisterMessage((*Vmxnet3Details)(nil), "vmxnet3.Vmxnet3Details")
	api.RegisterMessage((*Vmxnet3Dump)(nil), "vmxnet3.Vmxnet3Dump")
}

var Messages = []api.Message{
	(*Vmxnet3Create)(nil),
	(*Vmxnet3CreateReply)(nil),
	(*Vmxnet3Delete)(nil),
	(*Vmxnet3DeleteReply)(nil),
	(*Vmxnet3Details)(nil),
	(*Vmxnet3Dump)(nil),
}

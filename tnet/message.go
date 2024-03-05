package tnet

import "github.com/cpf2021-gif/gos/tiface"

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

// Message implements tiface.IMessage
var _ tiface.IMessage = (*Message)(nil)

func NewMsgPack(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

func (m *Message) SetDataLen(dlen uint32) {
	m.DataLen = dlen
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

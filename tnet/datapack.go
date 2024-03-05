package tnet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/cpf2021-gif/gos/tiface"
	"github.com/cpf2021-gif/gos/utils"
)

type DataPack struct{}

// DataPack implements tiface.IDataPack
var _ tiface.IDataPack = (*DataPack)(nil)

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen(4字节）+ ID(4字节)
	return 8
}

func (dp *DataPack) Pack(msg tiface.IMessage) ([]byte, error) {
	databuff := bytes.NewBuffer([]byte{})

	// write msg dataLen
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// write msg dataId
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// write msg data
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return databuff.Bytes(), nil
}

func (dp *DataPack) UnPack(data []byte) (tiface.IMessage, error) {
	read := bytes.NewReader(data)

	msg := &Message{}

	// read msg dataLen
	if err := binary.Read(read, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// read msg dataId
	if err := binary.Read(read, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if utils.GlobalConfig.Gos.MaxPacketSize > 0 && msg.GetMsgLen() > utils.GlobalConfig.Gos.MaxPacketSize {
		return nil, errors.New("too large msg data recv")
	}

	return msg, nil
}

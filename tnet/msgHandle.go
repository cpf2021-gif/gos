package tnet

import (
	"fmt"
	"strconv"

	"github.com/cpf2021-gif/gos/tiface"
)

type MsgHandle struct {
	Apis map[uint32]tiface.IRouter
}

// MsgHandle implements tiface.IMsgHandle
var _ tiface.IMsgHandle = (*MsgHandle)(nil)

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]tiface.IRouter),
	}
}

func (mh *MsgHandle) DoMsgHandler(request tiface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("not found router, msgId =", strconv.Itoa(int(request.GetMsgID())))
		return
	}

	handler.Do(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router tiface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api, msgId = " + strconv.Itoa(int(msgID)))
	}

	mh.Apis[msgID] = router
}

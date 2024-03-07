package tnet

import (
	"fmt"
	"strconv"

	"github.com/cpf2021-gif/gos/tiface"
	"github.com/cpf2021-gif/gos/utils"
)

type MsgHandle struct {
	// 存放每个MsgID 所对应的处理方法的map属性
	Apis map[uint32]tiface.IRouter
	// 负责Worker取任务的消息队列
	TaskQueue []chan tiface.IRequest
	// 业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

// MsgHandle implements tiface.IMsgHandle
var _ tiface.IMsgHandle = (*MsgHandle)(nil)

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]tiface.IRouter),
		WorkerPoolSize: utils.GlobalConfig.Gos.WorkerPoolSize,
		TaskQueue:      make([]chan tiface.IRequest, utils.GlobalConfig.Gos.WorkerPoolSize),
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

// 启动一个Worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	for i := range int(utils.GlobalConfig.Gos.WorkerPoolSize) {
		mh.TaskQueue[i] = make(chan tiface.IRequest, utils.GlobalConfig.Gos.MaxWorkerTaskLen)
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) startOneWorker(workerID int, taskQueue chan tiface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started ...")

	for {
		req := <-taskQueue
		mh.DoMsgHandler(req)
	}
}

// 将消息交给TaskQueue, 由Worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request tiface.IRequest) {
	// 1. 保证消息平均分配到不同的Worker
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), " request MsgID = ", request.GetMsgID(), " to WorkerID = ", workerID)

	// 2. 将消息发送给对应的Worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}

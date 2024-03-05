package tiface

import "net"

type IConnection interface {
	// 启动连接 让当前的连接准备开始工作
	Start()

	// 停止连接 结束当前连接的工作
	Stop()

	// 获取当前连接绑定的socket Conn
	GetTCPConnection() *net.TCPConn

	// 获取当前连接ID
	GetConnID() uint32

	// 获取远程客户端的状态 IP port
	RemoteAddr() net.Addr

	// 发送数据
	SendMsg(msgId uint32, data []byte) error
}

// 处理业务
type HandleFunc func(*net.TCPConn, []byte, int) error

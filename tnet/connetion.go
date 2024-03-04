package tnet

import (
	"fmt"
	"io"
	"net"

	"github.com/cpf2021-gif/gos/tiface"
	"github.com/cpf2021-gif/gos/utils"
)

type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnID uint32

	// 当前连接的状态
	isClosed bool

	// 该连接处理的方法Router
	Router tiface.IRouter

	// 告知当前连接已经退出/停止的channel
	ExitChan chan bool
}

// Connection implements tiface.IConnection
var _ tiface.IConnection = (*Connection)(nil)

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router tiface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}

	return c
}

// StartReader 启动连接的读数据业务
func (c *Connection) startReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中
		buf := make([]byte, utils.GlobalConfig.Gos.MaxPacketSize)
		_, err := c.Conn.Read(buf)

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("recv buf err ", err)
			continue
		}

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}

		// 执行注册路由
		go func(request tiface.IRequest) {
			c.Router.PreHandle(&req)
			c.Router.Handle(&req)
			c.Router.PostHandle(&req)
		}(&req)
	}
}

// Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)

	// 启动从当前连接的读数据的业务
	go c.startReader()
	// TODO 启动从当前连接读数据的业务
}

// Stop 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = true

	// 关闭socket连接
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true
	close(c.ExitChan)

	fmt.Println("Conn Stop() ConnID = ", c.ConnID)
}

// GetTCPConnection 获取连接的原始socket
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据给对方客户端
func (c *Connection) Send(data []byte) error {
	return nil
}

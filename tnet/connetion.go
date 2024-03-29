package tnet

import (
	"errors"
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

	// 无缓冲管道，用于读、写goroutine之间的消息通信
	msgChan chan []byte

	// 绑定MsgID和对应的处理业务API关系
	MsgHandle tiface.IMsgHandle

	// 告知当前连接已经退出/停止的channel
	ExitChan chan bool
}

// Connection implements tiface.IConnection
var _ tiface.IConnection = (*Connection)(nil)

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandle tiface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		msgChan:   make(chan []byte),
		MsgHandle: msgHandle,
		ExitChan:  make(chan bool, 1),
	}

	return c
}

// StartReader 启动连接的读数据业务
func (c *Connection) startReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 创建一个拆包解包对象
		dp := NewDataPack()

		// 读取客户端的Msg Head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		// 拆包，得到msgID和msgDataLen放在msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		// 根据dataLen，再次读取Data，存储到msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err: ", err)
				break
			}
		}

		msg.SetData(data)

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalConfig.Gos.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给Worker工作池处理即可
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中，找到注册绑定的Conn对应的router调用
			// 根据绑定好的MsgID找到对应处理api业务 执行
			go c.MsgHandle.DoMsgHandler(&req)
		}
	}
}

// StartWriter 启动连接的写数据业务
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

// 提供一个SendMsg方法，将我们要发送的数据先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}

	// 将data进行封包 MsgDataLen|MsgID|Data
	dp := NewDataPack()

	pack, err := dp.Pack(NewMsgPack(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id:", msgId)
		return errors.New("Pack error msg")
	}

	// 将数据发送到channel中
	c.msgChan <- pack

	return nil
}

// Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)

	// 启动从当前连接的读数据的业务
	go c.startReader()
	// TODO 启动从当前连接写数据的业务
	go c.StartWriter()
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

	fmt.Println("Conn Stop() ConnID = ", c.ConnID)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

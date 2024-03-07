package tnet

import (
	"fmt"
	"net"

	"github.com/cpf2021-gif/gos/tiface"
	"github.com/cpf2021-gif/gos/utils"
)

type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器的端口
	Port int
	// 当前Server添加一个MsgHandle，用来绑定MsgID和对应的处理业务API关系
	MsgHandle tiface.IMsgHandle
}

// Server implements tiface.IServer
var _ tiface.IServer = (*Server)(nil)

func (s *Server) Start() {
	go func() {
		fmt.Printf("[Gos] Server Name: %s, listenner at IP: %s, Port: %d is starting\n",
			utils.GlobalConfig.ServerCfg.Name,
			utils.GlobalConfig.ServerCfg.Host,
			utils.GlobalConfig.ServerCfg.TcpPort)
		fmt.Printf("[Gos] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
			utils.GlobalConfig.Gos.Version,
			utils.GlobalConfig.Gos.MaxConn,
			utils.GlobalConfig.Gos.MaxPacketSize)

		// 0. 启动worker工作池
		s.MsgHandle.StartWorkerPool()

		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("ResolveTCPAddr error: ", err)
			return
		}

		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("ListenTCP error: ", err)
			return
		}

		fmt.Println("Start [gos] Server v0.1 successfully, ", "Name: ", s.Name, " Listening...")
		var cid uint32 = 0

		// 3. 阻塞的等待客户端连接，处理客户端连接业务(读写)
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP error: ", err)
				continue
			}

			// 绑定业务方法和客户端连接
			dealConn := NewConnection(conn, cid, s.MsgHandle)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态或者一些已经开辟的连接信息进行停止或回收
}

func (s *Server) Serve() {
	// 异步启动服务器
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgId uint32, router tiface.IRouter) {
	s.MsgHandle.AddRouter(msgId, router)
	fmt.Println("Add Router successfully!", "msgId =", msgId)
}

// NewServer creates a new server
func NewServer() tiface.IServer {
	s := &Server{
		Name:      utils.GlobalConfig.ServerCfg.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalConfig.ServerCfg.Host,
		Port:      utils.GlobalConfig.ServerCfg.TcpPort,
		MsgHandle: NewMsgHandle(),
	}

	return s
}

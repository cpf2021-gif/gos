package tnet

import (
	"fmt"
	"io"
	"net"

	"github.com/cpf2021-gif/gos/tiface"
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
}

// Server implements tiface.IServer
var _ tiface.IServer = (*Server)(nil)

func (s *Server) Start() {
	go func() {
		fmt.Printf("[Start] Server Listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

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

		// 3. 阻塞的等待客户端连接，处理客户端连接业务(读写)
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP error: ", err)
				continue
			}

			// 已经与客户端建立连接，做一些业务，做一个最基本的回显服务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)

					if err == io.EOF {
						fmt.Println("client Close")
						break
					}

					if err != nil {
						fmt.Println("Receive Buffer error: ", err)
						continue
					}

					fmt.Printf("Receive Buffer: %s, cnt = %d\n", buf, cnt)

					// 回显
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("Write Buffer error: ", err)
						continue
					}
				}
			}()
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

// NewServer creates a new server
func NewServer(name string) tiface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}

	return s
}

package main

import (
	"fmt"

	"github.com/cpf2021-gif/gos/tiface"
	"github.com/cpf2021-gif/gos/tnet"
)

type PingRouter struct {
	tnet.BaseRouter
}

// Test PreHandle
func (pr *PingRouter) PreHandle(request tiface.IRequest) {
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping... ping... ping...\n"))
	if err != nil {
		fmt.Println("call back before ping... ping... ping... error")
	}
}

// Test Handle
func (pr *PingRouter) Handle(request tiface.IRequest) {
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping... ping... ping...\n"))
	if err != nil {
		fmt.Println("call back ping... ping... ping... error")
	}
}

// Test PostHandle
func (pr *PingRouter) PostHandle(request tiface.IRequest) {
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping... ping... ping...\n"))
	if err != nil {
		fmt.Println("call back after ping... ping... ping... error")
	}
}

func main() {
	// 创建一个server句柄
	s := tnet.NewServer("[gos] Server v0.1")

	// 配置路由
	s.AddRouter(&PingRouter{})

	// 开启服务
	s.Serve()
}

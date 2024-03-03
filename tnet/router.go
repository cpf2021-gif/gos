package tnet

import (
	"github.com/cpf2021-gif/gos/tiface"
)

type BaseRouter struct{}

// BaseRouter implements tiface.IRouter
var _ tiface.IRouter = (*BaseRouter)(nil)

func (br *BaseRouter) PreHandle(request tiface.IRequest) {}

func (br *BaseRouter) Handle(request tiface.IRequest) {}

func (br *BaseRouter) PostHandle(request tiface.IRequest) {}

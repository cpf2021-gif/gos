package tiface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(router IRouter)
}

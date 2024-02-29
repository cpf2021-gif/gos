package tiface

type IServer interface {
	Start()
	Stop()
	Serve()
}

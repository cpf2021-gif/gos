package tiface

type IRouter interface {
	// 处理业务之前的Hook
	PreHandle(request IRequest)

	// 处理业务的Hook
	Handle(request IRequest)

	// 处理业务之后的Hook
	PostHandle(request IRequest)
}

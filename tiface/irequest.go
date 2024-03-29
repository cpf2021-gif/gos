package tiface

type IRequest interface {
	// 得到当前连接
	GetConnection() IConnection

	// 得到请求的消息数据
	GetData() []byte

	// 得到请求的消息ID
	GetMsgID() uint32
}

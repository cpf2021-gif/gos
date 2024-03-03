package tnet

import "github.com/cpf2021-gif/gos/tiface"

type Request struct {
	// 已经和客户端建立好的连接
	conn tiface.IConnection

	// 客户端请求的数据
	data []byte
}

// Request implements tiface.IRequest
var _ tiface.IRequest = (*Request)(nil)

// GetConnection 得到当前连接
func (r *Request) GetConnection() tiface.IConnection {
	return r.conn
}

// GetData 得到客户端请求的数据
func (r *Request) GetData() []byte {
	return r.data
}

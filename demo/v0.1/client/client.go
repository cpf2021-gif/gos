package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/cpf2021-gif/gos/tnet"
)

func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)

	// 创建一个客户端
	conn, err := net.Dial("tcp", "0.0.0.0:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 发送封包message消息
		dp := tnet.NewDataPack()
		pack, err := dp.Pack(tnet.NewMsgPack(0, []byte("hello, gos!")))
		if err != nil {
			fmt.Println("Pack error:", err)
			break
		}

		_, err = conn.Write(pack)
		if err != nil {
			fmt.Println("write error:", err)
			break
		}

		// 服务器应该有message返回回来，msgId:1
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head error:", err)
			break
		}

		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("client unpack err:", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			// 读取包体
			msg := msgHead.(*tnet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Recv Server Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}

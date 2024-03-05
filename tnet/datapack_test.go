package tnet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestNewDataPack(t *testing.T) {
	/*
		模拟服务器
	*/
	done := make(chan struct{})
	recvNum := 0

	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err: ", err)
			}

			go func(conn net.Conn) {
				// 处理客户端请求
				// 创建拆包解包对象
				dp := NewDataPack()

				for {
					// 1 第一次从conn读，把包的Head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err: ", err)
						break
					}

					// 拆包，得到msgid和datalen放在msg消息中
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err: ", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						// 2 第二次从conn读，根据datalen再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						// 根据datalen的长度再从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err: ", err)
							return
						}

						// 完整的一个消息已经读取完毕
						fmt.Println("----> Recv MsgID: ", msg.Id, " datalen = ", msg.DataLen, " data = ", string(msg.Data))
						recvNum++
						if recvNum == 2 {
							done <- struct{}{}
						}
					}
				}
			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}

	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg一同发送
	// 封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err: ", err)
		return
	}
	// 封装第二个msg2包
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'w', 'o', 'r', 'l', 'd'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err: ", err)
		return
	}
	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)
	// 一次性发送给服务端
	_, err = conn.Write(sendData1)
	if err != nil {
		t.Fatal(err)
	}

	// 客户端阻塞

	<-done
	fmt.Println("TestNewDataPack test passed!")
}

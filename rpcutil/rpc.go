package rpcutil

import (
	"log"
	"net/rpc"
)

// ClientEnd 客户端的信息
type ClientEnd struct {
	Addr   string      // IP:port
	Client *rpc.Client //  与对端通信的rpc工具
}

// Call 向对端发送请求,使用这个函数来进行RPC通信, 请求运行对端名为methodName的函数, args是函数参数, reply是函数返回值
func (c *ClientEnd) Call(methodName string, args interface{}, reply interface{}) bool {
	if c.Client == nil {
		connect, err := TryConnect(c.Addr)
		if err != nil || connect == nil {
			return false
		}
		c.Client = connect
	}
	// 使用client来进行RPC的call通信, 在client上运行名为methodName的函数
	// args参数是methodName的参数，reply是返回的数据
	err := c.Client.Call(methodName, args, reply)
	if err != nil {
		log.Panicln(err)
		return false
	}
	return true
}

// TryConnect 尝试连接
func TryConnect(addr string) (*rpc.Client, error) {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return client, nil
}

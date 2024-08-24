package db

import (
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"kiritoabc/raft-kv-db/pkg/log"
	"kiritoabc/raft-kv-db/rpcutil"
	"os"
)

// ClientConfig 客户端配置
type ClientConfig struct {
	ClientEnd []struct {
		Ip   string `yaml:"ip"`
		Port string `yaml:"port"`
	} `yaml:"clientEnd"`
}

// KVClient 客户端维护的一些数据
type KVClient struct {
	// 服务器列表
	servers []*rpcutil.ClientEnd
	id      uuid.UUID
	servlen int
	leader  int
}

// generateUUID 生成一个全局唯一的ID
func generateUUID() uuid.UUID {
	id := uuid.NewV1()
	return id
}

// NewKVClient 创建一个KV-DB的客户端
func NewKVClient(servers []*rpcutil.ClientEnd) *KVClient {
	return &KVClient{
		servers: servers,
		id:      generateUUID(),
		servlen: len(servers),
		leader:  0,
	}
}

// Get 客户端Get函数
func (ck *KVClient) Get(key string) string {
	// 构造Get指令远程调用的参数
	args := &GetArgs{
		Key:    key,
		Id:     ck.id,
		Serial: generateUUID(),
	}
	reply := &GetReply{}
	log.Log.Debugf("[%v] 发送 Get 请求 {Key=%v Serial=%v}", ck.id, key, args.Serial)

	for {
		// RPC通信，请求服务器的Get指令
		if ok := ck.servers[ck.leader].Call(RPCGet, args, reply); !ok {
			log.Log.Debugf("[%v] 对 服务器 %v 的 Get 请求 (Key=%v Serial=%v) 超时",
				ck.id, ck.leader, key, args.Serial)

			// 切换领导者
			ck.leader = (ck.leader + 1) % ck.servlen
			continue
		}
		if reply.Err == OK {
			log.Log.Debugf("[%v] 收到对 %v 发送的 Get 请求 {Key=%v Serial=%v} 的响应，结果为 %v",
				ck.id, ck.leader, key, args.Serial, reply.Value)
			return reply.Value
		} else if reply.Err == ErrNoKey {
			// 没有对应的Key
			log.Log.Debugf("[%v] 收到对 %v 发送的 Get 请求 {Key=%v Serial=%v} 的响应，结果为 ErrNoKey",
				ck.id, ck.leader, key, args.Serial)
			return NoKeyValue
		} else if reply.Err == ErrWrongLeader {
			// 当前请求的服务器不是leader
			log.Log.Debug("错误的领导者")
			// 请求了错误的领导者，切换请求新的服务器
			ck.leader = (ck.leader + 1) % ck.servlen
			// ck.leader = 0
			continue
		} else {
			// 其他错误
			log.Log.Fatalf("[%v] 对 服务器 %v 的 Get 请求 (Key=%v Serial=%v) 收到一条空 Err",
				ck.id, ck.leader, key, args.Serial)
		}
	}
}

// Put 客户端Put接口
func (ck *KVClient) Put(key string, value string) {
	ck.putAppend(key, value, OpPut)
}

// Append 客户端Append接口
func (ck *KVClient) Append(key string, value string) {
	ck.putAppend(key, value, OpAppend)
}

// putAppend 客户端Put或者Append操作
func (ck *KVClient) putAppend(key string, value string, op string) {
	// 构造Put或者Append操作的参数
	args := &PutAppendArgs{
		Key:    key,
		Value:  value,
		Op:     op,
		Id:     ck.id,
		Serial: generateUUID(),
	}
	reply := &PutAppendReply{}
	log.Log.Debugf("[%v] 发送 PA 请求 {Op=%v Key=%v Value='%v' Serial=%v}", ck.id, op, key, value, args.Serial)

	for {
		// 与服务器进行RPC通信, 调用服务器的函数
		if ok := ck.servers[ck.leader].Call(RPCPutAppend, args, reply); !ok {
			log.Log.Debugf("[%v] 对 服务器 %v 的 PutAppend 请求 (Serial=%v Key=%v Value=%v op=%v) 超时",
				ck.id, ck.leader, args.Serial, key, value, op)
			ck.leader = (ck.leader + 1) % ck.servlen
			continue
		}
		if reply.Err == OK {
			log.Log.Debugf("[%v] 收到对 %v 发送的 PA 请求 {Op=%v Key=%v Value='%v' Serial=%v} 的响应，结果为 OK",
				ck.id, ck.leader, op, key, value, args.Serial)
			return
		} else if reply.Err == ErrWrongLeader {
			// 请求了错误的leader，更换请求leader
			ck.leader = (ck.leader + 1) % ck.servlen
			continue
		} else {
			log.Log.Fatalf("[%v] 对 服务器 %v 的 PutAppend 请求 (Serial=%v Key=%v Value=%v op=%v) 收到一条空 Err",
				ck.id, ck.leader, args.Serial, key, value, op)
		}
	}
}

// GetClientEnds 获取客户端通信实例
func GetClientEnds(path string) []*rpcutil.ClientEnd {
	config := getClientConfig(path)
	num := len(config.ClientEnd)
	if (num&1) == 0 || num < 3 {
		log.Log.Fatal("the number of servers must be odd and greater than or equal to 3")
	}

	clientEnds := make([]*rpcutil.ClientEnd, 0)
	for _, end := range config.ClientEnd {
		address := end.Ip + ":" + end.Port
		client, _ := rpcutil.TryConnect(address)

		clientEnds = append(clientEnds, &rpcutil.ClientEnd{
			Addr:   address,
			Client: client,
		})
	}
	return clientEnds
}

// getClientConfig 获取客户端配置
func getClientConfig(path string) *ClientConfig {
	if len(os.Args) != 1 {
		path = os.Args[1]
	}
	log.Log.Infof("installing config %s", path)
	cfg := ClientConfig{}
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Log.Fatalf("read config failed: %v", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Log.Fatalf("unmarshal config failed: %v", err)
	}
	return &cfg
}

package main

import (
	"github.com/spf13/viper"
	"kiritoabc/raft-kv-db/db"
	"kiritoabc/raft-kv-db/pkg/log"
	"kiritoabc/raft-kv-db/raft"
	"kiritoabc/raft-kv-db/rpcutil"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

// Config 配置文件结构体Config
type Config struct {
	// 服务器的Ip和对应的Port
	Servers []struct {
		Ip   string
		Port string
	}
	Me int `yaml:"me"`
}

const path = "config/server.yaml"

func main() {
	// 首先解析cfg文件
	serverCfg := getCfg(path)
	i := len(serverCfg.Servers)
	if (i & 1) == 0 {
		panic("总服务器数量必须为单数")
	}

	// 获取集群的所有服务器列表
	clientEnds := getClientEnds(serverCfg)
	// 持久化工具
	persister := raft.NewPersister()
	// 获取一个kvserver实例
	kvur := db.StartKVServer(clientEnds, serverCfg.Me, persister, i)

	// 首先将kvur进行注册，注册之后的方法才可以进行RPC调用
	// 使用rpc.Register进行注册
	if err := rpc.Register(kvur); err != nil {
		panic(err)
	}

	// 注册 HTTP 路由
	rpc.HandleHTTP()
	// 监听端口
	l, e := net.Listen("tcp", ":"+serverCfg.Servers[serverCfg.Me].Port)
	if e != nil {
		panic(e)
	}
	log.Log.Infof("Listen", serverCfg.Servers[serverCfg.Me].Port)

	// 启动一个Http服务，请求rpc的方法会交给rpc内部路由进行处理
	log.Log.Fatal(http.Serve(l, nil))
}

func getClientEnds(serverCfg Config) []*rpcutil.ClientEnd {
	// 创建clientEnds数组，初始长度为0
	clientEnds := make([]*rpcutil.ClientEnd, 0)
	// 遍历所有server
	for i, end := range serverCfg.Servers {
		// server的地址: IP:PORT
		address := end.Ip + ":" + end.Port
		var client *rpc.Client
		if i == serverCfg.Me {
			client = nil
		} else {
			// 如果不是自身，则尝试与别的服务器进行进行连接
			client, _ = rpcutil.TryConnect(address)
		}

		// 新建一个ClientEnd服务器
		tmpClientEnd := &rpcutil.ClientEnd{
			Addr:   address,
			Client: client,
		}

		clientEnds = append(clientEnds, tmpClientEnd)
	}
	return clientEnds
}

// getCfg 获取配置文件
func getCfg(path string) Config {
	cfg := Config{}
	// 如果没有指定config文件的地址，则使用默认的config文件
	if len(os.Args) != 1 {
		path = os.Args[1]
	}

	// 读取文件
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Log.Fatalf("Unable to decode into struct, %v", err)
	}

	return cfg
}

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

// config 配置文件路径
const config = "config/server.yaml"

// Config 配置文件结构体Config
type Config struct {
	// 服务器的Ip和对应的Port
	Servers []struct {
		Ip   string
		Port string
	}
	Me int `yaml:"me"`
}

func main() {
	serverCfg := getConf(config)
	num := len(serverCfg.Servers)
	// 服务器必须是单数
	if (num & 1) == 0 {
		log.Log.Fatalf("server number must be odd")
	}
	// 获取服务器列表
	clientEnds := getClientEnds(serverCfg)
	// 持久化工具
	persister := raft.NewPersister()
	// 获取一个kvserver实例
	kvur := db.StartKVServer(clientEnds, serverCfg.Me, persister, num)

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
	log.Log.Println("Listen", serverCfg.Servers[serverCfg.Me].Port)

	// 启动一个Http服务，请求rpc的方法会交给rpc内部路由进行处理
	http.Serve(l, nil)
}

// getClientEnds 获取客户端列表
func getClientEnds(serverCfg Config) []*rpcutil.ClientEnd {
	clientEnds := make([]*rpcutil.ClientEnd, 0)
	for i, server := range serverCfg.Servers {
		// 服务器的地址: IP:PORT
		address := server.Ip + ":" + server.Port
		var client *rpc.Client
		// 如果当前服务器是自己，则client为nil
		if i == serverCfg.Me {
			client = nil
		} else {
			client, _ = rpcutil.TryConnect(address)
		}

		clientEnds = append(clientEnds, &rpcutil.ClientEnd{
			Addr:   address,
			Client: client,
		})
	}

	return clientEnds
}

// getConf 获取配置文件
func getConf(path string) Config {
	if len(os.Args) != 1 {
		path = os.Args[1]
	}
	log.Log.Infof("installing config %s", path)
	cfg := Config{}
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Log.Fatalf("read config failed: %v", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Log.Fatalf("unmarshal config failed: %v", err)
	}
	return cfg
}

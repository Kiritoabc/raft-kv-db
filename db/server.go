package db

import (
	"github.com/google/uuid"
	"kiritoabc/raft-kv-db/pkg/log"
	"kiritoabc/raft-kv-db/raft"
	"kiritoabc/raft-kv-db/rpcutil"
	"sync"
)

// Op command定义 是客户端发送给Raft的命令
type Op struct {
	Type   string
	Key    string
	Value  string
	Serial uuid.UUID
}

// CommonReply 是Raft返回给客户端的命令执行结果
type CommonReply struct {
	Err    Err
	Key    string
	Value  string
	Serial *uuid.UUID
}

// Get 接口, 供Client进行RPC调用
func (kv *KVServer) Get(args *GetArgs, reply *GetReply) error {
	op := &Op{
		Type:   OpGet,
		Key:    args.Key,
		Value:  NoKeyValue,
		Serial: args.Serial,
	}
	reply.Err = ErrWrongLeader
	// 调用Raft节点的start接口执行具体的操作
	idx, _, isLeader := kv.rf.Start(*op)
	// 如果不是leader, 就返回
	if !isLeader {
		log.Log.Infof("%v 对于 %v 的 Get 请求 {Key=%v Serial=%v} 处理结果为 该服务器不是领导者 \n",
			kv.me, args.Id, args.Key, args.Serial)
		return nil
	}
	log.Log.Infof("%v 对于 %v 的 Get 请求 {Key=%v Serial=%v} 处理结果为 等待提交 \n",
		kv.me, args.Id, args.Key, args.Serial, idx)

	// Raft节点日志执行完毕, Apply到数据库进行具体的操作
	commonReply := &CommonReply{}
	// 在commonReplies中查找对应的reply
	find := kv.findReply(op, idx, commonReply)
	// 如果找到就设置对应的reply.Value
	if find == OK {
		reply.Value = commonReply.Value
		reply.Err = commonReply.Err
	}
	log.Log.Infof("%v 对于 %v 的 Get 请求 {Key=%v Serial=%v} 处理结果为 查找reply成功 \n",
		kv.me, args.Id, args.Key, args.Serial, reply.Err, len(reply.Value))

	return nil
}

// KVServer 实现了KVServer接口
type KVServer struct {
	mu sync.Mutex
	me int // 当前节点的索引号
	//rf      *raft.Raft         // 当前节点对应的Raft实例
	//applyCh chan raft.ApplyMsg // Apply通道，用来执行命令
	dead int32

	maxraftstate int // 如果日志变得很长，就备份一份快照

	data          map[string]string // 存储的数据
	commonReplies []*CommonReply    // 命令执行的返回值
}

// StartKVServer starts the server for a KV group.
func StartKVServer(servers []*rpcutil.ClientEnd, me int, persister *raft.Persister, maxraftstate int) *KVServer {

}

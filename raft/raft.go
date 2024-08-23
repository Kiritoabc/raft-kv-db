package raft

import (
	"kiritoabc/raft-kv-db/rpcutil"
	"sync"
)

const (
	FOLLOWER    = 0
	CANDIDATE   = 1
	LEADER      = 2
	CopyEntries = 3
	HeartBeat   = 4
)

// ApplyMsg 用于将日志命令Apply到状态机的数据结构
type ApplyMsg struct {
	CommandValid bool        // 命令是否有效
	Command      interface{} // 具体的命令
	CommandIndex int         // 命令索引
}

// LogEntry 日志条目
type LogEntry struct {
	Command interface{} // 具体的命令
	Term    int         // 日志所在任期
	Index   int         // 日志索引
}

// CommandState 命令的状态
type CommandState struct {
	Term  int // 命令所在任期
	Index int // 命令索引
}

// Raft 实现Raft一致性算法的结构体
type Raft struct {
	mu        sync.Mutex           // 互斥锁，用于共享变量的同步
	peers     []*rpcutil.ClientEnd // 服务器列表, 用来进行RPC通信的(ClientEnd实现参见package rpcutil)
	persister *Persister           // 持久化对象
	me        int                  // 当前服务器在服务器列表peers中的索引
	dead      int32                // 当前服务器是否宕机

	// 一个Raft节点必须要维护的一些数据结构
	currentTerm int        // 服务器最后一次知道的任期号（初始化为 0，持续递增）
	votedFor    int        // 在当前获得选票的候选人的 Id
	logs        []LogEntry // 日志条目集；每一个条目包含一个用户状态机执行的指令，和收到时的任期号

	// lastApplied的ID不能大于commitIndex的ID，因为只有commited状态的日志才能被提交到状态机执行
	commitIndex int // 已知的最大的已经被提交的日志条目的索引值
	lastApplied int // 最后被应用到状态机的日志条目索引值（初始化为 0，持续递增）

	nextIndex  []int // 对于每一个服务器，需要发送给他的下一个日志条目的索引值（初始化为领导人最后索引值加一）
	matchIndex []int // 对于每一个服务器，已经复制给他的日志的最高索引值（初始化为0）

	identity     int // 当前服务器的身份(follower, candidate, or leader)
	peersLen     int // 服务器的总数
	hbCnt        int
	applyCh      chan ApplyMsg                 // 用于提交命令的管道
	doAppendCh   chan int                      // 用于AppendEntries的管道(用于heartbeat或者是appendEntries)
	applyCmdLogs map[interface{}]*CommandState // 用于保存命令的字典
}

// GetState 返回当前服务器状态（当前任期号以及当前服务器是否是leader）
func (rf *Raft) GetState() (int, bool) {
	panic("implate me ")
}

// persist 将raft的状态进行持久化
// 需要持久化的数据有votedFor, currentTerm, 以及对应的日志logs
func (rf *Raft) persist() {
	panic("implate me ")
}

// readPersist 从之前保存的持久化状态恢复
func (rf *Raft) readPersist(data []byte) {
	panic("implate me ")
}

// AppendEntriesArgs AppendEntries函数所需要的参数
type AppendEntriesArgs struct {
	Term         int        // 领导人的任期号
	LeaderId     int        // 领导人的 Id，以便于跟随者重定向请求
	PrevLogIndex int        // 新的日志条目紧随之前的索引值
	PrevLogTerm  int        // PrevLogIndex 条目的任期号
	Entries      []LogEntry // 准备存储的日志条目（表示心跳时为空；一次性发送多个是为了提高效率）
	LeaderCommit int        // 领导人已经提交的日志的索引值, 用于更新follower的commitIndex, 然后进行日志的提交
}

// AppendEntriesReply AppendEntries函数的返回值
type AppendEntriesReply struct {
	Term         int  // 当前的任期号，用于领导人去更新自己
	Success      bool // 跟随者包含了匹配上 PrevLogIndex 和 PrevLogTerm 的日志时为真
	PrevLogIndex int  // 返回跟随者的最长的PrevLogIndex, 当不匹配时, 也就是follower的log长度
}

// RequestVoteArgs RequestVote函数的参数
type RequestVoteArgs struct {
	Term         int // 候选人的任期号
	CandidateId  int // 请求选票的候选人的 Id
	LastLogIndex int // 候选人的最后日志条目的索引值
	LastLogTerm  int // 候选人最后日志条目的任期号
	// 上面的任期号，日志索引以及日志任期号是用来判别服务器日志的更新状态
}

// RequestVoteReply RequestVote函数的返回值
type RequestVoteReply struct {
	Term        int  // 当前任期号，以便于候选人去更新自己的任期号
	VoteGranted bool // 候选人赢得了此张选票时为真
}

// RequestVote 函数，用于candidate请求其他服务器进行投票, 如果超过半数为其投票, 则当前candidate变为leader
// 投票选取leader的条件是：leader的Term是最大的，同时日志的Index也是最大的
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	panic("implate me ")
}

// AppendEntries 函数，有两个功能，HeartBeat和log replication功能 ==>  1. 心跳检测 2. 日志复制
func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	panic("implate me ")
}

// sendRequestVote RequestVote的RPC接口, 请求ID为server的投票
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	panic("implate me ")
}

// sendAppendEntries AppendEntries的RPC接口, 请求ID为server的日志复制和HeartBeat
func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	panic("implate me ")
}

// Start Raft节点执行命令, 调用此接口来是Raft节点执行对应的命令, command为具体的命令
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	panic("implate me ")
}

// Kill 设置Raft节点为dead状态
func (rf *Raft) Kill() {
	panic("implate me ")
}

// killed 判断当前Raft节点是否宕机
func (rf *Raft) killed() bool {
	panic("implate me ")
}

// sendLogEntry 向所有follower发送日志
func (rf *Raft) sendLogEntry(flag int) {
	panic("implate me ")
}

// setToFollower 设置当前Raft节点为FOLLOWER
func (rf *Raft) setToFollower() {
	panic("implate me ")
}

// apply 将commited状态的命令提交到状态机执行		（提交commit的命令）
func (rf *Raft) apply() {
	panic("implate me ")
}

// randomTimeout 生成一个随机超时时间。 （随机超时时间）
func randomTimeout(min, max int) int {
	panic("implate me ")
}

// NewRaft 返回一个Raft节点实例
/*
 参数：
	peers: 		表示集群中的所有服务器
	me:			表示当前服务器在peers中的索引
	persister:	一个持久化对象
	applyCh:	Apply通道, 将需要状态机Apply的命令传入到applyCh通道中, 数据库从applyCh通道中取命令行
 返回:
	一个Raft实例
*/
func NewRaft(peers []*rpcutil.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	panic("implate me ")
}

// leaderElection leader选举函数
func (rf *Raft) leaderElection(wonCh chan int, wgp *sync.WaitGroup) {
	panic("implate me ")
}

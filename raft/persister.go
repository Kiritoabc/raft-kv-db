package raft

import "sync"

// Persister 持久化对象
type Persister struct {
	mu        sync.Mutex // 互斥锁
	raftstate []byte     // raft状态变量
	snapshot  []byte     //raft快照变量
}

// NewPersister 创建持久化对象
func NewPersister() *Persister {
	return &Persister{}
}

// Copy Persister对象的拷贝
func (ps *Persister) Copy() *Persister {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	np := NewPersister()
	np.raftstate = ps.raftstate
	np.snapshot = ps.snapshot
	return np
}

// SaveRaftState 保存Raft状态
func (ps *Persister) SaveRaftState(state []byte) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.raftstate = state
}

// ReadRaftState 读取Raft状态
func (ps *Persister) ReadRaftState() []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.raftstate
}

// RaftStateSize 读取Raft的大小
func (ps *Persister) RaftStateSize() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return len(ps.raftstate)
}

// SaveStateAndSnapshot Save both Raft state and K/V snapshot as a single atomic action,
// to help avoid them getting out of sync.
func (ps *Persister) SaveStateAndSnapshot(state []byte, snapshot []byte) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.raftstate = state
	ps.snapshot = snapshot
}

// ReadSnapshot Read snapshot of the key-value server
func (ps *Persister) ReadSnapshot() []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.snapshot
}

// SnapshotSize Read snapshot of the key-value server
func (ps *Persister) SnapshotSize() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return len(ps.snapshot)
}

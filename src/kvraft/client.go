package kvraft

import (
	"crypto/rand"
	"math/big"
	"sync/atomic"

	"6.824/labrpc"
)

type Clerk struct {
	servers []*labrpc.ClientEnd // client terminals
	// You will have to modify this struct.
	leaderId int // raft leader peer id, to avoid researching leader every time
	// cliendId and sequenceNumber are used to identify the client's operation
	clientId       int64
	sequenceNumber int64
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(servers []*labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.servers = servers
	// You'll have to add code here.
	ck.leaderId = None
	ck.clientId = nrand()
	ck.sequenceNumber = 0
	return ck
}

// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
// Get操作是幂等的
func (ck *Clerk) Get(key string) string {

	// You will have to modify this function.
	DPrintf(false, "[clerk-start] Client try to Get Key = %s", key)
	atomic.AddInt64(&ck.sequenceNumber, 1) // mark client's operation
	// If leaderId is not -1, try to send Get RPC to the leader
	if ck.leaderId != None {
		args := GetArgs{
			Key:            key,
			ClientId:       ck.clientId,
			SequenceNumber: ck.sequenceNumber,
		}
		reply := GetReply{}
		ok := ck.servers[ck.leaderId].Call("KVServer.Get", &args, &reply)
		// current operation is successful
		if ok && reply.Err == "" {
			DPrintf(false, "[clerk-over] Client Get Key = %s, val = %s, old learderId = %d", key, reply.Value, ck.leaderId)
			return reply.Value
		}
	}
	// If leaderId is -1 or operation is failed, try to keep sending Get RPC to all servers until success
	for {
		for i := range ck.servers {
			args := GetArgs{
				Key:            key,
				ClientId:       ck.clientId,
				SequenceNumber: ck.sequenceNumber,
			}
			reply := GetReply{}
			ok := ck.servers[i].Call("KVServer.Get", &args, &reply)
			if ok && reply.Err == "" {
				ck.leaderId = i // find the leader
				DPrintf(false, "[clerk-over] Client Get Key = %s, val = %s, new leaderId = %d", key, reply.Value, ck.leaderId)
				return reply.Value
			}
		}
	}
}

// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.PutAppend", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) PutAppend(key string, value string, op string) {
	// You will have to modify this function.
	// DPrintf()
	atomic.AddInt64(&ck.sequenceNumber, 1) // mark client's operation
	// If leaderId is not -1, try to send PutAppend RPC to the leader
	if ck.leaderId != None {

	}
	for {

	}
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
}

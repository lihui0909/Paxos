package entity

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type Learner struct {
	lis net.Listener
	// 学习者id
	id int
	// 记录接受者已接受的提案：[接受者 id]请求消息
	acceptedMsg map[int]MsgArgs
}

func (l *Learner) Learn(args *MsgArgs, reply *MsgReply) error {
	a := l.acceptedMsg[args.From] //（a为接受者对应的历史提案）
	//(如果之前保存的这个接受者的提案编号小于这次的提案编号，更新历史记录）
	if a.Number < args.Number {
		l.acceptedMsg[args.From] = *args
		reply.Ok = true
	} else {
		reply.Ok = false
	}
	return nil
}

func (l *Learner) chosen() interface{} {
	//(记录每一个提案编号有多少个接受者）
	acceptCounts := make(map[int]int)
	// （记录每一个提案编号对应的提案）
	acceptMsg := make(map[int]MsgArgs)

	for _, accepted := range l.acceptedMsg {
		if accepted.Number != 0 {
			acceptCounts[accepted.Number]++
			acceptMsg[accepted.Number] = accepted
		}
	}

	for n, count := range acceptCounts {
		if count >= l.majority() {
			return acceptMsg[n].Value
		}
	}
	return nil
}

func (l *Learner) majority() int {
	return len(l.acceptedMsg)/2 + 1
}

func newLearner(id int, acceptorIds []int) *Learner {
	learner := &Learner{
		id:          id,
		acceptedMsg: make(map[int]MsgArgs),
	}
	for _, aid := range acceptorIds {
		learner.acceptedMsg[aid] = MsgArgs{
			Number: 0,
			Value:  nil,
		}
	}
	learner.server(id)
	return learner
}

func (l *Learner) server(id int) {
	rpcs := rpc.NewServer()
	rpcs.Register(l)
	addr := fmt.Sprintf(":%d", id)
	lis, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error : ", e)
	}
	l.lis = lis
	go func() {
		for {
			conn, err := l.lis.Accept()
			if err != nil {
				continue
			}
			go rpcs.ServeConn(conn)
		}
	}()
}

//关闭连接
func (l *Learner) close() {
	l.lis.Close()
}

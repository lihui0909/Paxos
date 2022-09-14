package entity

import "fmt"

type Proposer struct {
	//服务器 id
	id int
	//当前提议者已知的最大轮次（单调自增）
	round int
	//提案编号 = (轮次，服务器id)
	number int
	//接受者id列表（接受者RPC服务的端口）
	acceptors []int
}

func (p *Proposer) propose(v interface{}) interface{} {
	p.round++
	p.number = p.proposalNumber()

	//第一阶段(phase 1)(Perpare阶段）
	prepareCount := 0 //(记录返回promise的acceptor个数）
	maxNumber := 0    //(记录返回Promise中最大的提案编号）
	for _, aid := range p.acceptors {
		args := MsgArgs{
			Number: p.number,
			From:   p.id,
			To:     aid,
		}
		reply := new(MsgReply)
		Ok := call(fmt.Sprintf("127.0.0.1:%d", aid), "Acceptor.Prepare", args, reply)
		if !Ok {
			continue
		}

		if reply.Ok {
			prepareCount++
			if reply.Number > maxNumber {
				maxNumber = reply.Number
				v = reply.Value
			}
		}

		if prepareCount == p.majority() {
			break
		}
	}

	//第二阶段(phase 2)(Accept阶段）
	acceptCount := 0
	if prepareCount >= p.majority() {
		for _, aid := range p.acceptors {
			args := MsgArgs{
				Number: p.number,
				Value:  v,
				From:   p.id,
				To:     aid,
			}
			reply := new(MsgReply)
			ok := call(fmt.Sprintf("127.0.0.1:%d", aid), "Acceptor.Accept", args, reply)
			if !ok {
				continue
			}
			if reply.Ok {
				fmt.Printf("acceptCount:%d\n", acceptCount)
				acceptCount++
			}
		}
	}
	fmt.Printf("acceptCount:%d, v:%+v\n", acceptCount, v)

	if acceptCount >= p.majority() {
		//选择的提案值
		return v
	}
	return nil
}

func (p *Proposer) majority() int {
	return len(p.acceptors)/2 + 1
}

func (p *Proposer) proposalNumber() int {
	return p.round<<16 | p.id
}

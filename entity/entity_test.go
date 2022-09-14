package entity

import (
	"fmt"
	"testing"
)

func start(acceptorIds []int, learnerIds []int) ([]*Acceptor, []*Learner) {
	acceptors := make([]*Acceptor, 0)
	for _, aid := range acceptorIds {
		a := newAcceptor(aid, learnerIds)
		acceptors = append(acceptors, a)
	}

	learners := make([]*Learner, 0)
	for _, lid := range learnerIds {
		l := newLearner(lid, acceptorIds)
		learners = append(learners, l)
	}

	return acceptors, learners
}

func cleanup(acceptors []*Acceptor, learners []*Learner) {
	for _, a := range acceptors {
		a.close()
	}
	for _, l := range learners {
		l.close()
	}
}

func TestSingleProposer(t *testing.T) {
	// 1001、1002、1003是接受者id
	acceptorIds := []int{1001, 1002, 1003}
	// 2001 是学习者id
	learnerIds := []int{2001}
	acceptors, learners := start(acceptorIds, learnerIds)

	defer cleanup(acceptors, learners)

	//1是提议者id
	p := &Proposer{
		id:        1,
		acceptors: acceptorIds,
	}

	value := p.propose("hello xiaohui")
	fmt.Println("value:")
	fmt.Println(value)
	if value != "hello xiaohui" {
		t.Errorf("value = %s, excepted %s", value, "hello xiaohui")
	}

	learnValue := learners[0].chosen()
	if learnValue != value {
		t.Errorf("learnValue = %s, excepted %s", learnValue, "hello xiaohui")
	}
}

func TestTwoProposers(t *testing.T) {
	// 1001、1002、1003是接受者id
	acceptorIds := []int{1001, 1002, 1003}
	// 2001 是学习者id
	learnerIds := []int{2001}
	acceptors, learners := start(acceptorIds, learnerIds)

	defer cleanup(acceptors, learners)

	//1、2是提议者id
	p1 := &Proposer{
		id:        1,
		acceptors: acceptorIds,
	}
	value1 := p1.propose("hello world")

	p2 := &Proposer{
		id:        2,
		acceptors: acceptorIds,
	}
	value2 := p2.propose("hello google")

	if value1 != value2 {
		t.Errorf("value1 = %s, value2 = %s", value1, value2)
	}

	learnValue := learners[0].chosen()
	if learnValue != value1 {
		t.Errorf("learnValue = %s, excepted %s", learnValue, value1)
	}
}

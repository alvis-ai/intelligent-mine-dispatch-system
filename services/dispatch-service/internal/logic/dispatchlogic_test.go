package logic

import (
	"context"
	"testing"
)

func TestGetAssigner_FIFO(t *testing.T) {
	l := &DispatchLogic{}
	a := l.getAssigner("fifo")
	if _, ok := a.(*FIFOAssigner); !ok {
		t.Errorf("getAssigner('fifo') = %T, want *FIFOAssigner", a)
	}
}

func TestGetAssigner_NearestFirst(t *testing.T) {
	l := &DispatchLogic{}
	a := l.getAssigner("nearest_first")
	if _, ok := a.(*NearestFirstAssigner); !ok {
		t.Errorf("getAssigner('nearest_first') = %T, want *NearestFirstAssigner", a)
	}
}

func TestGetAssigner_WeightedRoundRobin(t *testing.T) {
	l := &DispatchLogic{}
	a := l.getAssigner("weighted_round_robin")
	if _, ok := a.(*WeightedRoundRobinAssigner); !ok {
		t.Errorf("getAssigner('weighted_round_robin') = %T, want *WeightedRoundRobinAssigner", a)
	}
}

func TestGetAssigner_Default(t *testing.T) {
	l := &DispatchLogic{}
	a := l.getAssigner("unknown")
	if _, ok := a.(*FIFOAssigner); !ok {
		t.Errorf("getAssigner('unknown') = %T, want *FIFOAssigner (default)", a)
	}
}

func TestGetAssigner_Empty(t *testing.T) {
	l := &DispatchLogic{}
	a := l.getAssigner("")
	if _, ok := a.(*FIFOAssigner); !ok {
		t.Errorf("getAssigner('') = %T, want *FIFOAssigner (default)", a)
	}
}

func TestGetAssigner_GeneticAlgorithm(t *testing.T) {
	l := &DispatchLogic{ctx: context.Background()}
	a := l.getAssigner("genetic_algorithm")
	ga, ok := a.(*GeneticAlgorithmAssigner)
	if !ok {
		t.Errorf("getAssigner('genetic_algorithm') = %T, want *GeneticAlgorithmAssigner", a)
	}
	_ = ga
}

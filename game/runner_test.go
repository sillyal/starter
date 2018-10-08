package game

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

type SequentialState struct {
	id int
}

func (state *SequentialState) Next() State {
	if state.id == 1 {
		return &SkippingState{id: state.id}
	} else if state.id > 3 {
		return nil
	}
	return &SequentialState{state.id + 1}
}

type SkippingState struct {
	id    int
	skipN *SkipN
}

func (state *SkippingState) Next() State {
	if state.skipN == nil {
		log.Panic("action is required")
	}
	return &SequentialState{state.id + state.skipN.n}
}

func (state *SkippingState) ActionsAvailable() []Action {
	return []Action{&SkipN{2}, &SkipN{3}}
}

func (state *SkippingState) TakeAction(action Action, reason Reason) State {
	state.skipN = action.(*SkipN)
	return state
}

func (state *SkippingState) Clone() ActionPending {
	return &SkippingState{state.id, state.skipN}
}

type SkipN struct {
	n int
}

type TakeFirst struct {
}

func (takeFirst *TakeFirst) ChooseAction(actionPending ActionPending) (Action, Reason) {
	actions := actionPending.ActionsAvailable()
	if len(actions) == 0 {
		return nil, takeFirst
	} else {
		return actions[0], takeFirst
	}
}

func TestRunner(t *testing.T) {
	runner := &Runner{Strategy: &TakeFirst{}}
	runner.Run(&SequentialState{0})
	assert.Equal(t, 5, len(runner.States))
	assertSequentialState(t, runner.States[0], 0)
	assertSequentialState(t, runner.States[1], 1)
	assertSkippingState(t, runner.States[2], 2)
	assertSequentialState(t, runner.States[3], 3)
	assertSequentialState(t, runner.States[4], 4)
}

func assertSequentialState(t *testing.T, state State, expectedId int) {
	sequentialState, ok := state.(*SequentialState)
	assert.True(t, ok)
	assert.Equal(t, expectedId, sequentialState.id)
}

func assertSkippingState(t *testing.T, state State, expectedSkipN int) {
	skippingState, ok := state.(*SkippingState)
	assert.True(t, ok)
	assert.Equal(t, expectedSkipN, skippingState.skipN.n)
}

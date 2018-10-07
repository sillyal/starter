package starter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type SimpleFeedback struct {
	count int
}

func (feedback *SimpleFeedback) Evaluate(state State, reference State) int {
	state2 := reference.(*TestState)
	if state1 := state.(*TestState); state1.winning == state2.playing {
		return 1
	} else {
		return -1
	}
}

type ActionFavors struct {
	player int
}

type TestState struct {
	playing int
	winning int
	actions []Action
}

func (state *TestState) ActionsAvailable() []Action {
	return state.actions
}

func (state *TestState) TakeAction(action Action) State {
	return &TestState{state.playing, action.(*ActionFavors).player, state.actions}
}

func (state *TestState) Next() State {
	return &TestState{playing: 2, winning: state.winning, actions: []Action{&ActionFavors{1}, &ActionFavors{1}}}
}

func (state *TestState) Clone() ActionPending {
	return &TestState{state.playing, state.winning, state.actions}
}

func TestMonteCarloTreeSearch_ChooseAction(t *testing.T) {
	feedback := &SimpleFeedback{}
	mcts := &MonteCarloTreeSearch{feedback, 2}
	state := &TestState{playing: 1, actions: []Action{&ActionFavors{1}, &ActionFavors{2}}}
	action := mcts.ChooseAction(state)
	assert.Equal(t, 1, action.(*ActionFavors).player)
}

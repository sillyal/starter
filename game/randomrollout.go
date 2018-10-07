package game

import (
	"math/rand"
	"time"
)

type RandomRollout struct {
}

func (rr *RandomRollout) ChooseAction(state ActionPending) Action {
	if actionsAvailable := state.ActionsAvailable(); len(actionsAvailable) > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		perm := r.Perm(len(actionsAvailable))
		return actionsAvailable[perm[0]]
	}
	return nil
}

package game

type State interface {
	Next() State
}

type Action interface {
}

type Reason interface {
}

type ActionPending interface {
	ActionsAvailable() []Action
	TakeAction(action Action, reason Reason) State
	Clone() ActionPending
}

type Strategy interface {
	ChooseAction(actionPending ActionPending) (Action, Reason)
}

type Runner struct {
	Strategy Strategy
	States   []State
}

func (runner *Runner) Run(state State) {
	if runner.Strategy == nil {
		panic("Strategy is required for a runner")
	}
	for state != nil {
		if actionPending, ok := state.(ActionPending); ok {
			action, reason := runner.Strategy.ChooseAction(actionPending)

			state = actionPending.TakeAction(action, reason)
		}
		runner.States = append(runner.States, state)
		state = state.Next()
	}
}

func (runner *Runner) End() State {
	if len(runner.States) == 0 {
		return nil
	}
	return runner.States[len(runner.States)-1]
}

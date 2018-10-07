package starter

type State interface {
	Next() State
}

type Action interface {
}

type ActionPending interface {
	ActionsAvailable() []Action
	TakeAction(action Action) State
	Clone() ActionPending
}

type Strategy interface {
	ChooseAction(actionPending ActionPending) Action
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
		runner.States = append(runner.States, state)

		if actionPending, ok := state.(ActionPending); ok {
			action := runner.Strategy.ChooseAction(actionPending)

			state = actionPending.TakeAction(action)
		}
		state = state.Next()
	}
}

func (runner *Runner) End() State {
	if len(runner.States) == 0 {
		return nil
	}
	return runner.States[len(runner.States)-1]
}

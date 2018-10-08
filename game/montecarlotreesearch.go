package game

import (
	"math"
	"sort"
)

type Feedback interface {
	Evaluate(state State, reference State) int
}

type MonteCarloTreeSearch struct {
	Feedback Feedback
	Round    int
}

type TrivialChoice struct {
	Count int `json:"choices"`
}

func (mcts *MonteCarloTreeSearch) ChooseAction(root ActionPending) (Action, Reason) {
	if actionsAvailable := root.ActionsAvailable(); len(actionsAvailable) < 2 {
		if len(actionsAvailable) == 1 {
			return actionsAvailable[0], &TrivialChoice{1}
		} else {
			return nil, &TrivialChoice{0}
		}
	}
	var nodes []*Node
	origin := &Node{state: root, children: []*Node{}}
	current := origin
	for count := 0; current != nil && count < mcts.Round; count += 1 {
		actions := current.state.ActionsAvailable()
		for _, action := range actions {
			child := &Node{parent: current, Action: action}
			current.children = append(current.children, child)

			nextState := current.state.Clone().TakeAction(action, mcts).Next()
			value := mcts.Feedback.Evaluate(nextState, root.(State))
			child.backfill(value)

			child.state = nextActionPendingState(nextState)
			if child.state != nil {
				nodes = append(nodes, child)
			}
		}

		if len(nodes) > 0 {
			if len(nodes) > 1 {
				sort.Sort(ByU(nodes))
			}
			current = nodes[0]
			nodes = nodes[1:]
		} else {
			current = nil
		}
	}

	if len(origin.children) > 0 {
		if len(origin.children) == 1 {
			return origin.children[0].Action, origin.children
		}
		sort.Sort(ByN(origin.children))
		best := origin.children[0].Action
		return best, origin.children
	}

	return nil, origin.children
}

type Node struct {
	parent   *Node
	Action   Action `json:"action"`
	state    ActionPending
	children []*Node
	W        int `json:"w"`
	N        int `json:"n"`
}

func (node *Node) backfill(W int) {
	node.W += W
	node.N += 1
	if node.parent != nil {
		node.parent.backfill(W)
	}
}

func (node *Node) U() float64 {
	return float64(node.W)/float64(node.N) +
		math.Sqrt(math.Log(float64(node.parent.N))/float64(node.N))
}

type ByU []*Node

func (byu ByU) Len() int {
	return len(byu)
}

func (byu ByU) Swap(i, j int) {
	byu[i], byu[j] = byu[j], byu[i]
}

func (byu ByU) Less(i, j int) bool {
	return byu[i].U() > byu[j].U()
}

type ByN []*Node

func (byn ByN) Len() int {
	return len(byn)
}

func (byn ByN) Swap(i, j int) {
	byn[i], byn[j] = byn[j], byn[i]
}

func (byn ByN) Less(i, j int) bool {
	return byn[i].N > byn[j].N
}

func nextActionPendingState(state State) ActionPending {
	for state != nil {
		if decision, ok := state.(ActionPending); ok {
			return decision
		}
		state = state.Next()
	}
	return nil
}

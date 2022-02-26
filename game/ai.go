package game

type AI interface {
	Think(state *GameState, agent *AgentState) *ActionMove
}

package game

import (
	"fmt"

	"github.com/statiolake/witness-counting-game/geom"
)

type GameState struct {
	field  *FieldState
	squads []*SquadState
	agents []*AgentState
}

type FieldState struct {
	minX, maxX, minY, maxY float64
	obsts                  []ObstructionState
}

type ObstructionState struct {
	lt geom.Coord
	rb geom.Coord
}

type SquadState struct {
	id    int
	name  string
	alive bool
}

type AgentState struct {
	id         int
	squad      int
	name       string
	kind       Kind
	pos        geom.Coord
	point      float64
	nextAction *ActionMove
}

type ActionMove struct {
	Dir geom.PolarVector
}

func NewSquadState(id int, squad SquadConfig) *SquadState {
	return &SquadState{
		id:    id,
		name:  squad.Name,
		alive: true,
	}
}

func NewAgentState(squadId int, id int, agent AgentConfig) *AgentState {
	return &AgentState{
		id:         id,
		squad:      squadId,
		name:       agent.Name,
		kind:       agent.Kind,
		pos:        agent.InitPos,
		point:      0,
		nextAction: &ActionMove{},
	}
}

func (a *AgentState) SetNextAction(action *ActionMove) {
	a.nextAction = action
}

func (a *AgentState) GetPos() geom.Coord {
	return a.pos
}

func NewGameState(config *GameConfig) GameState {
	obsts := []ObstructionState{}

	for _, obst := range config.Field.InitObsts {
		obsts = append(obsts, ObstructionState{
			lt: obst.LT,
			rb: obst.RB,
		})
	}

	field := &FieldState{
		minX:  config.Field.MinX,
		maxX:  config.Field.MaxX,
		minY:  config.Field.MinY,
		maxY:  config.Field.MaxY,
		obsts: obsts,
	}

	squads := []*SquadState{}
	agents := []*AgentState{}

	for squadId, squad := range config.Squads {
		squads = append(squads, NewSquadState(squadId, squad))
		for _, agent := range squad.Agents {
			// ID は squad ごとではなく完全にリストとして扱うので注意すべし
			id := len(agents)
			agents = append(agents, NewAgentState(squadId, id, agent))
		}
	}

	return GameState{
		field:  field,
		squads: squads,
		agents: agents,
	}
}

func (g *GameState) Step() error {
	for idx, agent := range g.agents {
		if idx != agent.id {
			panic(
				fmt.Sprintf(
					"internal error: index and id do not agree: %d and %d",
					idx,
					agent.id,
				),
			)
		}

		// 次の行動が設定されているのであれば適用する
		if _, err := g.ApplyActionFor(agent); err != nil {
			return fmt.Errorf("failed to execute a step: %w", err)
		}
	}

	return nil
}

func (g *GameState) ApplyActionFor(agent *AgentState) (bool, error) {
	if !g.isRegisteredAgent(agent) {
		return false, fmt.Errorf(
			"agent %s/%s is not registered",
			g.squads[agent.squad].name,
			agent.name,
		)
	}

	action := agent.nextAction
	if action == nil {
		// 移動しないが別にエラーではない
		return false, nil
	}

	// 実際に位置を移動する
	vec_dir := action.Dir.ToVector()
	new_pos := agent.pos.Add(vec_dir).AsCoord()

	if !g.field.MovableTo(agent, new_pos) {
		// 移動できないので何もしない
		return false, fmt.Errorf("cannot move to %s", new_pos.ToString())
	}

	agent.pos = new_pos
	return true, nil
}

func (g *GameState) isRegisteredAgent(agent *AgentState) bool {
	if agent.id >= len(g.agents) {
		return false
	}
	return g.agents[agent.id] == agent
}

func (f *FieldState) MovableTo(agent *AgentState, new_pos geom.Coord) bool {
	return (f.minX <= new_pos.X &&
		new_pos.X <= f.maxX &&
		f.minY <= new_pos.Y &&
		new_pos.Y <= f.maxY)
}

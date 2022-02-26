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

func (g *GameState) Step(ignoreError bool) error {
	if err := g.ProcessActions(true); err != nil {
		return err
	}

	g.MoveScore()

	return nil
}

func (g *GameState) ProcessActions(ignoreError bool) error {
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
		if _, err := agent.ApplyActionOn(g); !ignoreError && err != nil {
			return fmt.Errorf("failed to execute a step: %w", err)
		}
	}

	return nil
}

func (g *GameState) MoveScore() {
	// まずは各 Runner が何人から見られているかを数える (それによって一人の
	// Hunter がその Runner からもらえる得点がかわってくるので)
	watchers := make([][]*AgentState, len(g.agents))
	for _, a := range g.agents {
		if a.kind == Hunter {
			watchers[a.id] = g.GetCapturedRunners(a)
		} else if a.kind == Runner {
			watchers[a.id] = g.GetWatchingHunters(a)
		}
	}

	for _, a := range g.agents {
		var delta float64 = 0.0
		if a.kind == Hunter {
			// Hunter の報酬は各 Runner が提供してくれるスコアの和
			// 各 Runner はスコア 1.0 を見られているハンターへ等分する
			// TODO: 等分だと「わざと自分のハンターに見られることで相手に点数
			// が流出する量を減らす、という裏技が生まれてしまうので、よりいい
			// 感じのスコアを考えるべし
			for _, r := range watchers[a.id] {
				delta += 1.0 / float64(len(watchers[r.id]))
			}
		}
		if a.kind == Runner {
			// Runner は Hunter に対してスコアを提供する
			// 一人からでも見られている限り 1.0 を供出することになる
			if len(watchers[a.id]) > 0 {
				delta = -1.0
			}
		}
		a.point += delta
	}
}

func (g *GameState) GetWatchingHunters(agent *AgentState) []*AgentState {
	res := []*AgentState{}
	for _, hunter := range g.agents {
		if hunter.kind != Hunter {
			continue
		}

		if hunter.IsWatching(agent, g) {
			res = append(res, hunter)
		}
	}

	return res
}

func (g *GameState) GetCapturedRunners(hunter *AgentState) []*AgentState {
	res := []*AgentState{}
	for _, agent := range g.agents {
		if agent.kind != Runner {
			continue
		}

		if hunter.IsWatching(agent, g) {
			res = append(res, agent)
		}
	}

	return res
}

func (from *AgentState) IsWatching(to *AgentState, g *GameState) bool {
	// TODO: from と to の間に遮蔽物があるかどうかをチェックする
	return true
}

func (a *AgentState) ApplyActionOn(g *GameState) (bool, error) {
	if !a.isRegisteredAgentOn(g) {
		return false, fmt.Errorf(
			"agent %s/%s is not registered",
			g.squads[a.squad].name,
			a.name,
		)
	}

	action := a.nextAction
	if action == nil {
		// 移動しないが別にエラーではない
		return false, nil
	}

	// 実際に位置を移動する
	vec_dir := action.Dir.ToVector()
	new_pos := a.pos.Add(vec_dir).AsCoord()

	if !g.field.MovableTo(a, new_pos) {
		// 移動できないので何もしない
		return false, fmt.Errorf("cannot move to %s", new_pos.ToString())
	}

	a.pos = new_pos
	return true, nil
}

func (a *AgentState) isRegisteredAgentOn(g *GameState) bool {
	if a.id >= len(g.agents) {
		return false
	}
	return g.agents[a.id] == a
}

func (f *FieldState) MovableTo(agent *AgentState, new_pos geom.Coord) bool {
	return (f.minX <= new_pos.X &&
		new_pos.X <= f.maxX &&
		f.minY <= new_pos.Y &&
		new_pos.Y <= f.maxY)
}

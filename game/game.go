package game

import (
	"fmt"

	"github.com/statiolake/witness-counting-game/geom"
)

type Kind int

const (
	Hunter Kind = iota
	Runner
)

type Game struct {
	config *GameConfig
	field  *Field
	squads []*Squad
	agents []*Agent
}

type Field struct {
	rect  geom.Rect
	obsts []Obstruction
}

type Obstruction struct {
	rect geom.Rect
}

type Squad struct {
	id    int
	name  string
	alive bool
}

type Agent struct {
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

func NewGame(config *GameConfig) *Game {
	obsts := []Obstruction{}

	for _, obst := range config.field.obsts {
		obsts = append(obsts, Obstruction{rect: obst.rect})
	}

	field := &Field{
		rect:  config.field.rect,
		obsts: obsts,
	}

	squads := []*Squad{}
	agents := []*Agent{}

	for squadId, squad := range config.squads {
		squads = append(squads, NewSquad(squadId, squad))
		for _, agent := range squad.agents {
			// ID は squad ごとではなく完全にリストとして扱うので注意すべし
			id := len(agents)
			agents = append(agents, NewAgent(squadId, id, agent))
		}
	}

	return &Game{
		config: config,
		field:  field,
		squads: squads,
		agents: agents,
	}
}

func (g *Game) GetConfig() GameConfig {
	return *g.config
}

func NewSquad(id int, squad SquadConfig) *Squad {
	return &Squad{
		id:    id,
		name:  squad.name,
		alive: true,
	}
}

func NewAgent(squadId int, id int, agent AgentConfig) *Agent {
	return &Agent{
		id:         id,
		squad:      squadId,
		name:       agent.name,
		kind:       agent.kind,
		pos:        agent.initPos,
		point:      0,
		nextAction: &ActionMove{},
	}
}

func (a *Agent) GetPos() geom.Coord {
	return a.pos
}

func (g *Game) Step(ignoreError bool) error {
	if err := g.ProcessActions(true); err != nil {
		return err
	}

	g.MoveScore()

	return nil
}

func (g *Game) ProcessActions(ignoreError bool) error {
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

func (g *Game) MoveScore() {
	// まずは各 Runner が何人から見られているかを数える (それによって一人の
	// Hunter がその Runner からもらえる得点がかわってくるので)
	watchers := make([][]*Agent, len(g.agents))
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

func (g *Game) GetWatchingHunters(agent *Agent) []*Agent {
	res := []*Agent{}
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

func (g *Game) GetCapturedRunners(hunter *Agent) []*Agent {
	res := []*Agent{}
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

func (from *Agent) IsWatching(to *Agent, g *Game) bool {
	// TODO: from と to の間に遮蔽物があるかどうかをチェックする
	return true
}

func (a *Agent) ApplyActionOn(g *Game) (bool, error) {
	if !a.isRegisteredOn(g) {
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

func (f *Field) MovableTo(agent *Agent, new_pos geom.Coord) bool {
	return (f.rect.LT.X <= new_pos.X &&
		new_pos.X <= f.rect.RB.X &&
		f.rect.LT.Y <= new_pos.Y &&
		new_pos.Y <= f.rect.RB.Y)
}

func (a *Agent) isRegisteredOn(g *Game) bool {
	if a.id >= len(g.agents) {
		return false
	}
	return g.agents[a.id] == a
}

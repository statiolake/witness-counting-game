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
	Config GameConfig
	Field  Field
	Squads []Squad
	Agents []Agent
}

type Field struct {
	Rect  geom.Rect
	Obsts []Obstruction
}

type Obstruction struct {
	Rect geom.Rect
}

type Squad struct {
	Id    int
	Name  string
	Alive bool
}

type Agent struct {
	Id         int
	Squad      int
	Name       string
	Kind       Kind
	Pos        geom.Coord
	Point      float64
	NextAction *ActionMove
}

type ActionMove struct {
	Dir geom.PolarVector
}

func NewGame(config GameConfig) Game {
	obsts := []Obstruction{}

	for _, obst := range config.Field.Obsts {
		obsts = append(obsts, Obstruction(obst))
	}

	field := Field{
		Rect:  config.Field.Rect,
		Obsts: obsts,
	}

	squads := []Squad{}
	agents := []Agent{}

	for squadId, squad := range config.Squads {
		squads = append(squads, NewSquad(squadId, squad))
		for _, agent := range squad.Agents {
			// ID は squad ごとではなく完全にリストとして扱うので注意すべし
			id := len(agents)
			agents = append(agents, NewAgent(squadId, id, agent))
		}
	}

	return Game{
		Config: config,
		Field:  field,
		Squads: squads,
		Agents: agents,
	}
}

func (g *Game) Clone() Game {
	squads := make([]Squad, len(g.Squads))
	copy(squads, g.Squads)
	agents := make([]Agent, len(g.Agents))
	copy(agents, g.Agents)

	return Game{
		Config: g.Config.Clone(),
		Field:  g.Field.Clone(),
		Squads: squads,
		Agents: agents,
	}
}

func (f *Field) Clone() Field {
	obsts := make([]Obstruction, len(f.Obsts))
	copy(obsts, f.Obsts)
	return Field{
		Rect:  f.Rect,
		Obsts: obsts,
	}
}

func NewSquad(id int, squad SquadConfig) Squad {
	return Squad{
		Id:    id,
		Name:  squad.Name,
		Alive: true,
	}
}

func NewAgent(squadId int, id int, agent AgentConfig) Agent {
	return Agent{
		Id:         id,
		Squad:      squadId,
		Name:       agent.Name,
		Kind:       agent.Kind,
		Pos:        agent.InitPos,
		Point:      0,
		NextAction: &ActionMove{},
	}
}

func (g *Game) Step(ignoreError bool) error {
	if err := g.ProcessActions(true); err != nil {
		return err
	}

	g.MoveScore()

	return nil
}

func (g *Game) ProcessActions(ignoreError bool) error {
	for idx := range g.Agents {
		agent := &g.Agents[idx]
		if idx != agent.Id {
			panic(
				fmt.Sprintf(
					"internal error: index and id do not agree: %d and %d",
					idx,
					agent.Id,
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
	watchers := make([][]*Agent, len(g.Agents))
	for idx := range g.Agents {
		a := &g.Agents[idx]
		if a.Kind == Hunter {
			watchers[a.Id] = a.FindCapturedRunners(g)
		} else if a.Kind == Runner {
			watchers[a.Id] = a.FindWatchingHunters(g)
		}
	}

	for idx := range g.Agents {
		a := &g.Agents[idx]
		var delta float64 = 0.0
		if a.Kind == Hunter {
			// Hunter の報酬は各 Runner が提供してくれるスコアの和
			// 各 Runner はスコア 1.0 を見られているハンターへ等分する
			// TODO: 等分だと「わざと自分のハンターに見られることで相手に点数
			// が流出する量を減らす、という裏技が生まれてしまうので、よりいい
			// 感じのスコアを考えるべし
			for _, r := range watchers[a.Id] {
				delta += 1.0 / float64(len(watchers[r.Id]))
			}
		}
		if a.Kind == Runner {
			// Runner は Hunter に対してスコアを提供する
			// 一人からでも見られている限り 1.0 を供出することになる
			if len(watchers[a.Id]) > 0 {
				delta = -1.0
			}
		}
		a.Point += delta
	}
}

func (agent *Agent) FindWatchingHunters(g *Game) []*Agent {
	res := []*Agent{}
	for idx := range g.Agents {
		hunter := &g.Agents[idx]
		if hunter.Kind != Hunter {
			continue
		}

		if hunter.IsWatching(agent, g) {
			res = append(res, hunter)
		}
	}

	return res
}

func (hunter *Agent) FindCapturedRunners(g *Game) []*Agent {
	res := []*Agent{}
	for idx := range g.Agents {
		agent := &g.Agents[idx]
		if agent.Kind != Runner {
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
			g.Squads[a.Squad].Name,
			a.Name,
		)
	}

	action := a.NextAction
	if action == nil {
		// 移動しないが別にエラーではない
		return false, nil
	}

	// 実際に位置を移動する
	vec_dir := action.Dir.ToVector()
	new_pos := a.Pos.Add(vec_dir).AsCoord()

	if !g.Field.MovableTo(a, new_pos) {
		// 移動できないので何もしない
		return false, fmt.Errorf("cannot move to %s", new_pos.ToString())
	}

	a.Pos = new_pos
	return true, nil
}

func (f *Field) MovableTo(agent *Agent, new_pos geom.Coord) bool {
	return (f.Rect.LT.X <= new_pos.X &&
		new_pos.X <= f.Rect.RB.X &&
		f.Rect.LT.Y <= new_pos.Y &&
		new_pos.Y <= f.Rect.RB.Y)
}

func (a *Agent) isRegisteredOn(g *Game) bool {
	return a.Id < len(g.Agents) && &g.Agents[a.Id] == a
}

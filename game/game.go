package game

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/statiolake/witness-counting-game/geom"
)

type Kind int

const (
	Hunter Kind = iota
	Runner
)

type Game struct {
	Config        GameConfig
	Field         Field
	Squads        []Squad
	Agents        []Agent
	TimeRemaining int
}

// TODO: 渡す情報は考えるべし
type Knowledge struct {
	// フィールド情報
	Field Field
	// Squad の数
	NumSquads int
	// Agent の数
	NumAgents int
	// 自分
	Me Agent
	// 自分を見ている、あるいは自分から見える Agent (自分を含む)
	Watchers []Agent
}

type Field struct {
	Rect  geom.Rect
	Obsts []Obstruction
}

type Obstruction struct {
	Segment geom.Segment
}

type Squad struct {
	Id   int
	Name string
}

type Agent struct {
	Id        int // 全体で一意なエージェント番号
	InSquadId int // squad の中でのエージェント番号
	SquadId   int
	Name      string
	Kind      Kind
	Pos       geom.Coord
	Point     float64

	// ターンごとにリセットされる情報

	PointGains []PointGain
	Action     *ActionMove
}

// Runner が Hunter にポイントを提供するときは負の Gain として扱う。
// 「Runner は Hunter から -1.0 ポイントを獲得した」
type PointGain struct {
	// 本来なら *Agent などとして扱いたいところではあるが、再帰的な構造になる
	// と JSON と相性が悪いので諦める。
	//
	// Marshal をカスタムして JSON 出力時だけ整数値に置き換えることはできるか
	// もしれないが、その場合でも Unmarshal を実装するのは無理そう。PointGain
	// を Unmarshal するときに []Agent が必要ということになるため。
	AgentIdGainedFrom int
	Gain              float64
}

type ActionMove struct {
	Dir geom.PolarVector
}

func (config *GameConfig) BuildGame() Game {
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
		squads = append(squads, Squad{
			Id:   squadId,
			Name: squad.Name,
		})
		for inSquadId, agent := range squad.Agents {
			// ID は squad ごとではなく完全にリストとして扱うので注意すべし
			id := len(agents)
			agents = append(agents, Agent{
				Id:         id,
				InSquadId:  inSquadId,
				SquadId:    squadId,
				Name:       agent.Name,
				Kind:       agent.Kind,
				Pos:        agent.InitPos,
				Point:      0,
				PointGains: []PointGain{},
				Action:     nil,
			})
		}
	}

	return Game{
		Config:        config.Clone(),
		Field:         field,
		Squads:        squads,
		Agents:        agents,
		TimeRemaining: config.Time,
	}
}

func (g *Game) Clone() Game {
	squads := make([]Squad, len(g.Squads))
	copy(squads, g.Squads)
	agents := make([]Agent, 0, len(g.Agents))
	for idx := range g.Agents {
		agents = append(agents, g.Agents[idx].Clone())
	}

	return Game{
		Config:        g.Config.Clone(),
		Field:         g.Field.Clone(),
		Squads:        squads,
		Agents:        agents,
		TimeRemaining: g.TimeRemaining,
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

func (a *Agent) Clone() Agent {
	// ポインタなので NextAction を丁寧にコピーする必要がある
	var nextAction *ActionMove
	if a.Action != nil {
		action := *a.Action
		nextAction = &action
	}

	pointGains := make([]PointGain, len(a.PointGains))
	copy(pointGains, a.PointGains)

	return Agent{
		Id:         a.Id,
		InSquadId:  a.InSquadId,
		SquadId:    a.SquadId,
		Name:       a.Name,
		Kind:       a.Kind,
		Pos:        a.Pos,
		Point:      a.Point,
		PointGains: pointGains,
		Action:     nextAction,
	}
}

func (g *Game) GetKnowledgeFor(agent *Agent) Knowledge {
	field := g.Field.Clone()
	numSquads := len(g.Squads)
	numAgents := len(g.Agents)
	me := agent.Clone()
	var watchers []Agent
	for _, agent := range agent.FindWatchingAgents(g, nil, true) {
		watchers = append(watchers, agent.Clone())
	}

	return Knowledge{
		Field:     field,
		NumSquads: numSquads,
		NumAgents: numAgents,
		Me:        me,
		Watchers:  watchers,
	}
}

func (g *Game) IsFinished() bool {
	return g.TimeRemaining == 0
}

func (g *Game) CommitTurn() error {
	if g.IsFinished() {
		return fmt.Errorf("attempted to step an finished game")
	}

	// エラーは無視する (ゲーム中は基本的にエラーがあっても継続してほしい;
	// エージェントが不正な命令を出した場合の処理は無視)
	// TODO: 一発退場のような重たい罰にするべき？
	_ = g.processActions()

	g.moveScore()
	g.TimeRemaining--

	return nil
}

func (g *Game) DescribeAgent(agent *Agent) string {
	return fmt.Sprintf(
		"%s/%s",
		g.Squads[agent.SquadId].Name,
		agent.Name,
	)
}

func (agent *Agent) FindWatchingAgents(g *Game, targetKind *Kind, includeSquad bool) (res []*Agent) {
	for idx := range g.Agents {
		other := &g.Agents[idx]
		// ターゲットの種類が指定されている場合は一致していなければ終了
		if targetKind != nil && other.Kind != *targetKind {
			continue
		}

		// 自分の Squad を含まない場合は SquadId が同じものはスキップ
		if !includeSquad && other.SquadId == agent.SquadId {
			continue
		}

		if other.IsWatching(agent, g) {
			res = append(res, other)
		}
	}

	return
}

// 同じ Squad のメンバーを含めたい場合は includeSquad を true とする
func (runner *Agent) FindWatchingHunters(g *Game, includeSquad bool) []*Agent {
	if runner.Kind != Runner {
		panic("Find watching hunters called on Runner")
	}

	hunter := Hunter
	return runner.FindWatchingAgents(g, &hunter, includeSquad)
}

// 同じ Squad のメンバーを含めたい場合は includeSquad を true とする
func (hunter *Agent) FindWatchingRunners(g *Game, includeSquad bool) []*Agent {
	if hunter.Kind != Hunter {
		panic("Find watching runners called on Hunter")
	}

	runner := Runner
	return hunter.FindWatchingAgents(g, &runner, includeSquad)
}

func (from *Agent) IsWatching(to *Agent, g *Game) bool {
	for _, obst := range g.Field.Obsts {
		ftseg := geom.Segment{
			A: from.Pos,
			B: to.Pos,
		}

		// from-to を結ぶ線分と遮蔽物がぶつかるのであればこの二者はお互いに見
		// えていない。
		if obst.Segment.Crosses(ftseg) {
			return false
		}
	}

	return true
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

// エージェントの NextAction やポイント変動情報をリセットする
func (g *Game) StartTurn() {
	for idx := range g.Agents {
		g.Agents[idx].startTurn()
	}
}

func (a *Agent) startTurn() {
	a.Action = nil
	a.PointGains = []PointGain{}
}

func (g *Game) processActions() (errs error) {
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

		// 次の行動が設定されているのであれば適用する。
		if _, err := agent.applyActionOn(g); err != nil {
			errs = multierror.Append(errs, fmt.Errorf(
				"failed to execute a step in agent %s: %w",
				g.DescribeAgent(agent),
				err,
			))
		}
	}

	return
}

func (a *Agent) applyActionOn(g *Game) (bool, error) {
	if !a.isRegisteredOn(g) {
		return false, fmt.Errorf(
			"agent %s is not registered",
			g.DescribeAgent(a),
		)
	}

	action := a.Action
	if action == nil {
		// 移動しないが別にエラーではない
		return false, nil
	}

	// 移動速度は Speed までに制限する
	if action.Dir.R >= g.Config.Speed {
		action.Dir.R = g.Config.Speed
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

func (g *Game) moveScore() {
	// TODO: Hunter か Runner のどちらかから見るだけでよくない (それにより不整
	// 合も減らせる)
	// まずは各 Runner が何人から見られているかを数える (それによって一人の
	// Hunter がその Runner からもらえる得点がかわってくるので)
	watchers := make([][]*Agent, len(g.Agents))
	for idx := range g.Agents {
		a := &g.Agents[idx]
		switch a.Kind {
		case Hunter:
			watchers[a.Id] = a.FindWatchingRunners(g, false)
		case Runner:
			watchers[a.Id] = a.FindWatchingHunters(g, false)
		}
	}

	for idx := range g.Agents {
		a := &g.Agents[idx]
		delta := 0.0

		switch a.Kind {
		case Hunter:
			// Hunter の報酬は各 Runner が提供してくれるスコアの和
			// 各 Runner はスコア 1.0 を見られているハンターへ等分する
			for _, runner := range watchers[a.Id] {
				// この Hunter が見ている Runner から点数をもらう。もらえる点
				// 数は Runner が何人の Hunter から見られているかに依存する。
				// TODO: ここ単に等分で OK ？
				gain := 1.0 / float64(len(watchers[runner.Id]))
				delta += gain
				a.PointGains = append(a.PointGains, PointGain{
					AgentIdGainedFrom: runner.Id,
					Gain:              gain,
				})
			}

		case Runner:
			// Runner は Hunter に対してスコアを提供する
			// 一人からでも見られている限り 1.0 を供出することになる
			if len(watchers[a.Id]) > 0 {
				delta = -1.0

				// 見られている全員に対してスコアを供出する
				each := delta / float64(len(watchers[a.Id]))
				for _, hunter := range watchers[a.Id] {
					a.PointGains = append(a.PointGains, PointGain{
						AgentIdGainedFrom: hunter.Id,
						Gain:              each,
					})
				}
			}
		}

		a.Point += delta
	}
}

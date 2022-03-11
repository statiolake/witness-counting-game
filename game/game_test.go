package game

import (
	"fmt"
	"math"
	"testing"

	"github.com/statiolake/witness-counting-game/geom"
)

func TestStateSetup(t *testing.T) {
	t.Run("SquadsIndexIdAgreement", func(t *testing.T) {
		g := dummyGame()
		for idx, squad := range g.Squads {
			if idx != squad.ID {
				t.Errorf(
					"Squad id and index do not agree: %d and %d",
					idx,
					squad.ID,
				)
			}
		}
	})

	t.Run("AgentsIndexIdAgreement", func(t *testing.T) {
		g := dummyGame()
		for idx, agent := range g.Agents {
			if idx != agent.ID {
				t.Errorf(
					"Agent id and index do not agree: %d and %d",
					idx,
					agent.ID,
				)
			}
		}
	})

	t.Run("ConfigGameAgentKindAgreement", func(t *testing.T) {
		g := dummyGame()

		for idx := range g.Agents {
			agent := &g.Agents[idx]
			agentConfig := &g.Config.
				Squads[agent.SquadID].
				Agents[agent.InSquadID]

			if agent.Kind != agentConfig.Kind {
				t.Fatalf(
					"agent kind is not set: %v but config %v",
					agent.Kind, agentConfig.Kind,
				)
			}
		}
	})
}

func TestApplyActionOn(t *testing.T) {
	t.Run("RightAbove45", func(t *testing.T) {
		g := dummyGame()
		agent := &g.Agents[0]

		// 右上斜め 45 度にこのエージェントを動かしてみる
		agent.Action = &ActionMove{
			Dir: geom.NewPolarVector(1, math.Pi/4),
		}

		ok, err := agent.applyActionOn(&g)
		if err != nil {
			t.Fatalf("action didn't apply: %v", err)
		}

		if !ok {
			t.Fatalf("action not applied even though not nil")
		}

		expected := geom.NewCoord(1/math.Sqrt(2), 1/math.Sqrt(2))
		actual := agent.Pos
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}
	})

	t.Run("OutsideOfField", func(t *testing.T) {
		g := dummyGame()
		agent := &g.Agents[0]

		// まずは画面右端へ
		agent.Pos.X = g.Field.Rect.RB.X

		// さらに右へ
		agent.Action = &ActionMove{
			Dir: geom.NewPolarVector(1, 0),
		}

		ok, err := agent.applyActionOn(&g)

		if err == nil {
			t.Fatalf("too far moving accepted")
		}

		if ok {
			t.Fatalf("error but returned true")
		}
	})

	t.Run("DoNotMoveTooFast", func(t *testing.T) {
		g := dummyGame()
		agent := &g.Agents[0]

		// 遠くへ移動しようとしてみる
		agent.Action = &ActionMove{
			Dir: geom.NewPolarVector(1e5, 0),
		}

		ok, err := agent.applyActionOn(&g)

		if err != nil {
			t.Fatalf("too fast move caused error: %v", err)
		}

		if !ok {
			t.Fatalf("too fast move was cancelled")
		}

		// 実際にはスピード程度に抑えられていることを確認する
		expected := geom.NewCoord(1, 0)
		actual := agent.Pos
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}
	})

	t.Run("InvalidAgent", func(t *testing.T) {
		g := dummyGame()
		agent := Agent{
			ID:         0,
			InSquadID:  0,
			SquadID:    0,
			Name:       "",
			Kind:       0,
			Pos:        geom.NewCoord(0, 0),
			Point:      0,
			PointGains: []PointGain{},
			Action:     &ActionMove{Dir: geom.NewPolarVector(1, 0)},
		}

		ok, err := agent.applyActionOn(&g)

		if err == nil {
			t.Fatalf("invalid agent accepted")
		}

		if ok {
			t.Fatalf("error but returned true")
		}
	})

	t.Run("RunSpecifiedTime", func(t *testing.T) {
		count := 0

		g := dummyGame()
		for !g.IsFinished() {
			if count > g.Config.Time {
				t.Fatalf(
					"Turn elapsed without stopping game (%d turns (in %d turns) left)",
					g.TimeRemaining,
					g.Config.Time,
				)
			}
			g.StartTurn()
			g.Agents[0].Action = &ActionMove{
				Dir: geom.NewPolarVector(0.01, 0),
			}

			if err := g.CommitTurn(); err != nil {
				t.Fatalf("commit turn failed: %v", err)
			}

			count++
		}

		expected := geom.NewCoord(1, 0)
		actual := g.Agents[0].Pos
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}

		if count != g.Config.Time {
			t.Fatalf("Only %d turn (in %d turn) passed", count, g.Config.Time)
		}

		if !g.IsFinished() {
			t.Fatalf("IsFinished() is inconsistent")
		}
	})

	t.Run("RunSpecifiedTimeErrornous", func(t *testing.T) {
		count := 0

		g := dummyGame()
		for !g.IsFinished() {
			if count > g.Config.Time {
				t.Fatalf(
					"Turn elapsed without stopping game (%d turns (in %d turns) left)",
					g.TimeRemaining,
					g.Config.Time,
				)
			}
			g.StartTurn()
			tooFastSpeed := g.Config.Field.Rect.RB.X * 2 /
				float64(g.Config.Time)
			g.Agents[0].Action = &ActionMove{
				Dir: geom.NewPolarVector(tooFastSpeed, 0),
			}

			if err := g.CommitTurn(); err != nil {
				t.Fatalf("commit turn failed: %v", err)
			}

			count++
		}

		expected := geom.NewCoord(g.Config.Field.Rect.RB.X, 0)
		actual := g.Agents[0].Pos
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}

		if count != g.Config.Time {
			t.Fatalf("Only %d turn (in %d turn) passed", count, g.Config.Time)
		}

		if !g.IsFinished() {
			t.Fatalf("IsFinished() is inconsistent")
		}
	})
}

func TestTurn(t *testing.T) {
	t.Run("OneAgentHasAction", func(t *testing.T) {
		g := dummyGame()

		g.StartTurn()

		agent := &g.Agents[0]

		// 右上斜め 45 度にこのエージェントを動かしてみる
		agent.Action = &ActionMove{
			Dir: geom.NewPolarVector(1, math.Pi/4),
		}

		if err := g.CommitTurn(); err != nil {
			t.Fatalf("commit turn failed: %v", err)
		}

		expected := geom.NewCoord(1/math.Sqrt(2), 1/math.Sqrt(2))
		actual := agent.Pos
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}
	})

	t.Run("TwoAgentHaveAction", func(t *testing.T) {
		g := dummyGame()

		g.StartTurn()

		// 右上斜め 45 度にこのエージェントを動かしてみる
		ag0 := &g.Agents[0]
		ag0.Action = &ActionMove{
			Dir: geom.NewPolarVector(1, math.Pi/4),
		}

		// 左上斜め 45 度にこのエージェントを動かしてみる
		ag1 := &g.Agents[1]
		ag1.Action = &ActionMove{
			Dir: geom.NewPolarVector(1, 3*math.Pi/4),
		}

		if err := g.CommitTurn(); err != nil {
			t.Fatalf("commit turn failed: %v", err)
		}

		// ag0 の位置を確認
		expected := geom.NewCoord(1/math.Sqrt(2), 1/math.Sqrt(2))
		actual := ag0.Pos
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}

		// ag1 の位置を確認
		expected = geom.NewCoord(-1/math.Sqrt(2), 1/math.Sqrt(2))
		actual = ag1.Pos
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}

		// スコアを確認
		// 現状は障害物がないのでそのまま人数が出てくるはず。
		for idx := range g.Agents {
			a := &g.Agents[idx]
			if a.Kind == Hunter {
				// Hunter の場合、cfg.squads - 1 の数だけ Runner がいるから、
				// それぞれから 1/(cfg.squads - 1) だけもらっていて、結局 +1
				// squads - 1 なのはフレンドリーファイアがないため。
				if !eq(a.Point, 1.0) {
					t.Fatalf("expected %v but actual %v", 1.0, a.Point)
				}

				if len(a.PointGains) != len(g.Config.Squads)-1 {
					// いまは障害物がないので必ず Runner の数と Hunter が見て
					// いる Runner の数が一致しているはず。
					// ただしフレンドリーファイアはないので squads の数からは
					// 1 減る。
					t.Fatalf(
						"point was gained from %d runners but squads are %d",
						len(a.PointGains), len(g.Config.Squads),
					)
				}
			}

			if a.Kind == Runner {
				// Runner の場合、誰かには見られているので結局 -1
				if !eq(a.Point, -1.0) {
					t.Fatalf("expected %v but actual %v", 1.0, a.Point)
				}

				if len(a.PointGains) != len(g.Config.Squads)-1 {
					// いまは障害物がないので必ず Hunter の数と Runner が見ら
					// れた Hunter の数が一致しているはず。
					t.Fatalf(
						"point was given to %d hunters but squads are %d",
						len(a.PointGains), len(g.Config.Squads),
					)
				}
			}
		}
	})

	t.Run("Obstruction", func(t *testing.T) {
		// 次のような位置関係のゲームを作る。
		//
		//      +h |
		//      *h | *r
		//      +r |
		//
		// *: squad-01
		// +: squad-02
		//
		// このとき squad-01 は壁の向こう側にいるので得点を失わず、
		// squad-02 が一方的に得点を吸われる状況になっていてほしい
		g := DefaultGameConfig().
			WithSquadAdded(
				NewSquadConfig("squad-01").
					WithAgentAdded(
						NewAgentConfig("agent-01h", Hunter).
							WithInitPos(geom.NewCoord(-1, 0)),
					).
					WithAgentAdded(
						NewAgentConfig("agent-01r", Runner).
							WithInitPos(geom.NewCoord(1, 0)),
					),
			).
			WithSquadAdded(
				NewSquadConfig("squad-02").
					WithAgentAdded(
						NewAgentConfig("agent-02h", Hunter).
							WithInitPos(geom.NewCoord(-1, 1)),
					).
					WithAgentAdded(
						NewAgentConfig("agent-02r", Runner).
							WithInitPos(geom.NewCoord(-1, -1)),
					),
			).
			// 中央の遮蔽物を追加する
			WithFieldConfig(
				DefaultFieldConfig().
					WithObstructionAdded(
						ObstructionConfig{
							Segment: geom.NewSegment(
								geom.NewCoord(0, 2),
								geom.NewCoord(0, -2),
							),
						},
					),
			).
			BuildGame()

		// ターンを実行する
		g.StartTurn()
		if err := g.CommitTurn(); err != nil {
			t.Fatalf("commit turn failed: %v", err)
		}

		// 得点の変動を確認する
		hunter1 := &g.Agents[0]
		runner1 := &g.Agents[1]
		hunter2 := &g.Agents[2]
		runner2 := &g.Agents[3]

		asserts := []struct {
			name          string
			agent         *Agent
			numPointGains int
			point         float64
		}{
			{"hunter1", hunter1, 1, 1.0},
			{"hunter2", hunter2, 0, 0.0},
			{"runner1", runner1, 0, 0.0},
			{"runner2", runner2, 1, -1.0},
		}

		for _, assert := range asserts {
			if len(assert.agent.PointGains) != assert.numPointGains {
				t.Fatalf(
					"unexpected point gain for %s (%s): %v",
					assert.name, g.DescribeAgent(assert.agent),
					assert.agent.PointGains,
				)
			}

			if !eq(assert.agent.Point, assert.point) {
				t.Fatalf(
					"unexpected point for %s (%s): %v",
					assert.name, g.DescribeAgent(assert.agent),
					assert.agent.Point,
				)
			}
		}

		// Squad に保存されているスコアも確認
		totalPointsPerSquad := make([]float64, len(g.Squads))
		for idx := range g.Agents {
			a := &g.Agents[idx]
			totalPointsPerSquad[a.SquadID] += a.Point
		}

		for idx := range g.Squads {
			expected := totalPointsPerSquad[idx]
			actual := g.Squads[idx].TotalPoint
			if !eq(expected, actual) {
				t.Fatalf(
					"squad point does not agree: expected %f but actual %f",
					expected,
					actual,
				)
			}
		}
	})

	t.Run("Knowledge", func(t *testing.T) {
		// 次のような位置関係のゲームを作る。
		//
		//      +h |
		//      *h | *r
		//      +r |
		//
		// *: squad-01
		// +: squad-02
		//
		// このとき +h, *h, +r はお互いが分かるが、*r は誰のことも見えていない。
		g := DefaultGameConfig().
			WithSquadAdded(
				NewSquadConfig("squad-01").
					WithAgentAdded(
						NewAgentConfig("agent-01h", Hunter).
							WithInitPos(geom.NewCoord(-1, 0)),
					).
					WithAgentAdded(
						NewAgentConfig("agent-01r", Runner).
							WithInitPos(geom.NewCoord(1, 0)),
					),
			).
			WithSquadAdded(
				NewSquadConfig("squad-02").
					WithAgentAdded(
						NewAgentConfig("agent-02h", Hunter).
							WithInitPos(geom.NewCoord(-1, 1)),
					).
					WithAgentAdded(
						NewAgentConfig("agent-02r", Runner).
							WithInitPos(geom.NewCoord(-1, -1)),
					),
			).
			// 中央の遮蔽物を追加する
			WithFieldConfig(
				DefaultFieldConfig().
					WithObstructionAdded(
						ObstructionConfig{
							Segment: geom.NewSegment(
								geom.NewCoord(0, 2),
								geom.NewCoord(0, -2),
							),
						},
					),
			).
			BuildGame()

		{
			hunter1 := &g.Agents[0]
			knowledge := g.GetKnowledgeFor(hunter1)
			if len(knowledge.Watchers) != 3 {
				t.Fatalf(
					"invalid number of watchers: knows %v but should know %d",
					knowledge.Watchers, 3,
				)
			}
		}

		{
			runner1 := &g.Agents[1]
			knowledge := g.GetKnowledgeFor(runner1)
			if len(knowledge.Watchers) != 1 {
				t.Fatalf(
					"invalid number of watchers: knows %v but should know %d",
					knowledge.Watchers, 1,
				)
			}
		}
	})
}

func eq(a, b float64) bool {
	return math.Abs(a-b) < 1e-8
}

func dummyGame() Game {
	numSquads := 5
	config := DefaultGameConfig()
	for i := 1; i < numSquads; i++ {
		name := fmt.Sprintf("squad-%02d", i)
		agentBase := fmt.Sprintf("agent-%02d", i)
		config.WithSquadAdded(
			NewSquadConfig(name).
				WithAgentAdded(NewAgentConfig(agentBase+"h", Hunter)).
				WithAgentAdded(NewAgentConfig(agentBase+"r", Runner)),
		)
	}
	return config.BuildGame()
}

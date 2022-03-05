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
			if idx != squad.Id {
				t.Errorf(
					"Squad id and index do not agree: %d and %d",
					idx,
					squad.Id,
				)
			}
		}
	})

	t.Run("AgentsIndexIdAgreement", func(t *testing.T) {
		g := dummyGame()
		for idx, agent := range g.Agents {
			if idx != agent.Id {
				t.Errorf(
					"Agent id and index do not agree: %d and %d",
					idx,
					agent.Id,
				)
			}
		}
	})
}

func TestApplyActionFor(t *testing.T) {
	t.Run("RightAbove45", func(t *testing.T) {
		g := dummyGame()
		agent := &g.Agents[0]

		// 右上斜め 45 度にこのエージェントを動かしてみる
		agent.NextAction = &ActionMove{
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

		// 遠くへ移動しようとしてみる
		agent.NextAction = &ActionMove{
			Dir: geom.NewPolarVector(1e5, 0),
		}

		ok, err := agent.applyActionOn(&g)

		if err == nil {
			t.Fatalf("too far moving accepted")
		}

		if ok {
			t.Fatalf("error but returned true")
		}
	})

	t.Run("InvalidAgent", func(t *testing.T) {
		g := dummyGame()
		agent := Agent{
			Id:         0,
			SquadId:    0,
			Name:       "",
			Kind:       0,
			Pos:        geom.NewCoord(0, 0),
			Point:      0,
			NextAction: &ActionMove{Dir: geom.NewPolarVector(1, 0)},
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
			g.Agents[0].NextAction = &ActionMove{
				Dir: geom.NewPolarVector(0.01, 0),
			}
			g.Step()
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
			too_fast_speed := g.Config.Field.Rect.RB.X * 2 /
				float64(g.Config.Time)
			g.Agents[0].NextAction = &ActionMove{
				Dir: geom.NewPolarVector(too_fast_speed, 0),
			}
			g.Step()
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

func TestStep(t *testing.T) {
	t.Run("OneAgentHasAction", func(t *testing.T) {
		g := dummyGame()
		agent := &g.Agents[0]

		// 右上斜め 45 度にこのエージェントを動かしてみる
		agent.NextAction = &ActionMove{
			Dir: geom.NewPolarVector(1, math.Pi/4),
		}

		if err := g.Step(); err != nil {
			t.Fatalf("step failed: %v", err)
		}

		expected := geom.NewCoord(1/math.Sqrt(2), 1/math.Sqrt(2))
		actual := agent.Pos
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}
	})

	t.Run("TwoAgentHaveAction", func(t *testing.T) {
		g := dummyGame()

		// 右上斜め 45 度にこのエージェントを動かしてみる
		ag0 := &g.Agents[0]
		ag0.NextAction = &ActionMove{
			Dir: geom.NewPolarVector(1, math.Pi/4),
		}

		// 左上斜め 45 度にこのエージェントを動かしてみる
		ag1 := &g.Agents[1]
		ag1.NextAction = &ActionMove{
			Dir: geom.NewPolarVector(1, 3*math.Pi/4),
		}

		if err := g.Step(); err != nil {
			t.Fatalf("step failed: %v", err)
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
		for _, a := range g.Agents {
			if a.Kind == Hunter {
				// Hunter の場合、 cfg.squads の数だけ Runner がいるから、そ
				// れぞれから 1/cfg.squads だけもらっていて、結局 +1
				if !eq(a.Point, 1.0) {
					t.Fatalf("expected %v but actual %v", 1.0, a.Point)
				}
			}

			if a.Kind == Runner {
				// Runner の場合、誰かには見られているので結局 -1
				if !eq(a.Point, -1.0) {
					t.Fatalf("expected %v but actual %v", 1.0, a.Point)
				}
			}
		}
	})
}

func eq(a, b float64) bool {
	return math.Abs(a-b) < 1e-8
}

func createConfig() GameConfig {
	numSquads := 5

	squads := []SquadConfig{}
	for i := 1; i <= numSquads; i++ {
		name := fmt.Sprintf("squad-%02d", i)
		agentBase := fmt.Sprintf("agent-%02d", i)
		squads = append(squads, SquadConfig{
			Name: name,
			Agents: []AgentConfig{
				{
					Name:    agentBase + "h",
					Kind:    Hunter,
					InitPos: geom.NewCoord(0, 0),
				},
				{
					Name:    agentBase + "r",
					Kind:    Runner,
					InitPos: geom.NewCoord(0, 0),
				},
			},
		})
	}

	return GameConfig{
		Field: FieldConfig{
			Rect:  geom.NewRectFromPoints(-50.0, -50.0, 50.0, 50.0),
			Obsts: []ObstructionConfig{},
		},
		Squads: squads,
		Speed:  0,
		Time:   100,
	}
}

func dummyGame() Game {
	return NewGame(createConfig())
}

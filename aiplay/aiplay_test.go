package aiplay

import (
	"fmt"
	"math"
	"testing"

	"github.com/statiolake/witness-counting-game/game"
	"github.com/statiolake/witness-counting-game/geom"
)

type constAI struct {
	Dir geom.PolarVector
}

func (ai *constAI) Init(config game.GameConfig) error {
	return nil
}

func (ai *constAI) Think(
	knowledge game.Knowledge,
	agent game.Agent,
) (*game.ActionMove, error) {
	return &game.ActionMove{Dir: ai.Dir}, nil
}

func TestConstAI(t *testing.T) {
	t.Run("AIActionApplied", func(t *testing.T) {
		m := createAIPlay()

		if err := m.Step(); err != nil {
			t.Fatalf("failed to step: %v", err)
		}

		for idx := range m.Game.Agents {
			agent := &m.Game.Agents[idx]
			expected := geom.NewCoord(1.0, 0.0)
			actual := agent.Pos
			if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
				t.Fatalf("expected %v but actual %v", expected, actual)
			}
		}
	})
}

func TestStepAll(t *testing.T) {
	t.Run("StepAll", func(t *testing.T) {
		g := createAIPlay()
		snapshots, err := g.StepAll()
		if err != nil {
			t.Fatalf("StepAll() failed: %v", err)
		}

		if len(snapshots) != g.Game.Config.Time+1 {
			t.Fatalf(
				"only %d snapshots (in %d turn) returned",
				len(snapshots), g.Game.Config.Time,
			)
		}

		final := &snapshots[len(snapshots)-1]

		if !final.IsFinished() {
			t.Fatalf("last snapshot is not finished game")
		}

		expected := geom.NewCoord(
			math.Min(1.0*float64(g.Game.Config.Time), g.Game.Field.Rect.RB.X), 0,
		)
		actual := g.Game.Agents[0].Pos
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}

		// スコアを確認
		for idx := range g.Game.Agents {
			agent := &g.Game.Agents[idx]
			finalAgent := &final.Agents[idx]

			if eq(agent.Point, 0.0) {
				t.Fatalf(
					"no point move on agent %s",
					g.Game.DescribeAgent(agent),
				)
			}

			if !eq(agent.Point, finalAgent.Point) {
				t.Fatalf(
					"last game and final snapshot does not agree on Point"+
						" of agent %s: %f vs %f",
					g.Game.DescribeAgent(agent),
					agent.Point, finalAgent.Point,
				)
			}
		}
	})
}

func eq(a, b float64) bool {
	return math.Abs(a-b) < 1e-8
}

func createAIPlay() AIPlay {
	numSquads := 5

	config := DefaultAIPlayConfig()
	for i := 1; i <= numSquads; i++ {
		name := fmt.Sprintf("squad-%02d", i)
		agentBase := fmt.Sprintf("agent-%02d", i)

		squad := NewSquadConfig(name)
		squad.AddAgent(
			game.NewAgentConfig(agentBase+"h", game.Hunter),
			&constAI{Dir: geom.NewPolarVector(1, 0)},
		)

		squad.AddAgent(
			game.NewAgentConfig(agentBase+"r", game.Runner),
			&constAI{Dir: geom.NewPolarVector(1, 0)},
		)

		config.AddSquad(squad)
	}

	return NewAIPlay(config)
}

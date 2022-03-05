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

func (ai *constAI) Think(knowledge game.Knowledge, agent *game.Agent) error {
	agent.NextAction = &game.ActionMove{
		Dir: ai.Dir,
	}
	return nil
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

		if !snapshots[len(snapshots)-1].IsFinished() {
			t.Fatalf("last snapshot is not finished game")
		}

		expected := geom.NewCoord(
			math.Min(1.0*float64(g.Game.Config.Time), g.Game.Field.Rect.RB.X), 0,
		)
		actual := g.Game.Agents[0].Pos
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}
	})
}

func eq(a, b float64) bool {
	return math.Abs(a-b) < 1e-8
}

func createAIPlay() AIPlay {
	return NewAIPlay(createConfig())
}

func createConfig() AIPlayConfig {
	numSquads := 5

	config := AIPlayConfig{
		GameConfig: createGameConfig(),
		AIs:        []AI{},
	}

	for i := 1; i <= numSquads; i++ {
		name := fmt.Sprintf("squad-%02d", i)
		squad := SquadConfig{
			Name:   name,
			Agents: []game.AgentConfig{},
			AIs:    []AI{},
		}

		agentBase := fmt.Sprintf("agent-%02d", i)
		squad.AddAgent(
			game.AgentConfig{
				Name:    agentBase + "h",
				Kind:    game.Hunter,
				InitPos: geom.NewCoord(0, 0),
			},
			&constAI{Dir: geom.NewPolarVector(1, 0)},
		)

		squad.AddAgent(
			game.AgentConfig{
				Name:    agentBase + "r",
				Kind:    game.Runner,
				InitPos: geom.NewCoord(0, 0),
			},
			&constAI{Dir: geom.NewPolarVector(1, 0)},
		)
		config.AddSquad(squad)
	}

	return config
}

func createGameConfig() game.GameConfig {
	return game.GameConfig{
		Field: game.FieldConfig{
			Rect:  geom.NewRectFromPoints(-50.0, -50.0, 50.0, 50.0),
			Obsts: []game.ObstructionConfig{},
		},
		Squads: []game.SquadConfig{},
		Speed:  1,
		Time:   100,
	}
}

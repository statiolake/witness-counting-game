package game

import (
	"fmt"
	"math"
	"testing"

	"github.com/statiolake/witness-counting-game/geom"
)

func TestStateSetup(t *testing.T) {
	state := dummyGameState()

	t.Run("SquadsIndexIdAgreement", func(t *testing.T) {
		for idx, squad := range state.squads {
			if idx != squad.id {
				t.Errorf(
					"Squad id and index do not agree: %d and %d",
					idx,
					squad.id,
				)
			}
		}
	})

	t.Run("AgentsIndexIdAgreement", func(t *testing.T) {
		for idx, agent := range state.agents {
			if idx != agent.id {
				t.Errorf(
					"Agent id and index do not agree: %d and %d",
					idx,
					agent.id,
				)
			}
		}
	})
}

func TestApplyActionFor(t *testing.T) {
	state := dummyGameState()
	t.Run("RightAbove45", func(t *testing.T) {
		agent := state.agents[0]

		// 右上斜め 45 度にこのエージェントを動かしてみる
		agent.SetNextAction(&ActionMove{
			Dir: geom.NewPolarVector(1, math.Pi/4),
		})

		ok, err := state.ApplyActionFor(state.agents[0])
		if err != nil {
			t.Fatalf("action didn't apply: %v", err)
		}

		if !ok {
			t.Fatalf("action not applied even though not nil")
		}

		expected := geom.NewCoord(1/math.Sqrt(2), 1/math.Sqrt(2))
		actual := agent.GetPos()
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}
	})

	t.Run("OutsideOfField", func(t *testing.T) {
		agent := state.agents[0]

		// 遠くへ移動しようとしてみる
		agent.SetNextAction(&ActionMove{
			Dir: geom.NewPolarVector(1e5, 0),
		})

		ok, err := state.ApplyActionFor(state.agents[0])

		if err == nil {
			t.Fatalf("too far moving accepted")
		}

		if ok {
			t.Fatalf("error but returned true")
		}
	})

	t.Run("InvalidAgent", func(t *testing.T) {
		agent := AgentState{
			id:         0,
			squad:      0,
			name:       "",
			kind:       0,
			pos:        geom.NewCoord(0, 0),
			point:      0,
			nextAction: &ActionMove{Dir: geom.NewPolarVector(1, 0)},
		}

		ok, err := state.ApplyActionFor(&agent)

		if err == nil {
			t.Fatalf("invalid agent accepted")
		}

		if ok {
			t.Fatalf("error but returned true")
		}
	})
}

func TestStep(t *testing.T) {
	t.Run("OneAgentHasAction", func(t *testing.T) {
		state := dummyGameState()
		agent := state.agents[0]

		// 右上斜め 45 度にこのエージェントを動かしてみる
		agent.SetNextAction(&ActionMove{
			Dir: geom.NewPolarVector(1, math.Pi/4),
		})

		if err := state.Step(); err != nil {
			t.Fatalf("step failed: %v", err)
		}

		expected := geom.NewCoord(1/math.Sqrt(2), 1/math.Sqrt(2))
		actual := agent.GetPos()
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}
	})

	t.Run("TwoAgentHaveAction", func(t *testing.T) {
		state := dummyGameState()

		// 右上斜め 45 度にこのエージェントを動かしてみる
		ag0 := state.agents[0]
		ag0.nextAction = &ActionMove{
			Dir: geom.NewPolarVector(1, math.Pi/4),
		}

		// 左上斜め 45 度にこのエージェントを動かしてみる
		ag1 := state.agents[1]
		ag1.nextAction = &ActionMove{
			Dir: geom.NewPolarVector(1, 3*math.Pi/4),
		}

		if err := state.Step(); err != nil {
			t.Fatalf("step failed: %v", err)
		}

		// ag0 の位置を確認
		expected := geom.NewCoord(1/math.Sqrt(2), 1/math.Sqrt(2))
		actual := ag0.GetPos()
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}

		// ag1 の位置を確認
		expected = geom.NewCoord(-1/math.Sqrt(2), 1/math.Sqrt(2))
		actual = ag1.GetPos()
		if !eq(actual.X, expected.X) || !eq(actual.Y, expected.Y) {
			t.Fatalf("expected %v but actual %v", expected, actual)
		}
	})
}

func eq(a, b float64) bool {
	return math.Abs(a-b) < 1e-8
}

func dummyConfig() GameConfig {
	numSquads := 5

	field := FieldConfig{
		MinX:      -10.0,
		MaxX:      10.0,
		MinY:      -10.0,
		MaxY:      10.0,
		InitObsts: []ObstructionConfig{},
	}

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
		Field:  field,
		Squads: squads,
		Speed:  0,
	}
}

func dummyGameState() GameState {
	config := dummyConfig()
	return NewGameState(&config)
}

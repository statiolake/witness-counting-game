package game

import (
	"testing"

	"github.com/statiolake/witness-counting-game/geom"
)

func TestGameInitialize(t *testing.T) {
	_ = newGame()
}

func createConfig() *GameConfig {
	return &GameConfig{
		Field: FieldConfig{
			MinX:      -50.0,
			MaxX:      50.0,
			MinY:      -50.0,
			MaxY:      50.0,
			InitObsts: []ObstructionConfig{},
		},
		Squads: []SquadConfig{
			{
				Name: "squad01",
				Agents: []AgentConfig{
					{
						Name:    "agent01h",
						Kind:    Hunter,
						InitPos: geom.NewCoord(0, 0),
						AI:      nil,
					},
					{
						Name:    "agent01r",
						Kind:    Runner,
						InitPos: geom.NewCoord(0, 0),
						AI:      nil,
					},
				},
			},
		},
		Speed: 0,
	}
}

func newGame() *Game {
	return NewGame(createConfig())
}

package main

import (
	"encoding/json"
	"fmt"

	"github.com/statiolake/witness-counting-game/aiplay"
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

func main() {
	m := createAIPlay()

	snapshots, err := m.StepAll()

	if err != nil {
		panic(err)
	}

	result, err := json.MarshalIndent(snapshots, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(result))
}

func createAIPlay() aiplay.AIPlay {
	return aiplay.NewAIPlay(createConfig())
}

func createConfig() aiplay.AIPlayConfig {
	numSquads := 5

	config := aiplay.AIPlayConfig{
		GameConfig: createGameConfig(),
		AIs:        []aiplay.AI{},
	}

	for i := 1; i <= numSquads; i++ {
		name := fmt.Sprintf("squad-%02d", i)
		squad := aiplay.SquadConfig{
			Name:   name,
			Agents: []game.AgentConfig{},
			AIs:    []aiplay.AI{},
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

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

func (ai *constAI) Think(
	knowledge game.Knowledge,
	agent game.Agent,
) (*game.ActionMove, error) {
	return &game.ActionMove{Dir: ai.Dir}, nil
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
	numSquads := 5

	config := aiplay.DefaultAIPlayConfig()

	for i := 1; i <= numSquads; i++ {
		name := fmt.Sprintf("squad-%02d", i)
		agentBase := fmt.Sprintf("agent-%02d", i)

		squad := aiplay.NewSquadConfig(name)
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

	return aiplay.NewAIPlay(config)
}

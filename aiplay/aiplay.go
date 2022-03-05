package aiplay

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/statiolake/witness-counting-game/game"
)

type AI interface {
	Init(config game.GameConfig) error
	Think(knowledge game.Knowledge, agent *game.Agent) error
}

type AIPlay struct {
	Game game.Game
	AIs  []AI
}

func NewAIPlay(config AIPlayConfig) AIPlay {
	game := game.NewGame(config.GameConfig)
	return AIPlay{
		Game: game,
		AIs:  config.AIs,
	}
}

func (g *AIPlay) Step() error {
	if err := g.decideActions(); err != nil {
		return fmt.Errorf("failed to decide actions: %w", err)
	}

	if err := g.Game.Step(); err != nil {
		return fmt.Errorf("failed to step game: %w", err)
	}

	return nil
}

func (g *AIPlay) decideActions() (errs error) {
	for idx := range g.AIs {
		agent := &g.Game.Agents[idx]
		ai := g.AIs[idx]
		err := ai.Think(g.Game.GetKnowledgeFor(agent), agent)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf(
				"agent %s: %w",
				g.Game.DescribeAgent(agent),
				err,
			))
		}
	}

	return
}

package aiplay

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/statiolake/witness-counting-game/game"
)

type AI interface {
	Init(config game.GameConfig) error
	Think(knowledge game.Knowledge, agent game.Agent) (*game.ActionMove, error)
}

type AIPlay struct {
	Game game.Game
	AIs  []AI
}

func (config *AIPlayConfig) BuildAIPlay() AIPlay {
	game := config.GameConfig.BuildGame()
	return AIPlay{
		Game: game,
		AIs:  config.AIs,
	}
}

func (g *AIPlay) StepAll() (snapshots []game.Game, err error) {
	snapshots = append(snapshots, g.Game.Clone())
	for !g.Game.IsFinished() {
		if err = g.Step(); err != nil {
			return
		}
		snapshots = append(snapshots, g.Game.Clone())
	}
	return
}

func (g *AIPlay) Step() error {
	g.Game.StartTurn()

	if err := g.decideActions(); err != nil {
		return fmt.Errorf("failed to decide actions: %w", err)
	}

	if err := g.Game.CommitTurn(); err != nil {
		return fmt.Errorf("failed to step game: %w", err)
	}

	return nil
}

func (g *AIPlay) decideActions() (errs error) {
	for idx := range g.AIs {
		agent := &g.Game.Agents[idx]
		ai := g.AIs[idx]
		action, err := ai.Think(
			g.Game.GetKnowledgeFor(agent),
			agent.Clone(),
		)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf(
				"agent %s: %w",
				g.Game.DescribeAgent(agent),
				err,
			))
		} else {
			agent.Action = action
		}
	}

	return
}

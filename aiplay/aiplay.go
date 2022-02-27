package aiplay

import "github.com/statiolake/witness-counting-game/game"

type AIPlay struct {
	game *game.Game
	ais  []AI
}

func NewAIPlay(config AIPlayConfig) AIPlay {
	for _, ai := range config.ais {
		ai.Init(&config.gameConfig)
	}

	return AIPlay{
		game: game.NewGame(&config.gameConfig),
		ais:  config.ais,
	}
}

func (g *AIPlay) Step() bool {
	// 各 AI に設定させる
	for id, agent := range g.game.GetAgents() {
		g.ais[id].Think(Knowledge{}, agent)
	}

	return g.game.Step()
}

func (g *AIPlay) StepAll() {
	for g.Step() {
	}
}

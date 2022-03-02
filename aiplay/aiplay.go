package aiplay

import "github.com/statiolake/witness-counting-game/game"

type AIPlay struct {
	game *game.Game
	ais  []AI
}

func NewAIPlay(config AIPlayConfig) AIPlay {
	game := game.NewGame(config.GameConfig)
	return AIPlay{
		game: &game,
		ais:  config.AIs,
	}
}

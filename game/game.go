package game

type Kind int

const (
	Hunter Kind = iota
	Runner
)

type Game struct {
	config GameConfig
	state  GameState
}

func NewGame(config *GameConfig) *Game {
	return &Game{
		state: NewGameState(config),
	}
}

func (g *Game) GetConfig() GameConfig {
	return g.config
}

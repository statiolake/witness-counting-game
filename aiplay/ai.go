package aiplay

import "github.com/statiolake/witness-counting-game/game"

type AI interface {
	Init(config game.GameConfig)
	Think(knowledge Knowledge, agent *game.Agent)
}

// TODO: 渡す情報は考えるべし
type Knowledge struct {
}

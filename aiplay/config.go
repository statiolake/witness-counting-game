package aiplay

import "github.com/statiolake/witness-counting-game/game"

type AIPlayConfig struct {
	gameConfig game.GameConfig
	ais        []AI
}

type SquadConfig struct {
	name   string
	agents []game.AgentConfig
	ais    []AI
}

func NewAIPlayConfig(field game.FieldConfig) AIPlayConfig {
	return AIPlayConfig{
		gameConfig: game.NewGameConfig(field),
		ais:        []AI{},
	}
}

func NewSquadConfig(name string) SquadConfig {
	return SquadConfig{
		name:   name,
		agents: []game.AgentConfig{},
		ais:    []AI{},
	}
}

func (c *AIPlayConfig) AddSquad(squad SquadConfig) {
	squadConfig := game.NewSquadConfig(squad.name)

	for _, agent := range squad.agents {
		squadConfig.WithAgent(agent)
	}

	for _, ai := range squad.ais {
		c.ais = append(c.ais, ai)
	}

	c.gameConfig.WithSquad(squadConfig)
}

func (c *SquadConfig) AddAgent(agent game.AgentConfig, ai AI) {
	c.agents = append(c.agents, agent)
	c.ais = append(c.ais, ai)
}

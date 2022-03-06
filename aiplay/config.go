package aiplay

import "github.com/statiolake/witness-counting-game/game"

type AIPlayConfig struct {
	GameConfig game.GameConfig
	AIs        []AI
}

type SquadConfig struct {
	Name   string
	Agents []game.AgentConfig
	AIs    []AI
}

func DefaultAIPlayConfig() AIPlayConfig {
	return AIPlayConfig{
		GameConfig: game.DefaultGameConfig(),
		AIs:        []AI{},
	}
}

func (c *AIPlayConfig) AddSquad(squad SquadConfig) *AIPlayConfig {
	squadConfig := game.SquadConfig{
		Name:   squad.Name,
		Agents: []game.AgentConfig{},
	}

	for _, agent := range squad.Agents {
		squadConfig.Agents = append(squadConfig.Agents, agent)
	}

	for _, ai := range squad.AIs {
		c.AIs = append(c.AIs, ai)
	}

	c.GameConfig.Squads = append(c.GameConfig.Squads, squadConfig)

	return c
}

func NewSquadConfig(name string) SquadConfig {
	return SquadConfig{
		Name:   name,
		Agents: []game.AgentConfig{},
		AIs:    []AI{},
	}
}

func (c *SquadConfig) AddAgent(agent game.AgentConfig, ai AI) *SquadConfig {
	c.Agents = append(c.Agents, agent)
	c.AIs = append(c.AIs, ai)
	return c
}

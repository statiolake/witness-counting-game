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

func (c *AIPlayConfig) AddSquad(squad *SquadConfig) {
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
}

func (c *SquadConfig) AddAgent(agent game.AgentConfig, ai AI) {
	c.Agents = append(c.Agents, agent)
	c.AIs = append(c.AIs, ai)
}

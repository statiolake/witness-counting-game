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

func DefaultAIPlayConfig() *AIPlayConfig {
	return &AIPlayConfig{
		GameConfig: *game.DefaultGameConfig(),
		AIs:        []AI{},
	}
}

func (c *AIPlayConfig) WithSquadAdded(squad *SquadConfig) *AIPlayConfig {
	c.AIs = append(c.AIs, squad.AIs...)

	squadConfig := game.NewSquadConfig(squad.Name)
	for idx := range squad.Agents {
		squadConfig.WithAgentAdded(&squad.Agents[idx])
	}
	c.GameConfig.WithSquadAdded(squadConfig)

	return c
}

func NewSquadConfig(name string) *SquadConfig {
	return &SquadConfig{
		Name:   name,
		Agents: []game.AgentConfig{},
		AIs:    []AI{},
	}
}

func (c *SquadConfig) WithAgentAdded(agent *game.AgentConfig, ai AI) *SquadConfig {
	c.Agents = append(c.Agents, agent.Clone())
	c.AIs = append(c.AIs, ai)
	return c
}

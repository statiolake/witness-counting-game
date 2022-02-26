package game

import "github.com/statiolake/witness-counting-game/geom"

type GameConfig struct {
	field  FieldConfig
	squads []SquadConfig
	speed  float64
}

type FieldConfig struct {
	rect  geom.Rect
	obsts []ObstructionConfig
}

type ObstructionConfig struct {
	rect geom.Rect
}

type SquadConfig struct {
	name   string
	agents []AgentConfig
}

type AgentConfig struct {
	name    string
	kind    Kind
	initPos geom.Coord
}

func NewGameConfig(field FieldConfig, speed float64) GameConfig {
	return GameConfig{
		field:  field,
		squads: []SquadConfig{},
		speed:  speed,
	}
}

func NewFieldConfig(width, height float64) FieldConfig {
	return FieldConfig{
		rect:  geom.NewRectFromPoints(-width/2, -height/2, width/2, height/2),
		obsts: []ObstructionConfig{},
	}
}

func NewSquadConfig(name string) SquadConfig {
	return SquadConfig{
		name:   name,
		agents: []AgentConfig{},
	}
}

func NewAgentConfig(name string, kind Kind, initPos geom.Coord) AgentConfig {
	return AgentConfig{
		name:    name,
		kind:    kind,
		initPos: initPos,
	}
}

func (c *GameConfig) GetFieldConfig() *FieldConfig {
	return &c.field
}

func (c *GameConfig) GetSquadConfigs() []SquadConfig {
	return c.squads
}

func (c *GameConfig) GetSpeed() float64 {
	return c.speed
}

func (c *FieldConfig) GetRect() geom.Rect {
	return c.rect
}

func (c *FieldConfig) GetObstructionConfigs() []ObstructionConfig {
	return c.obsts
}

func (c *ObstructionConfig) GetRect() geom.Rect {
	return c.rect
}

func (c *FieldConfig) AddObstructionConfig(obst ObstructionConfig) {
	c.obsts = append(c.obsts, obst)
}

func (c *GameConfig) AddSquad(squad SquadConfig) {
	c.squads = append(c.squads, squad)
}

func (c *SquadConfig) AddAgent(agent AgentConfig) {
	c.agents = append(c.agents, agent)
}

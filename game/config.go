package game

import "github.com/statiolake/witness-counting-game/geom"

type GameConfig struct {
	Field  FieldConfig
	Squads []SquadConfig
	Speed  float64
	Time   int
}

type FieldConfig struct {
	Rect  geom.Rect
	Obsts []ObstructionConfig
}

type ObstructionConfig struct {
	Line geom.Line
}

type SquadConfig struct {
	Name   string
	Agents []AgentConfig
}

type AgentConfig struct {
	Name    string
	Kind    Kind
	InitPos geom.Coord
}

func DefaultGameConfig() GameConfig {
	return GameConfig{
		Field:  DefaultFieldConfig(),
		Squads: []SquadConfig{},
		Speed:  1.0,
		Time:   100,
	}
}

func (c *GameConfig) WithFieldConfig(field *FieldConfig) *GameConfig {
	c.Field = *field
	return c
}

func (c *GameConfig) AddSquad(squad SquadConfig) *GameConfig {
	c.Squads = append(c.Squads, squad)
	return c
}

func (c *GameConfig) WithSpeed(speed float64) *GameConfig {
	c.Speed = speed
	return c
}

func (c *GameConfig) WithTime(time int) *GameConfig {
	c.Time = time
	return c
}

func DefaultFieldConfig() FieldConfig {
	return FieldConfig{
		Rect:  geom.NewRectFromPoints(-50.0, -50.0, 50.0, 50.0),
		Obsts: []ObstructionConfig{},
	}
}

func (c *FieldConfig) WithRect(rect geom.Rect) *FieldConfig {
	c.Rect = rect
	return c
}

func (c *FieldConfig) AddObstruction(obst ObstructionConfig) *FieldConfig {
	c.Obsts = append(c.Obsts, obst)
	return c
}

func NewSquadConfig(name string) SquadConfig {
	return SquadConfig{
		Name:   name,
		Agents: []AgentConfig{},
	}
}

func (c *SquadConfig) AddAgent(agent AgentConfig) *SquadConfig {
	c.Agents = append(c.Agents, agent)
	return c
}

func NewAgentConfig(name string, kind Kind) AgentConfig {
	return AgentConfig{
		Name:    name,
		Kind:    kind,
		InitPos: geom.NewCoord(0, 0),
	}
}

func (c *AgentConfig) WithInitPos(pos geom.Coord) *AgentConfig {
	c.InitPos = pos
	return c
}

func (c *GameConfig) Clone() GameConfig {
	squads := make([]SquadConfig, len(c.Squads))
	copy(squads, c.Squads)
	return GameConfig{
		Field:  c.Field.Clone(),
		Squads: squads,
		Speed:  c.Speed,
	}
}

func (c *FieldConfig) Clone() FieldConfig {
	obsts := make([]ObstructionConfig, len(c.Obsts))
	copy(obsts, c.Obsts)
	return FieldConfig{
		Rect:  c.Rect,
		Obsts: obsts,
	}
}

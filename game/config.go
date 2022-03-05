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

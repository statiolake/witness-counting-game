package game

import "github.com/statiolake/witness-counting-game/geom"

type GameConfig struct {
	Field  FieldConfig
	Squads []SquadConfig
	Speed  float64
}

type FieldConfig struct {
	MinX, MaxX, MinY, MaxY float64
	InitObsts              []ObstructionConfig
}

type ObstructionConfig struct {
	LT, RB geom.Coord
}

type SquadConfig struct {
	Name   string
	Agents []AgentConfig
}

type AgentConfig struct {
	Name    string
	Kind    Kind
	InitPos geom.Coord
	AI      AI
}

package game

import "github.com/statiolake/witness-counting-game/geom"

type Snapshot struct {
	Field         Field
	Squads        []SquadSnapshot
	Agents        []AgentSnapshot
	TimeRemaining int
}

type SquadSnapshot struct {
	Id          int
	Name        string
	TotalPoints float64
	PointsGain  float64
	PointsLost  float64
}

func squadSnapshotFromSquad(squad *Squad) SquadSnapshot {
	panic("not implemented")
}

type AgentSnapshot struct {
	Id         int
	Squad      int
	Name       string
	Kind       Kind
	Pos        geom.Coord
	Points     float64
	PointDelta []PointDelta
}

func agentSnapshotFromAgent(agent *Agent) AgentSnapshot {
	panic("not implemented")
}

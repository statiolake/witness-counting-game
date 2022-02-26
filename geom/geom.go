package geom

import (
	"fmt"
	"math"
)

type Vector struct {
	X, Y float64
}

type Coord struct {
	Vector
}

type PolarVector struct {
	R, T float64
}

type PolarCoord struct {
	PolarVector
}

type Rect struct {
	LT, RB Coord
}

type Direction struct {
	t float64
}

func NewVector(x, y float64) Vector {
	return Vector{
		X: x,
		Y: y,
	}
}

func (v Vector) ToString() string {
	return fmt.Sprintf("(%f, %f)", v.X, v.Y)
}

func (v Vector) AsCoord() Coord {
	return Coord{
		Vector: v,
	}
}

func (a Vector) Add(b Vector) Vector {
	return NewVector(
		a.X+b.X,
		a.Y+b.Y,
	)
}

func NewCoord(x, y float64) Coord {
	return Coord{
		Vector: NewVector(x, y),
	}
}

func NewRect(lt, rb Coord) Rect {
	return Rect{
		LT: lt,
		RB: rb,
	}
}

func NewRectFromPoints(minX, minY, maxX, maxY float64) Rect {
	return Rect{
		LT: NewCoord(minX, minY),
		RB: NewCoord(maxX, maxY),
	}
}

func (v Coord) AsVector() Vector {
	return v.Vector
}

func NewPolarVector(r, t float64) PolarVector {
	return PolarVector{
		R: r,
		T: t,
	}
}

func (p PolarVector) ToVector() Vector {
	return NewVector(
		p.R*math.Cos(p.T),
		p.R*math.Sin(p.T),
	)
}

func (d Direction) ToPolarVector() PolarVector {
	return NewPolarVector(1, d.t)
}

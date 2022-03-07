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

type Segment struct {
	A, B Coord
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
	return Coord{Vector: v}
}

func (a Vector) Add(b Vector) Vector {
	return NewVector(a.X+b.X, a.Y+b.Y)
}

func (a Vector) Sub(b Vector) Vector {
	return a.Add(b.MulScalar(-1))
}

func (a Vector) MulScalar(r float64) Vector {
	return NewVector(r*a.X, r*a.Y)
}

func (a Vector) Cross(b Vector) float64 {
	return a.X*b.Y - a.Y*b.X
}

func NewCoord(x, y float64) Coord {
	return Coord{Vector: NewVector(x, y)}
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

func NewSegment(a Coord, b Coord) Segment {
	return Segment{A: a, B: b}
}

func (v Coord) AsVector() Vector {
	return v.Vector
}

func NewPolarVector(r, t float64) PolarVector {
	return PolarVector{R: r, T: t}
}

func (p PolarVector) ToVector() Vector {
	return NewVector(
		p.R*math.Cos(p.T),
		p.R*math.Sin(p.T),
	)
}

// a, b, c が反時計回りかどうかを返す。
// 1: 半時計回り
// 0: 直線上
// -1: 時計回り
func CCW(a, b, c Coord) int {
	cross := (b.Sub(a.Vector)).Cross(c.Sub(a.Vector))
	if math.Abs(cross) < 1e-8 {
		return 0
	} else if cross > 0 {
		return 1
	} else {
		return -1
	}
}

func (a Segment) Crosses(b Segment) bool {
	// a を基準にして b の端点が異なる側にあり、かつ
	// b を基準にして a の端点が異なる側にあれば
	// a と b は交差する。
	//
	// どちらかが線分上にある場合は「見えている」と考えることにすると CCW == 0
	// の場合は false を返せばよいことになる。
	return CCW(a.A, a.B, b.A)*CCW(a.A, a.B, b.B) < 0 &&
		CCW(b.A, b.B, a.A)*CCW(b.A, b.B, a.B) < 0
}

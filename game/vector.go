package game

import (
	"math"
)

type Vector struct {
	X, Y float32
}

var (
	ZeroVector = Vector{0, 0}
	OneVector  = Vector{1, 1}
)

// --------------------------------------------------------
// Constructors
// --------------------------------------------------------

func V2(coords ...float32) *Vector {
	l := len(coords)
	if l == 0 {
		return &Vector{X: 0, Y: 0}
	}
	if l == 1 {
		return &Vector{X: coords[0], Y: coords[0]}
	}
	return &Vector{X: coords[0], Y: coords[1]}
}

// Stack-friendly version (no allocation)
func Vec(x, y float32) Vector {
	return Vector{X: x, Y: y}
}

func (v *Vector) Clone() *Vector {
	return &Vector{v.X, v.Y}
}

// --------------------------------------------------------
// Basic mutation operations
// --------------------------------------------------------

func (v *Vector) Set(x, y float32) {
	v.X = x
	v.Y = y
}

func (v *Vector) Add(o *Vector) {
	v.X += o.X
	v.Y += o.Y
}

func (v *Vector) Sub(o *Vector) {
	v.X -= o.X
	v.Y -= o.Y
}

func (v *Vector) Mul(s float32) {
	v.X *= s
	v.Y *= s
}

func (v *Vector) Div(s float32) {
	if s != 0 {
		inv := 1.0 / s
		v.X *= float32(inv)
		v.Y *= float32(inv)
	}
}

// --------------------------------------------------------
// Non-mutating versions
// --------------------------------------------------------

func (v *Vector) Added(o *Vector) *Vector {
	return &Vector{v.X + o.X, v.Y + o.Y}
}

func (v *Vector) Subbed(o *Vector) *Vector {
	return &Vector{v.X - o.X, v.Y - o.Y}
}

func (v *Vector) Muled(s float32) *Vector {
	return &Vector{v.X * s, v.Y * s}
}

func (v *Vector) Divved(s float32) *Vector {
	if s == 0 {
		return &Vector{0, 0}
	}
	inv := 1.0 / s
	return &Vector{v.X * float32(inv), v.Y * float32(inv)}
}

// --------------------------------------------------------
// Vector math
// --------------------------------------------------------

func (v *Vector) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v *Vector) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector) Normalize() {
	l2 := v.X*v.X + v.Y*v.Y
	if l2 > 0 {
		inv := 1 / float32(math.Sqrt(float64(l2)))
		v.X *= inv
		v.Y *= inv
	}
}

func (v *Vector) Normalized() *Vector {
	l2 := v.X*v.X + v.Y*v.Y
	if l2 == 0 {
		return &Vector{0, 0}
	}
	inv := 1 / float32(math.Sqrt(float64(l2)))
	return &Vector{v.X * inv, v.Y * inv}
}

func (v *Vector) Dot(o *Vector) float32 {
	return v.X*o.X + v.Y*o.Y
}

func (v *Vector) Cross(o *Vector) float32 {
	return v.X*o.Y - v.Y*o.X
}

// --------------------------------------------------------
// Rotation
// --------------------------------------------------------

func (v *Vector) Angle(angle float32) {
	sin, cos := math.Sincos(float64(angle))
	s := float32(sin)
	c := float32(cos)

	oldX := v.X
	oldY := v.Y

	v.X = oldX*c - oldY*s
	v.Y = oldX*s + oldY*c
}

func (v *Vector) Angled(angle float32) *Vector {
	sin, cos := math.Sincos(float64(angle))
	s := float32(sin)
	c := float32(cos)

	return &Vector{
		X: v.X*c - v.Y*s,
		Y: v.X*s + v.Y*c,
	}
}

// Degree-based versions (optional)
func (v *Vector) AngleDeg(deg float32) {
	v.Angle(deg * (math.Pi / 180))
}

func (v *Vector) AngledDeg(deg float32) *Vector {
	return v.Angled(deg * (math.Pi / 180))
}

// --------------------------------------------------------
// Lerp / Clamp / Misc helpers
// --------------------------------------------------------

func (v *Vector) Lerp(o *Vector, t float32) {
	v.X = v.X + (o.X-v.X)*t
	v.Y = v.Y + (o.Y-v.Y)*t
}

func (v *Vector) Clamp(max float32) {
	l2 := v.LengthSquared()
	max2 := max * max
	if l2 > max2 {
		scale := max / float32(math.Sqrt(float64(l2)))
		v.X *= scale
		v.Y *= scale
	}
}

func (v *Vector) Clamped(max float32) *Vector {
	l2 := v.LengthSquared()
	max2 := max * max
	if l2 > max2 {
		scale := max / float32(math.Sqrt(float64(l2)))
		return &Vector{v.X * scale, v.Y * scale}
	}
	return &Vector{v.X, v.Y}
}

func (v *Vector) SquareClamp(max float32) {
	v.X = clampf(v.X, -max, max)
	v.Y = clampf(v.Y, -max, max)
}

func (v *Vector) SquaredClamped(max float32) *Vector {
	return &Vector{clampf(v.X, -max, max), clampf(v.Y, -max, max)}
}

func (v *Vector) Dist(o *Vector) float32 {
	dx := v.X - o.X
	dy := v.Y - o.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

func (v *Vector) SqrDist(o *Vector) float32 {
	dx := v.X - o.X
	dy := v.Y - o.Y
	return dx*dx + dy*dy
}

func (v *Vector) Eq(o *Vector) bool {
	return v.X == o.X && v.Y == o.Y
}

func (v *Vector) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

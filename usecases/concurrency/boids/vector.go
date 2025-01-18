package main

import "math"

type Vector struct {
	x, y float64
}

func (v Vector) Add(v2 Vector) Vector {
	return Vector{v.x + v2.x, v.y + v2.y}
}

func (v Vector) Subtract(v2 Vector) Vector {
	return Vector{v.x - v2.x, v.y - v2.y}
}

func (v Vector) Multiply(v2 Vector) Vector {
	return Vector{v.x * v2.x, v.y * v2.y}
}

func (v Vector) MultiplyV(val float64) Vector {
	return Vector{v.x * val, v.y * val}
}

func (v Vector) Divide(v2 Vector) Vector {
	return Vector{v.x / v2.x, v.y / v2.y}
}

func (v Vector) DivideV(val float64) Vector {
	return Vector{v.x / val, v.y / val}
}

func (v Vector) limit(lower, upper float64) Vector {
	return Vector{math.Max(math.Min(v.x, upper), lower), math.Max(math.Min(v.y, upper), lower)}
}

func (v Vector) Distance(v2 Vector) float64 {
	return math.Sqrt(math.Pow(v.x-v2.x, 2) + math.Pow(v.y-v2.y, 2))
}

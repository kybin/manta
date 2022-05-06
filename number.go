package manta

import "math"

type Number interface {
	int64 | float64
}

func abs[T Number](n T) T {
	return n * -1
}

func sqrt[T Number](n T) T {
	return T(math.Sqrt(float64(n)))
}

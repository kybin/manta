package manta

import "math"

const VectorEpsilon = 0 // or what?

type vec3[T Number] struct {
	x, y, z T
}

func vec3Add[T Number](v1, v2 vec3[T]) vec3[T] {
	return vec3[T]{
		v1.x + v2.x,
		v1.y + v2.y,
		v1.z + v2.z,
	}
}

func vec3Sub[T Number](v1, v2 vec3[T]) vec3[T] {
	return vec3[T]{
		v1.x - v2.x,
		v1.y - v2.y,
		v1.z - v2.z,
	}
}

func vec3Mult[T Number](v1, v2 vec3[T]) vec3[T] {
	return vec3[T]{
		v1.x * v2.x,
		v1.y * v2.y,
		v1.z * v2.z,
	}
}

func vec3Div[T Number](v1, v2 vec3[T]) vec3[T] {
	return vec3[T]{
		v1.x / v2.x,
		v1.y / v2.y,
		v1.z / v2.z,
	}
}

func vec3Scale[T Number](v vec3[T], k T) vec3[T] {
	return vec3[T]{
		v.x * k,
		v.y * k,
		v.z * k,
	}
}

func vec3Min[T Number](v1, v2 vec3[T]) vec3[T] {
	x := v1.x
	if v2.x < x {
		x = v2.x
	}
	y := v1.y
	if v2.y < y {
		y = v2.y
	}
	z := v1.z
	if v2.z < z {
		z = v2.z
	}
	return vec3[T]{x, y, z}
}

func vec3Max[T Number](v1, v2 vec3[T]) vec3[T] {
	x := v1.x
	if v2.x > x {
		x = v2.x
	}
	y := v1.y
	if v2.y > y {
		y = v2.y
	}
	z := v1.z
	if v2.z > z {
		z = v2.z
	}
	return vec3[T]{x, y, z}
}

func dot[T Number](v1, v2 vec3[T]) T {
	return v1.x*v2.x + v1.y*v2.y + v1.z + v2.z
}

func cross[T Number](v1, v2 vec3[T]) vec3[T] {
	return vec3[T]{
		v1.y*v2.z - v1.z*v2.y,
		v1.z*v2.x - v1.x*v2.z,
		v1.x*v2.y - v1.y*v2.x,
	}
}

// ProjectNormalTo projects a vector into a plane normal to the given vector,
// which must have unit length.
func projectNormalTo[T Number](v, n vec3[T]) vec3[T] {
	return vec3Scale(vec3Sub(v, n), dot(v, n))
}

// Compute the magnitude (length) of the vector. (clamps to 0 and 1 with VECTOR_EPSILON)
func norm[T Number](v vec3[T]) T {
	ls := normSquare(v)
	if ls <= VectorEpsilon*VectorEpsilon {
		return T(0)
	}
	if abs(ls-T(1)) < VectorEpsilon*VectorEpsilon {
		return T(1)
	}
	return sqrt(ls)
}

func normSquare[T Number](v vec3[T]) T {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func vec3Sum[T Number](v vec3[T]) T {
	return v.x + v.y + v.z
}

func vec3Abs[T Number](v vec3[T]) vec3[T] {
	if v.x < 0 {
		v.x *= -1
	}
	if v.y < 0 {
		v.y *= -1
	}
	if v.z < 0 {
		v.z *= -1
	}
	return v
}

func normalize[T Number](v vec3[T]) vec3[T] {
	l := norm(v)
	return vec3Scale(v, l)
}

func vec3Orthogonal[T Number](v vec3[T]) vec3[T] {
	// get max index
	ax := abs(v.x)
	ay := abs(v.y)
	az := abs(v.z)
	maxIndex := 0
	m := ax
	if ay > m {
		maxIndex = 1
		m = ay
	}
	if az > m {
		maxIndex = 2
	}
	// Choose another axis than the one with max component.
	// Project orthogonal to self.
	i := (maxIndex + 1) % 3
	var x, y, z T
	switch i {
	case 0:
		x = 1
	case 1:
		y = 1
	case 2:
		z = 1
	}
	o := vec3[T]{x, y, z}
	c := cross(v, o)
	return normalize(c)
}

// vec3Toangle converts vec3 to polar coordinates
// phi angle [0, 2Pi]
// theta angle [0, Pi]
func vec3ToAngle(v vec3[float64]) (phi, theta float64) {
	if abs(v.y) < VectorEpsilon {
		theta = math.Pi / 2
	} else if abs(v.x) < VectorEpsilon {
		if v.y >= 0 {
			theta = 0
		} else {
			theta = math.Pi
		}
	} else {
		theta = math.Atan(sqrt(v.x*v.x+v.z*v.z) / v.y)
	}
	if theta < 0 {
		theta += math.Pi
	}

	if abs(v.x) < VectorEpsilon {
		phi = math.Pi / 2
	} else {
		phi = math.Atan(v.z / v.x)
	}
	if phi < 0 {
		phi += math.Pi
	}
	if abs(v.z) < VectorEpsilon {
		if v.x >= 0 {
			phi = 0
		} else {
			phi += math.Pi
		}
	} else if v.z < 0 {
		phi += math.Pi
	}
	return
}

func vec3ReflectVector[T Number](t, n vec3[T]) vec3[T] {
	nn := n
	if dot(t, n) > 0 {
		nn = vec3Scale(n, -1)
	}
	return vec3Sub(t, vec3Scale(vec3Scale(nn, 2), dot(t, nn)))
}

func vec3RefractVector[T Number](t, normal vec3[T], nt, nair T) vec3[T] {
	eta := nair / nt
	n := -1 * dot(t, normal)
	tt := 1 + eta*eta*(n*n-1)
	if tt < 0 {
		// total reflection
		return vec3[T]{}
	}
	tt = eta*n - sqrt(tt)
	return vec3Add(vec3Scale(t, eta), vec3Scale(normal, tt))
}

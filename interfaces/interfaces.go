package interfaces

import (
	"math"

	"github.com/shoenig/test/internal/constraints"
)

// EqualsFunc represents a type implementing the Equals method.
type EqualsFunc[A any] interface {
	Equals(A) bool
}

// LessFunc represents any type implementing the Less method.
type LessFunc[A any] interface {
	Less(A) bool
}

// Map represents any map type where keys are comparable.
type Map[K comparable, V any] interface {
	~map[K]V
}

// MapEqualsFunc represents any map type where keys are comparable and values implement .Equals method.
type MapEqualsFunc[K comparable, V EqualsFunc[V]] interface {
	~map[K]V
}

// Number is float, integer, or complex.
type Number interface {
	constraints.Ordered
	constraints.Float | constraints.Integer | constraints.Complex
}

// Numeric returns false if n is Inf/NaN.
//
// Always returns true for integral values.
func Numeric[N Number](n N) bool {
	check := func(f float64) bool {
		if math.IsNaN(f) {
			return false
		} else if math.IsInf(f, 0) {
			return false
		}
		return true
	}
	return check(float64(n))
}

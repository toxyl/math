package math

// Number is a constraint that covers most common numeric types.
type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

// Float is a constraint for floating-point types.
type Float interface {
	float32 | float64
}

// Matrix represents a two-dimensional slice of float64 values.
type Matrix [][]float64

package math

import "fmt"

// Add returns the sum of x and y.
func Add[N Number](x, y N) N {
	return x + y
}

// Sub returns the difference between x and y.
func Sub[N Number](x, y N) N {
	return x - y
}

// Blend blends two numbers based on the proportion p.
func Blend[N Number](x, y N, p float64) N {
	return N(float64(x)*(1-p) + float64(y)*p)
}

// Lerp performs linear interpolation between a and b with weight t.
// The parameter t is typically in range [0,1].
func Lerp[N Number](a, b N, t float64) N {
	return N(float64(a) + t*(float64(b)-float64(a)))
}

// Avg calculates the average of a variadic number of values.
func Avg[N Number](x ...N) N {
	res := 0.0
	for i, n := range x {
		if i == 0 {
			res = float64(n)
			continue
		}
		res = (res + float64(n)) / 2.0
	}
	return N(res)
}

// Clamp restricts x to be within the range [min, max].
func Clamp[N Number](x, min, max N) N {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

// Delta calculates the difference between x and y.
func Delta[N Number](x, y N) N {
	maxVal := Max(x, y)
	minVal := Min(x, y)
	if maxVal >= 0 && minVal <= 0 {
		return N(Abs(minVal) + maxVal)
	}
	return Abs(Sub(maxVal, minVal))
}

// Wrap wraps the value within the interval [min, max).
func Wrap[N Number](val, min, max N) N {
	if min > max {
		min, max = max, min
	}
	rangeSize := Sub(max, min)
	val = Mod(Sub(val, min), rangeSize)
	if val < 0 {
		val = Add(val, rangeSize)
	}
	return Add(val, min)
}

// Normalize converts a value from a given range [min, max] to a normalized value in [0,1].
// If min equals max, it returns 0.
func Normalize[N Number](val, min, max N) N {
	if max == min {
		return 0
	}
	return (val - min) / (max - min)
}

// Denormalize converts a normalized value in [0,1] to the target range [min, max].
// The normalized value is clamped between 0 and 1.
func Denormalize[N Number](norm, min, max N) N {
	if norm < 0 {
		norm = 0
	} else if norm > 1 {
		norm = 1
	}
	return norm*(max-min) + min
}

// FormatNumber formats a number with appropriate suffixes for large values.
func FormatNumber[N Number](number N, decimals int) string {
	suffix := " "
	divisor := 1.0
	n := float64(number)
	switch {
	case n >= 1_000_000_000_000 || n <= -1_000_000_000_000:
		suffix = "t"
		divisor = 1_000_000_000_000
	case n >= 1_000_000_000 || n <= -1_000_000_000:
		suffix = "b"
		divisor = 1_000_000_000
	case n >= 1_000_000 || n <= -1_000_000:
		suffix = "m"
		divisor = 1_000_000
	case n >= 1_000 || n <= -1_000:
		suffix = "k"
		divisor = 1_000
	}
	return fmt.Sprintf("%*.*f", 5+decimals, decimals, n/divisor) + suffix
}

// MaxLenStr returns the length of the longest string from the provided arguments.
func MaxLenStr(strs ...string) int {
	l := 0
	for _, s := range strs {
		l = Max(l, len(s))
	}
	return l
}

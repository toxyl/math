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

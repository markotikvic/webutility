package numbersutility

import (
	"fmt"
)

// ClampInt64 ...
func ClampInt64(v, min, max int64) int64 {
	if v < min {
		return min
	} else if v > max {
		return max
	}

	return v
}

// InRangeInt64 ...
func InRangeInt64(v, min, max int64) bool {
	return (v >= min && v <= max)
}

// Int64ToString ...
func Int64ToString(i int64) string {
	return fmt.Sprintf("%d", i)
}

// Int64PtrToString ...
func Int64PtrToString(i *int64) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%d", *i)
}

// BoolToInt64 ...
func BoolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

// Int64ToBool ...
func Int64ToBool(i int64) bool {
	return i != 0
}

func MaxInt(vars ...int) (max int) {
	max = vars[0]
	for _, v := range vars {
		if v > max {
			max = v
		}
	}
	return max
}

func MinInt64(vars ...int64) (min int64) {
	min = vars[0]
	for _, v := range vars {
		if v < min {
			min = v
		}
	}
	return min
}

func FormatInt64(i int64) string {
	res := ""
	for i >= 1000 {
		rem := i % 1000
		i = i / 1000
		//res = res + fmt.Sprintf(".%03d", rem)
		res = fmt.Sprintf(".%03d", rem) + res
	}
	res = fmt.Sprintf("%d", i) + res
	return res
}

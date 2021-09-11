package stringutility

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const sanitisationPatern = "\"';&*<>=\\`:"

// SanitiseString removes characters from s found in patern and returns new modified string.
func SanitiseString(s string) string {
	return ReplaceAny(s, sanitisationPatern, "")
}

// IsWrappedWith ...
func IsWrappedWith(src, begin, end string) bool {
	return strings.HasPrefix(src, begin) && strings.HasSuffix(src, end)
}

// ParseInt64Arr ...
func ParseInt64Arr(s, sep string) (arr []int64) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	parts := strings.Split(s, sep)
	arr = make([]int64, len(parts))
	for i, p := range parts {
		num := StringToInt64(p)
		arr[i] = num
	}

	return arr
}

// Int64SliceToString ...
func Int64SliceToString(arr []int64) (s string) {
	if len(arr) == 0 {
		return ""
	}

	s += fmt.Sprintf("%d", arr[0])
	for i := 1; i < len(arr); i++ {
		s += fmt.Sprintf(",%d", arr[i])
	}

	return s
}

// CombineStrings ...
func CombineStrings(s1, s2, glue string) string {
	s1 = strings.TrimSpace(s1)
	s2 = strings.TrimSpace(s2)

	if s1 != "" && s2 != "" {
		s1 += glue + s2
	} else {
		s1 += s2
	}

	return s1
}

// ReplaceAny replaces any of the characters from patern found in s with r and returns a new resulting string.
func ReplaceAny(s, patern, r string) (n string) {
	n = s
	for _, c := range patern {
		n = strings.Replace(n, string(c), r, -1)
	}
	return n
}

// StringToBool ...
func StringToBool(s string) bool {
	res, _ := strconv.ParseBool(s)
	return res
}

// BoolToString ...
func BoolToString(b bool) string {
	return fmt.Sprintf("%b", b)
}

// StringSliceContains ...
func StringSliceContains(slice []string, s string) bool {
	for i := range slice {
		if slice[i] == s {
			return true
		}
	}
	return false
}

func SplitString(s, sep string) (res []string) {
	parts := strings.Split(s, sep)
	for _, p := range parts {
		if p != "" {
			res = append(res, p)
		}
	}
	return res
}

// StringAt ...
func StringAt(s string, index int) string {
	if len(s)-1 < index || index < 0 {
		return ""
	}

	return string(s[index])
}

func StringAtRune(s string, index int) string {
	str := []rune(s)
	if index < len(str) {
		return string(str[index])
	}
	return ""
}

// SplitText ...
func SplitText(s string, maxLen int) (lines []string) {
	runes := []rune(s)

	i, start, sep, l := 0, 0, 0, 0
	for i = 0; i < len(runes); i++ {
		c := runes[i]

		if unicode.IsSpace(c) {
			sep = i
		}

		if c == '\n' {
			if start != sep {
				lines = append(lines, string(runes[start:sep]))
			}
			start = i
			sep = i
			l = 0
		} else if l >= maxLen {
			if start != sep {
				lines = append(lines, string(runes[start:sep]))
				sep = i
				start = i - 1
				l = 0
			}
		} else {
			l++
		}
	}
	if start != i-1 {
		lines = append(lines, string(runes[start:i-1]))
	}

	return lines
}

func CutTextWith(txt string, maxLen int, tail string) string {
	if len(txt) < maxLen || len(txt) <= len(tail) {
		return txt
	}

	return txt[:maxLen-3] + tail
}

func LimitTextWith(txt string, maxLen int, tail string) string {
	if len(txt) <= maxLen {
		return txt
	}

	return txt[:maxLen] + tail
}

// SplitStringAtWholeWords ...
func SplitStringAtWholeWords(s string, maxLen int) (res []string) {
	parts := strings.Split(s, " ")

	res = append(res, parts[0])
	i := 0
	for j := 1; j < len(parts); j++ {
		p := strings.TrimSpace(parts[j])
		if len(p) > maxLen {
			// TODO(marko): check if maxLen is >= 3
			p = p[0 : maxLen-3]
			p += "..."
		}
		if len(res[i])+len(p)+1 <= maxLen {
			res[i] += " " + p
		} else {
			res = append(res, p)
			i++
		}
	}

	return res
}

// StringToInt64 ...
func StringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// StringToFloat64 ...
func StringToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func StringToValidInt64(s string) (int64, bool) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return i, false
	}
	return i, true
}

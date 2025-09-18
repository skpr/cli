package slice

import "strings"

// Contains a string within a slice.
func Contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}

	return false
}

// Equal tells whether a and b contain the same elements.
// This helper function allows an equal check which is dependant of the keys.
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// If "b" contains all the same values as "a".
	for _, v := range a {
		if !Contains(b, v) {
			return false
		}
	}

	// If "a" contains all the same values as "b".
	for _, v := range b {
		if !Contains(a, v) {
			return false
		}
	}

	return true
}

// Remove a string from a slice.
func Remove(slice []string, s string) []string {
	var result []string

	for _, item := range slice {
		if item == s {
			continue
		}

		result = append(result, item)
	}

	return result
}

// AppendSlice to an existing slice without adding duplicates.
func AppendSlice(slice, extra []string) []string {
	for _, e := range extra {
		slice = AppendIfMissing(slice, e)
	}

	return slice
}

// AppendIfMissing to an existing slice.
func AppendIfMissing(slice []string, i string) []string {
	if Contains(slice, i) {
		return slice
	}

	return append(slice, i)
}

// ToMap converts a slice of strings into a map using the provided delimiter.
func ToMap(slice []string, delimiter string) map[string]string {
	m := make(map[string]string, len(slice))

	for _, s := range slice {
		sl := strings.Split(s, delimiter)

		if len(sl) != 2 {
			continue
		}

		m[sl[0]] = sl[1]
	}

	return m
}

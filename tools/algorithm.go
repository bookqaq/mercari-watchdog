package tools

import "reflect"

// Delete item in ordered array src that filter when reflect.DeepEqual(item, value) != true, return lenght deleted
func DeleteInvalidItem[T any](src []T, value T) int {
	deleted, formerpt, length := 0, 0, len(src)
	for i := 0; i < length; i++ {
		if reflect.DeepEqual(src[i], value) {
			deleted++
		} else {
			src[formerpt] = src[i]
			formerpt++
		}
	}
	return deleted
}

// Find longest common item in slice, return common item's length
func LongestCommon[T comparable](arr1, arr2 []T) int {
	m := make(map[T]int, len(arr1))
	ans := 0

	for _, v := range arr1 {
		m[v]++
	}

	for _, v := range arr2 {
		m[v]--
	}

	for _, v := range m {
		if v == 0 {
			ans++
		}
	}

	return ans
}

// Replace rune to new in s if rune in old
// TODO: Change old to map[rune]struct{} for a faster find speed
func StringMultipleReplacer(s string, old []rune, new rune) string {
	r := []rune(s)
	for i, v := range r {
		for _, item := range old {
			if v == item {
				r[i] = ' '
				break
			}
		}
	}
	return string(r)
}

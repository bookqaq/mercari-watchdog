package tools

import (
	"fmt"
	"reflect"
	"strings"
)

func ConcatKeyword(keywords []string) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, keywords[0])
	for _, item := range keywords[1:] {
		fmt.Fprintf(&builder, " %s", item)
	}
	return builder.String()
}

// Delete item in ordered array src that return true in reflect.DeepEqual(item, value), return lenght deleted
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

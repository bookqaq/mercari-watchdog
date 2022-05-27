package utils

import (
	"fmt"
	"strings"
)

type PushMsg struct {
	Dst int64
	S   []string
}

func ConcatKeyword(keywords []string) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, keywords[0])
	for _, item := range keywords[1:] {
		fmt.Fprintf(&builder, " %s", item)
	}
	return builder.String()
}

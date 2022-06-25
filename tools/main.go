package tools

type PushMsg struct {
	Dst int64    // target group number
	S   []string // messages, item in array will be sent separately
}

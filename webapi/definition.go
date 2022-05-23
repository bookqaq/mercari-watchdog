package webapi

type genericPostReply struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
	Auth    string      `json:"auth"`
}

type taskAddSubmitData struct {
	Auth string `json:"auth"`
	Data struct {
		KeywordsOrig string   `json:"keywords_orig"`
		Keywords     []string `json:"keywords"`
		MustMatch    []string `json:"mustMatch"`
		Owner        int64    `json:"owner"`
		Group        int64    `json:"group"`
		Interval     int      `json:"interval"`
		TargetPrice  [2]int   `json:"targetPrice"`
		MaxPage      int      `json:"maxPage"`
	} `json:"data"`
}

package utils

type AnalysisTask struct {
	TaskID      int32    `bson:"taskID"`      // a unique id for every search.
	Owner       int32    `bson:"owner"`       // a QQ id to the person owning the task
	Keywords    []string `bson:"keywords"`    // an array about search keywords, Keywords.0 is used to send request to mercari
	Interval    int      `bson:"interval"`    // seconds between task run repeatedly
	TargetPrice int      `bson:"targetPrice"` // Optional, items that prices lower than this will be selected into result. default(0 or negative number) means any price is accepted.
	MaxPage     int      `bson:"maxPage"`     // I perfer 30 items in one page(same as mercari default), so only page count is provided to control total items.
	Sort        string   `bson:"sort"`
	Order       string   `bson:"order"`
}

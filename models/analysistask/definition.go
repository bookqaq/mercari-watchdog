package analysistask

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalysisTask struct {
	ID          primitive.ObjectID `bson:"_id"`
	TaskID      int32              `bson:"taskID"       json:"taskID"`   // a unique id for every search.
	Owner       int64              `bson:"owner"        json:"owner"`    // a QQ id to the person owning the task
	Group       int64              `bson:"group"        json:"group"`    // owner's group id
	Keywords    []string           `bson:"keywords"     json:"keywords"` // an array about search keywords, Keywords.0 is used to send request to mercari
	MustMatch   []string           `bson:"mustMatch"    json:"mustMatch"`
	Interval    int                `bson:"interval"     json:"interval"`    // seconds between task run repeatedly
	TargetPrice [2]int             `bson:"targetPrice"  json:"targetPrice"` // items that prices lower than this will be selected into result. default(0 or negative number) means any price is accepted.
	MaxPage     int                `bson:"maxPage"      json:"maxPage"`     // I perfer 30 items in one page(same as mercari default), so only page count is provided to control total items.
	Sort        string             `bson:"sort"`
	Order       string             `bson:"order"`
}

func (t *AnalysisTask) FormatSimplifiedChinese() string {
	return fmt.Sprintf("任务ID:%v\n号主QQ:%v\n关键词:%s\n目标价格:%v~%v\n最大页数:%v 搜索间隔:%v",
		t.TaskID, t.Owner, concatKeyword(t.Keywords), t.TargetPrice[0], t.TargetPrice[1], t.MaxPage, t.Interval)
}

func concatKeyword(keywords []string) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, keywords[0])
	for _, item := range keywords[1:] {
		fmt.Fprintf(&builder, " %s", item)
	}
	return builder.String()
}

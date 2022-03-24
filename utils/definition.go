package utils

import (
	"fmt"

	"github.com/bookqaq/goForMercari/mercarigo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalysisTask struct {
	ID          primitive.ObjectID `bson:"_id"`
	TaskID      int32              `bson:"taskID"`      // a unique id for every search.
	Owner       int64              `bson:"owner"`       // a QQ id to the person owning the task
	Group       int64              `bson:"group"`       // owner's group id
	Keywords    []string           `bson:"keywords"`    // an array about search keywords, Keywords.0 is used to send request to mercari
	Interval    int                `bson:"interval"`    // seconds between task run repeatedly
	TargetPrice [2]int             `bson:"targetPrice"` // Optional, items that prices lower than this will be selected into result. default(0 or negative number) means any price is accepted.
	MaxPage     int                `bson:"maxPage"`     // I perfer 30 items in one page(same as mercari default), so only page count is provided to control total items.
	Sort        string             `bson:"sort"`
	Order       string             `bson:"order"`
}

type AnalysisData struct {
	ID       primitive.ObjectID      `bson:"_id"`
	Keywords []string                `bson:"keyword"` // a primary-key-like value
	TaskID   int32                   `bson:"taskID"`
	Time     int64                   `bson:"time"`   // unix time
	Length   int                     `bson:"length"` // data amount in total
	Data     []mercarigo.MercariItem `bson:"data"`
}

type GroupWhitelist struct {
	ID    primitive.ObjectID `bson:"_id"`
	Group int64              `bson:"group"`
}

type PushMsg struct {
	Dst int64
	S   []string
}

func ConcatKeyword(keywords []string) string {
	kwstring := keywords[0]
	for _, item := range keywords[1:] {
		kwstring += " "
		kwstring += item
	}
	return kwstring
}

func (t *AnalysisTask) FormatSimplifiedChinese() string {
	return fmt.Sprintf("任务ID:%v\n号主QQ:%v\n关键词:%s\n目标价格:%v~%v\n最大页数:%v 搜索间隔:%v",
		t.TaskID, t.Owner, ConcatKeyword(t.Keywords), t.TargetPrice[0], t.TargetPrice[1], t.MaxPage, t.Interval)
}

func (d *AnalysisData) FormatSimplifiedChinese() []string {
	res := make([]string, 1)
	res[0] = fmt.Sprintf("任务ID:%v\n关键词:%s\n时间:%v\n蹲到符合要求的结果数为%v", d.TaskID, ConcatKeyword(d.Keywords), d.Time, d.Length)

	if d.Length > 0 {
		for _, item := range d.Data {
			tmp := fmt.Sprintf("[CQ:image,file=%s]名称:%s\n价格:%vyen\n状态:%s\n更新时间:%v\n链接:%s",
				item.ImageURL[0], item.ProductName, item.Price, item.Status, item.Updated, item.GetProductURL())
			res = append(res, tmp)
		}
	}
	return res
}

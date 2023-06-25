package analysisdata

import (
	"fmt"
	"strings"
	"time"

	"bookq.xyz/mercari-watchdog/tools"
	wrapperv1 "github.com/bookqaq/mer-wrapper/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalysisData struct {
	ID       primitive.ObjectID      `bson:"_id"`
	Keywords []string                `bson:"keyword"`
	TaskID   int32                   `bson:"taskID"` // a primary-key-alike value
	Time     int64                   `bson:"time"`   // unix time
	Length   int                     `bson:"length"` // length of data
	Data     []wrapperv1.MercariItem `bson:"data"`
}

// privleged one display items in mulitple messages
func (d *AnalysisData) PrivlegedFormatSimplifiedChinese() []string {
	location, _ := time.LoadLocation("Asia/Shanghai")
	res := make([]string, 1, 6)
	res[0] = fmt.Sprintf("任务ID:%v\n关键词:%s\n时间:%s\n蹲到符合要求的结果数为%v",
		d.TaskID, tools.ConcatKeyword(d.Keywords), time.Unix(d.Time, 0).In(location).Format("2006-01-02 15:04:05"), d.Length)

	if d.Length > 0 {
		for _, item := range d.Data {
			tmp := fmt.Sprintf("[CQ:image,file=%s]名称:%s\n价格:%vyen\n更新时间:%s\n链接:%s",
				item.ImageURL[0], item.ProductName, item.Price,
				time.Unix(item.Updated, 0).In(location).Format("2006-01-02 15:04:05"), item.ProductId)
			res = append(res, tmp)
		}
	}
	return res
}

// normal one display items in one message
func (d *AnalysisData) FormatSimplifiedChinese() string {
	location, _ := time.LoadLocation("Asia/Shanghai")
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("任务ID:%v\n关键词:%s\n时间:%s\n",
		d.TaskID, tools.ConcatKeyword(d.Keywords), time.Unix(d.Time, 0).In(location).Format("2006-01-02 15:04:05")))

	if d.Length > 0 {
		for _, item := range d.Data {
			builder.WriteString(fmt.Sprintf("\n[CQ:image,file=%s]\n名称:%s\n价格:%vyen\n链接:%s",
				item.ImageURL[0], item.ProductName, item.Price, item.ProductId))
		}
	}
	return builder.String()
}

// The layout string used by the Parse function and Format method
// shows by example how the reference time should be represented.
// We stress that one must show how the reference time is formatted,
// not a time of the user's choosing. Thus each layout string is a
// representation of the time stamp,
//	Jan 2 15:04:05 2006 MST
// An easy way to remember this value is that it holds, when presented
// in this order, the values (lined up with the elements above):
//	  1 2  3  4  5    6  -7

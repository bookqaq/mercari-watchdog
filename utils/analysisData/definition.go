package analysisdata

import (
	"fmt"
	"time"

	"bookq.xyz/mercari-watchdog/utils/tools"
	"github.com/bookqaq/goForMercari/mercarigo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalysisData struct {
	ID       primitive.ObjectID      `bson:"_id"`
	Keywords []string                `bson:"keyword"`
	TaskID   int32                   `bson:"taskID"` // a primary-key-alike value
	Time     int64                   `bson:"time"`   // unix time
	Length   int                     `bson:"length"` // data amount in total
	Data     []mercarigo.MercariItem `bson:"data"`
}

func (d *AnalysisData) FormatSimplifiedChinese() []string {
	location, _ := time.LoadLocation("Asia/Shanghai")
	res := make([]string, 1, 6)
	res[0] = fmt.Sprintf("任务ID:%v\n关键词:%s\n时间:%v\n蹲到符合要求的结果数为%v", d.TaskID, tools.ConcatKeyword(d.Keywords), d.Time, d.Length)

	if d.Length > 0 {
		for _, item := range d.Data {
			tmp := fmt.Sprintf("[CQ:image,file=%s]名称:%s\n价格:%vyen\n状态:%s\n更新时间:%v\n链接:%s",
				item.ImageURL[0], item.ProductName, item.Price, item.Status,
				time.Unix(item.Updated, 0).In(location).Format("2022-05-23 12:46:09"), item.GetProductURL())
			res = append(res, tmp)
		}
	}
	return res
}

package utils

import (
	"github.com/bookqaq/goForMercari/mercarigo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalysisTask struct {
	ID          primitive.ObjectID `bson:"_id"`
	TaskID      int32              `bson:"taskID"`      // a unique id for every search.
	Owner       string             `bson:"owner"`       // a QQ id to the person owning the task
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

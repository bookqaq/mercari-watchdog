package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"bookq.xyz/mercariWatchdog/tasks"
	"bookq.xyz/mercariWatchdog/utils"
	"github.com/bookqaq/goForMercari/mercarigo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	//debug_addAnalysisData()
	tasks.Boot()
}

func debug_addAnalysisData() {
	tmp, err := mercarigo.Mercari_search("sasakure", "created_time", "desc", "", 30, 3)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	data := utils.AnalysisData{}
	data.Data = make([]mercarigo.MercariItem, 0)
	for _, item := range tmp {
		if strings.Contains(item.ProductName, "lasah") {
			data.Data = append(data.Data, item)
		}
	}

	data.ID = primitive.NewObjectID()
	data.Keywords = make([]string, 2)
	data.Keywords[0], data.Keywords[1] = "sasakure", "lasah"
	data.TaskID = 213411
	data.Time = time.Now().Unix()
	data.Length = len(data.Data)
	coll := utils.Connect().Client().Database("mercariWatchdogDatabase").Collection("AnalysisData")
	res, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v", res.InsertedID)
}

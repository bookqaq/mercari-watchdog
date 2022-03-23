package utils

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bep/debounce"
	"github.com/bookqaq/goForMercari/mercarigo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var lock sync.Mutex

var db *mongo.Database
var AddTaskChannel chan AnalysisTask
var tasksToAdd []AnalysisTask

func Init() {
	db = Connect()
	AddTaskChannel = make(chan AnalysisTask, 10)
	tasksToAdd = make([]AnalysisTask, 0)
	go addTaskBuffer()
}

// DB oriented:

func GetAllTasks(interval int) ([]AnalysisTask, error) {
	coll := db.Collection("AnalysisTask")
	cursor, err := coll.Find(context.TODO(), bson.D{primitive.E{Key: "interval", Value: interval}})
	if err != nil {
		return nil, err
	}
	var result []AnalysisTask
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetTasksByQQ(qq int64) ([]AnalysisTask, error) {
	var result []AnalysisTask
	coll := db.Collection("AnalysisTask")
	cursor, err := coll.Find(context.TODO(), bson.D{primitive.E{Key: "owner", Value: qq}})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteTask(taskID int32) error {
	coll := db.Collection("AnalysisTask")
	res, err := coll.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "taskID", Value: taskID}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("未找到可删除的任务")
	}
	return nil
}

func GetDataDB(taskid int32) (AnalysisData, error) {
	coll := db.Collection("AnalysisData")
	var result AnalysisData
	err := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "taskID", Value: taskid}}).Decode(&result)
	if err != nil {
		return AnalysisData{}, err
	}
	return result, nil
}

func UpdateDataDB(data AnalysisData) error {
	coll := db.Collection("AnalysisData")
	res := coll.FindOneAndReplace(context.TODO(), bson.D{primitive.E{Key: "_id", Value: data.ID}}, data)
	if err := res.Err(); err != nil {
		return err
	}
	return nil
}

func InsertDataDB(task AnalysisTask) error {
	result, err := mercarigo.Mercari_search(task.Keywords[0], task.Sort, task.Order, "", 30, 3)
	if err != nil {
		return err
	}
	result = KeywordFilter(task, result)
	result = PriceFilter(task, result)

	adata := AnalysisData{}
	adata.ID = primitive.NewObjectID()
	adata.TaskID = task.TaskID
	adata.Data = append(adata.Data, result...)
	adata.Length = len(result)
	adata.Time = time.Now().Unix()
	adata.Keywords = append(adata.Keywords, task.Keywords...)

	coll := db.Collection("AnalysisData")
	_, err = coll.InsertOne(context.TODO(), adata)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDataDB(taskID int32) error {
	coll := db.Collection("AnalysisData")
	res, err := coll.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "taskID", Value: taskID}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("未找到可删除的历史数据")
	}
	return nil
}

func FindWhitelist(group int64) (bool, error) {
	coll := db.Collection("GroupWhitelist")
	res := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "group", Value: group}})
	if res.Err() != nil {
		return false, res.Err()
	}
	return true, nil
}

// filters

func KeywordFilter(task AnalysisTask, data []mercarigo.MercariItem) []mercarigo.MercariItem {
	keywords := task.Keywords[1:]
	for _, keyword := range keywords {
		tmp := make([]mercarigo.MercariItem, 0)
		for _, item := range data {
			if strings.Contains(item.ProductName, keyword) {
				tmp = append(tmp, item)
			}
			data = tmp
		}
	}
	return data
}

func PriceFilter(task AnalysisTask, data []mercarigo.MercariItem) []mercarigo.MercariItem {
	result := make([]mercarigo.MercariItem, 0)
	if task.TargetPrice[0] >= 0 && task.TargetPrice[1] >= task.TargetPrice[0] {
		for _, item := range data {
			if item.Price >= task.TargetPrice[0] && item.Price <= task.TargetPrice[1] {
				result = append(result, item)
			}
		}
	} else {
		return data
	}
	return result
}

// Cache and insert tasks.

func addTaskBuffer() {
	debounced := debounce.New(2 * time.Second)
	for {
		newtask := <-AddTaskChannel
		debounced(addTasks)
		tasksToAdd = append(tasksToAdd, newtask)
	}
}

func addTasks() {
	lock.Lock()
	data := make([]interface{}, len(tasksToAdd))
	for i, task := range tasksToAdd {
		task.ID = primitive.NewObjectID()
		data[i] = task
	}
	tasksToAdd = make([]AnalysisTask, 0)
	lock.Unlock()
	coll := db.Collection("AnalysisTask")
	_, err := coll.InsertMany(context.TODO(), data)
	if err != nil {
		fmt.Printf("err inserting tasks from buffer, %s\n", err)
	}
}

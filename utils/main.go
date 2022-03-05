package utils

import (
	"context"
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

// keyword filter

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
		data[i] = task
	}
	tasksToAdd = make([]AnalysisTask, 0)
	coll := db.Collection("AnalysisTask")
	coll.InsertMany(context.TODO(), data)

	lock.Unlock()
}

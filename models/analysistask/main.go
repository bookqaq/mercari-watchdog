package analysistask

import (
	"context"
	"fmt"
	"sync"
	"time"

	"bookq.xyz/mercari-watchdog/database"
	"github.com/bep/debounce"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var lock sync.Mutex
var AddTaskChannel = make(chan AnalysisTask, 10)
var tasksToAdd = make([]AnalysisTask, 0, 10)

func IfExist(taskID int32) bool {
	coll := database.DB.Collection("AnalysisTask")
	if err := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "taskID", Value: taskID}}).Err(); err == mongo.ErrNoDocuments {
		return false
	}
	return true
}

func GetAll(interval int) ([]AnalysisTask, error) {
	coll := database.DB.Collection("AnalysisTask")
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

// Return all AnalysisTasks owned by user
func GetByQQ(qq int64, group int64) ([]AnalysisTask, error) {
	var result []AnalysisTask
	coll := database.DB.Collection("AnalysisTask")
	cursor, err := coll.Find(context.TODO(), bson.D{primitive.E{Key: "owner", Value: qq}, primitive.E{Key: "group", Value: group}})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Delete(taskID int32, qq int64) error {
	coll := database.DB.Collection("AnalysisTask")
	res, err := coll.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "taskID", Value: taskID}, primitive.E{Key: "owner", Value: qq}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("未找到可删除的任务")
	}
	return nil
}

// Cache and insert.

// implement debounce when inserting tasks
func AddTaskBuffer() {
	debounced := debounce.New(2 * time.Second)
	for {
		newtask := <-AddTaskChannel
		debounced(addTasks)
		tasksToAdd = append(tasksToAdd, newtask)
	}
}

// add multiple tasks in buffer by coll.InsertMany
func addTasks() {
	lock.Lock()
	data := make([]interface{}, len(tasksToAdd))
	for i, task := range tasksToAdd {
		task.ID = primitive.NewObjectID()
		data[i] = task
	}
	tasksToAdd = make([]AnalysisTask, 0)
	lock.Unlock()
	coll := database.DB.Collection("AnalysisTask")
	_, err := coll.InsertMany(context.TODO(), data)
	if err != nil {
		fmt.Printf("err inserting tasks from buffer, %s\n", err)
	}
}

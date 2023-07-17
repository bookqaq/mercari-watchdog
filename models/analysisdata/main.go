package analysisdata

import (
	"context"
	"fmt"
	"log"
	"time"

	"bookq.xyz/mercari-watchdog/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var optionsDoUpsert = options.FindOneAndReplace().SetUpsert(true)

func GetOne(taskid int32) (AnalysisData, error) {
	coll := database.DB.Collection("AnalysisData")
	var result AnalysisData
	err := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "taskID", Value: taskid}}).Decode(&result)
	if err != nil {
		return AnalysisData{}, err
	}
	return result, nil
}

func Update(data AnalysisData) error {
	coll := database.DB.Collection("AnalysisData")
	res := coll.FindOneAndReplace(context.TODO(), bson.D{primitive.E{Key: "_id", Value: data.ID}}, data, optionsDoUpsert)
	if err := res.Err(); err != nil {
		return err
	}
	return nil
}

func Insert(data AnalysisData) error {
	coll := database.DB.Collection("AnalysisData")
	_, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}
	return nil
}

func Delete(taskID int32) error {
	coll := database.DB.Collection("AnalysisData")
	res, err := coll.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "taskID", Value: taskID}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("未找到可删除的历史数据")
	}
	return nil
}

func DeleteByGroup(group int64) error {
	coll := database.DB.Collection("AnalysisData")
	_, err := coll.DeleteMany(context.TODO(), bson.D{primitive.E{Key: "group", Value: group}})
	return err
}

// Refresh data.Time to time.Now().Unix() for newest result
func RenewAll() {
	coll := database.DB.Collection("AnalysisData")
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "time", Value: time.Now().Unix()}}}}
	res, err := coll.UpdateMany(context.TODO(), bson.D{{}}, update)
	if err != nil {
		panic(err)
	}
	log.Printf("Found %d, Renew %d AnalysisDatas to current time\n", res.MatchedCount, res.ModifiedCount)
}

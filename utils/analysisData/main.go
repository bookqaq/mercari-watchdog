package analysisdata

import (
	"context"
	"fmt"

	"bookq.xyz/mercariWatchdog/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	res := coll.FindOneAndReplace(context.TODO(), bson.D{primitive.E{Key: "_id", Value: data.ID}}, data)
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

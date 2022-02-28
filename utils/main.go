package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

func Init() {
	db = Connect()
}

func SearchAllTasks(interval int) ([]AnalysisTask, error) {
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

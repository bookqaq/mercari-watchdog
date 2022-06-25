package fetchdata

import (
	"context"
	"time"

	"bookq.xyz/mercari-watchdog/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Insert(data TaskAddFetchData) error {
	coll := database.DB.Collection("TaskAddFetch")
	_, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}
	return nil
}

func GetOne(auth string) (TaskAddFetchData, error) {
	coll := database.DB.Collection("TaskAddFetch")
	res := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "auth", Value: auth}})
	if err := res.Err(); err != nil {
		return TaskAddFetchData{}, err
	}
	var result TaskAddFetchData
	err := res.Decode(&result)
	if err != nil {
		return TaskAddFetchData{}, err
	}
	return result, nil
}

func IfExist(auth string) bool {
	coll := database.DB.Collection("TaskAddFetch")
	if err := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "auth", Value: auth}}).Err(); err != nil {
		return false
	}
	return true
}

func Delete(auth string) bool {
	coll := database.DB.Collection("TaskAddFetch")
	if err := coll.FindOneAndDelete(context.TODO(), bson.D{primitive.E{Key: "auth", Value: auth}}).Err(); err != nil {
		return false
	}
	return true
}

// remove documents that expired
func ClearExpired() {
	coll := database.DB.Collection("TaskAddFetch")
	_, err := coll.DeleteMany(
		context.TODO(), bson.D{{
			Key: "expire",
			Value: bson.D{primitive.E{
				Key:   "$lte",
				Value: time.Now().Unix()}}}},
	)
	if err != nil && err != mongo.ErrNoDocuments {
		panic(err)
	}
}

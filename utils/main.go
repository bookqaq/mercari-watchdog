package utils

import (
	"context"

	"bookq.xyz/mercari-watchdog/database"
	"bookq.xyz/mercari-watchdog/utils/analysistask"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Init() {
	database.Connect()
	go analysistask.AddTaskBuffer()
}

// Check group invite whitelist
func FindWhitelist(group int64) (bool, error) {
	coll := database.DB.Collection("GroupWhitelist")
	res := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "group", Value: group}})
	if res.Err() != nil {
		return false, res.Err()
	}
	return true, nil
}

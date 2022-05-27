package group

import (
	"context"

	"bookq.xyz/mercari-watchdog/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var coll = database.DB.Collection("GroupSettings")

// Check group invite whitelist
func FindWhitelist(group int64) (bool, error) {
	res := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "group", Value: group}})
	if res.Err() != nil {
		return false, res.Err()
	}
	return true, nil
}

func Get(group int64) (Settings, error) {
	var result Settings
	res := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "group", Value: group}})
	if err := res.Decode(&result); err != nil {
		return Settings{}, err
	}
	return result, nil
}

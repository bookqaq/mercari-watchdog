package group

import (
	"context"

	"bookq.xyz/mercari-watchdog/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Check group invite whitelist
func FindWhitelist(group int64) (bool, error) {
	coll := database.DB.Collection("GroupSettings")
	res := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "group", Value: group}})
	if res.Err() != nil {
		return false, res.Err()
	}
	return true, nil
}

func Get(group int64) (Settings, error) {
	coll := database.DB.Collection("GroupSettings")
	var result Settings
	res := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "group", Value: group}})
	if err := res.Decode(&result); err != nil {
		return Settings{}, err
	}
	return result, nil
}

package blacklist

import (
	"context"

	"bookq.xyz/mercari-watchdog/database"
	"go.mongodb.org/mongo-driver/bson"
)

func BlockedSellerGetAll() ([]BlockedSeller, error) {
	coll := database.DB.Collection("BlackList")
	var result []BlockedSeller
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

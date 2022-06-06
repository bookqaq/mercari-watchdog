package group

import "go.mongodb.org/mongo-driver/bson/primitive"

type Settings struct {
	ID        primitive.ObjectID `bson:"_id"`
	Group     int64              `bson:"group"`
	Privleged bool               `bson:"privleged"`
}

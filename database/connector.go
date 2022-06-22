package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Connect() {

	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017/?readPreference=primary&ssl=false")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Printf("Connect to database fail.\n%s", err)
		os.Exit(-1)
	}
	db := client.Database("merbotV3DB")

	DB = db
}

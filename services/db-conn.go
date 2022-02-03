package services

import (
	"context"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DbInstance *mongo.Client
var dbConnOnce sync.Once

func dBConn() *mongo.Client {
	dbConnOnce.Do(func() {
		clientOptions := options.Client().
			ApplyURI(os.Getenv("MONGO_DB_CONN_URL"))
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}
		DbInstance = client
	})

	return DbInstance
}

func DbDisconnect() {
	if err := DbInstance.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
}

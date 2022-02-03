package services

import (
	"context"
	"sync"

	"github.com/malhotra-sidharth/url-shortener-go/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type IDbOperations interface {
	InsertOne(document *models.UrlRecord) (interface{}, error)
}

type db struct {
	client *mongo.Client
}

func newDB(client *mongo.Client) IDbOperations {
	dbOnce.Do(func() {
		urlsCollection = client.Database("url_shortener").Collection(("urls"))
	})
	return &db{
		client: client,
	}
}

var urlsCollection *mongo.Collection
var dbOnce sync.Once

func (db *db) InsertOne(document *models.UrlRecord) (interface{}, error) {
	result, err := urlsCollection.InsertOne(context.TODO(), document)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

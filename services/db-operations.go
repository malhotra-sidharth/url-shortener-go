package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/malhotra-sidharth/url-shortener-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IDbOperations interface {
	InsertOne(document *models.UrlRecord) (interface{}, error)
	FindOneById(id []byte) (*models.UrlRecord, error)
	UpdatedOneById(id []byte, update bson.D) (*int64, error)
	DeleteOneById(id []byte) (*int64, error)
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

func (db *db) FindOneById(id []byte) (*models.UrlRecord, error) {
	filter := bson.D{{"_id", bson.D{{"$eq", id}}}}
	var result *models.UrlRecord
	if err := urlsCollection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil
}

func (db *db) UpdatedOneById(id []byte, update bson.D) (*int64, error) {
	updatedRecord, err := urlsCollection.UpdateByID(context.TODO(), id, update)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &updatedRecord.ModifiedCount, nil
}

func (db *db) DeleteOneById(id []byte) (*int64, error) {
	filter := bson.D{{"_id", id}}
	result, err := urlsCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	return &result.DeletedCount, nil
}

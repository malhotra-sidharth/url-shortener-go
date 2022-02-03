package services

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"time"

	"github.com/malhotra-sidharth/url-shortener-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IShortener interface {
	Create(url string) (*string, error)
	ResolveUrl(id string) (*models.UrlRecord, error)
}

type shortener struct {
	db IDbOperations
}

func newShortener(db IDbOperations) IShortener {
	return &shortener{
		db: db,
	}
}

func generateId(val string) []byte {
	sha := sha1.New()
	sha.Write([]byte(val))
	return sha.Sum(nil)
}

func isExpired(currentTime uint64, lastVisited uint64) bool {
	secondsInYear := 31536000
	return (currentTime-lastVisited >= uint64(secondsInYear))
}

func (shortener *shortener) Create(url string) (*string, error) {
	// check if already exists
	id := generateId(url)
	result, err := shortener.db.FindOneById(id)
	if result != nil {
		return nil, errors.New("Records already exists")
	}

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	currentTime := uint64(time.Now().Unix())

	// insert new record
	document := &models.UrlRecord{Id: id, IdString: hex.EncodeToString(id), FullUrl: url, Visited: []uint64{}, CreatedAt: currentTime, LastVisited: currentTime}
	_, insertionErr := shortener.db.InsertOne(document)
	if insertionErr != nil {
		return nil, insertionErr
	}
	return &document.IdString, nil
}

func (shortener *shortener) ResolveUrl(id string) (*models.UrlRecord, error) {
	byteId, err := hex.DecodeString(id)
	if err != nil {
		return nil, err
	}
	record, err := shortener.db.FindOneById(byteId)

	if err != nil {
		return nil, err
	}

	visited := uint64(time.Now().Unix())

	if isExpired(visited, record.LastVisited) {
		return nil, errors.New("Expired Url")
	}

	update := bson.D{
		{
			"$push", bson.D{
				{"visited", visited},
			},
		},
		{
			"$set", bson.D{
				{"lastVisited", visited},
			},
		},
	}

	updatedCount, updatedErr := shortener.db.UpdatedOneById(byteId, update)

	if updatedErr != nil {
		return nil, updatedErr
	}

	if *updatedCount > 0 {
		record.Visited = nil
		return record, nil
	}

	return nil, errors.New("Record Not Found")
}

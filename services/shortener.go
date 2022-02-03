package services

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"time"

	"github.com/malhotra-sidharth/url-shortener-go/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type IShortener interface {
	Create(url string) (*string, error)
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

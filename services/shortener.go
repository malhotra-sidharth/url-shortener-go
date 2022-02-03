package services

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"sort"
	"time"

	"github.com/malhotra-sidharth/url-shortener-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IShortener interface {
	Create(url string) (*string, error)
	ResolveUrl(id string) (*models.UrlRecord, error)
	DeleteUrl(id string) (*int64, error)
	AccessCount(id string) (*models.Analytics, error)
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

func sortUintSliceDesc(visited []uint64) {
	sort.Slice(visited, func(i, j int) bool {
		return visited[i] > visited[j]
	})
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

func (shortener *shortener) DeleteUrl(id string) (*int64, error) {
	byteId, err := hex.DecodeString(id)
	if err != nil {
		return nil, err
	}

	return shortener.db.DeleteOneById(byteId)
}

func (shortener *shortener) AccessCount(id string) (*models.Analytics, error) {
	byteId, err := hex.DecodeString(id)
	if err != nil {
		return nil, err
	}
	record, err := shortener.db.FindOneById(byteId)

	if err != nil {
		return nil, err
	}

	currentTimestamp := time.Now().Unix()
	timestampLastDay := currentTimestamp - 86400
	timestampLastWeek := currentTimestamp - 604800

	analytics := &models.Analytics{
		Url:      *record,
		LastDay:  0,
		LastWeek: 0,
		AllTime:  0,
	}

	sortUintSliceDesc(record.Visited)

	for _, timestamp := range record.Visited {
		if timestamp >= uint64(timestampLastDay) {
			analytics.LastDay++
			analytics.LastWeek++
		} else if timestamp >= uint64(timestampLastWeek) {
			analytics.LastWeek++
		} else {
			break
		}
	}

	analytics.AllTime = uint64(len(record.Visited))
	analytics.Url.Visited = nil
	return analytics, nil
}

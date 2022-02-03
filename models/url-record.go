package models

type UrlRecord struct {
	Id          []byte   `bson:"_id" json:"id"`
	IdString    string   `bson:"idString, omitempty" json:"idString"`
	FullUrl     string   `bson:"fullUrl" json:"fullUrl"`
	CreatedAt   uint64   `bson:"createdAt" json:"createdAt"`
	Visited     []uint64 `bson:"visited, omitempty" json:"visited"`
	LastVisited uint64   `bson:"lastVisited, omitempty" json:"lastVisited"`
}

type CreateShortUrlPayload struct {
	FullUrl string `json:"url" binding:"required"`
}

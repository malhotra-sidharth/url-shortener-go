package models

type Analytics struct {
	Url      UrlRecord `json:"url"`
	LastDay  uint64    `json:"lastDay"`
	LastWeek uint64    `json:"lastWeek"`
	AllTime  uint64    `json:"allTime"`
}

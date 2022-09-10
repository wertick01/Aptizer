package models

import "time"

type News struct {
	ID     int64     `json:"id"`
	Text   string    `json:"text"`
	Photo  string    `json:"photo"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
	Author []*User   `json:"author"`
	Tag    []*Tag    `json:"tag"`
}

type Tag struct {
	TagID string `json:"tag_id"`
	Tag   string `json:"tag"`
}
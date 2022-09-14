package models

type News struct {
	ID           int64   `json:"id"`
	Text         string  `json:"text"`
	Photo        string  `json:"photo"`
	Title        string  `json:"title"`
	Date         int64   `json:"date"`
	Author       *User   `json:"author"`
	Participants []*User `json:"participants"`
	Tag          []*Tag  `json:"tag"`
	Updated_at   int64   `json:"updated_at"`
}

type Tag struct {
	TagID int64  `json:"tag_id"`
	Tag   string `json:"tag"`
}

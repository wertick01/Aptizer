package db

import (
	"database/sql"

	"aptizer.com/internals/app/models"
	"github.com/huandu/go-sqlbuilder"
)

type NewsStorage struct {
	database *sql.DB
	smth     string
}

func NewNewsStorage(db *sql.DB) *NewsStorage {
	storage := new(NewsStorage)
	storage.database = db
	return storage
}

func (news *NewsStorage) List() ([]*models.News, error) {
	rows, err := news.database.Query(
		sqlbuilder.Select(
			"n.id",
			"n.headline_date",
			"n.title",
			"tg.tag_id",
			"tg.tag",
			"n.headline_text",
			"n.photo",
			"n.author_id",
			"u.username",
			"u.usersurname",
			"u.userpatrynomic",
		).From(
			"tags_news tn",
		).JoinWithOption(
			sqlbuilder.LeftJoin,
			"news n",
			"tn.headline_id = n.id",
		).JoinWithOption(
			sqlbuilder.InnerJoin,
			"tags tg",
			"tn.tag_id=tg.tag_id",
		).JoinWithOption(
			sqlbuilder.InnerJoin,
			"users u",
			"n.author_id=u.userid",
		).String(),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listNews []*models.News

	for rows.Next() {
		headline := &models.News{}
		author := &models.User{}
		tag := &models.Tag{}
		if err := rows.Scan(
			&headline.ID,
			&headline.Date,
			&headline.Title,
			&tag.TagID,
			&tag.Tag,
			&headline.Text,
			&headline.Photo,
			&author.UserID,
			&author.Name,
			&author.Surname,
			&author.Patrynomic,
		); err != nil {
			return nil, err
		}
		headline.Author = append(headline.Author, author)
		headline.Tag = append(headline.Tag, tag)

		listNews = append(listNews, headline)
	}

	lasttag := &models.Tag{}
	var headline_id int64 = 0
	var counter int = 0

	for key, value := range listNews {
		if value.ID == headline_id {
			counter++
			listNews = remove(listNews, key-counter)
			value.Tag = append(value.Tag, lasttag)
			value.Tag = Reverse(value.Tag)
		}
		lasttag = value.Tag[0]
		headline_id = value.ID
	}
	return listNews, nil
}

func Reverse(input []*models.Tag) []*models.Tag {
	var output []*models.Tag
	for i := len(input) - 1; i >= 0; i-- {
		output = append(output, input[i])
	}
	return output
}

func remove(slice []*models.News, s int) []*models.News {
	return append(slice[:s], slice[s+1:]...)
}

package db

import (
	"database/sql"
	"fmt"
	"time"

	"aptizer.com/internal/app/models"
	"github.com/doug-martin/goqu/v9"
	"github.com/huandu/go-sqlbuilder"
)

type NewsStorage struct {
	database *sql.DB
	smth     string
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
			"pt.userid",
			"us.username",
			"us.usersurname",
			"us.userpatrynomic",
			"us.userphoto",
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
		).JoinWithOption(
			sqlbuilder.InnerJoin,
			"participants pt",
			"pt.headline_id=n.id",
		).JoinWithOption(
			sqlbuilder.InnerJoin,
			"users us",
			"us.userid=pt.userid",
		).String(),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listNews []*models.News

	for rows.Next() {
		headline := &models.News{Author: &models.User{}}
		participant := &models.User{}
		tag := &models.Tag{}
		if err := rows.Scan(
			&headline.ID,
			&headline.Date,
			&headline.Title,
			&tag.TagID,
			&tag.Tag,
			&headline.Text,
			&headline.Photo,
			&headline.Author.UserID,
			&headline.Author.Name,
			&headline.Author.Surname,
			&headline.Author.Patrynomic,
			&participant.UserID,
			&participant.Name,
			&participant.Surname,
			&participant.Patrynomic,
			&participant.Photo,
		); err != nil {
			return nil, err
		}
		headline.Participants = append(headline.Participants, participant)
		headline.Tag = append(headline.Tag, tag)

		listNews = append(listNews, headline)
	}
	fmt.Println(listNews)

	// var tag_mass []int64
	//counter := 0

	for key, value := range listNews[1:] {
		fmt.Println(key)
		if value.ID == listNews[key].ID {
			for _, i := range listNews[key].Tag {
				value.Tag = append(value.Tag, i)
				//tag_mass = append(tag_mass, )
			}
			for _, i := range listNews[key].Participants {
				value.Participants = append(value.Participants, i)
			}
			//listNews = remove(listNews, key)
			//counter++
		}
		fmt.Println(value)
		tag_mass := []int64{}
		for m, k := range value.Tag[1:] {
			tag_mass = append(tag_mass, value.Tag[m].TagID)
			if checker(tag_mass, k.TagID) {
				value.Tag = tagRemove(value.Tag, int64(m))
			}
		}
		fmt.Println(value)
		//counter = 0
		user_mass := []int64{}
		for g, _ := range value.Participants {
			user_mass = append(user_mass, value.Participants[g].UserID)
			fmt.Println(user_mass)
		}
		fmt.Println(value.Tag, value.Participants, user_mass)
		fmt.Println()
		// for _, j := range value.Participants {
		// 	fmt.Printf("%v ", j)
		// }
		// fmt.Println()
	}

	return listNews, nil
}

// var CurrentTime = time.Now().Format("2006-01-02 15:04:05")

func (news *NewsStorage) Create(headline *models.News) (*models.News, error) {
	headline.Date = time.Now().Unix()

	dialect := goqu.Dialect("mysql")

	insertquery, _, err := dialect.Insert(
		"users",
	).Prepared(
		true,
	).Rows(
		goqu.Record{
			"headline_date": headline.Date,
			"title":         headline.Title,
			"headline_text": headline.Text,
			"photo":         headline.Photo,
			"author_id":     headline.Author,
		},
	).ToSQL()
	if err != nil {
		return nil, err
	}

	if _, err := news.database.Exec(
		insertquery,
		headline.Date,
		headline.Title,
		headline.Text,
		headline.Photo,
		headline.Author,
	); err != nil {
		return nil, err
	}
	if err := news.database.QueryRow(
		`SELECT LAST_INSERT_ID()`,
	).Scan(
		&headline.ID,
	); err != nil {
		return nil, err
	}

	headline, err = news.CreateParticipants(headline)
	if err != nil {
		return nil, err
	}

	headline, err = news.CreateTags(headline)
	if err != nil {
		return nil, err
	}

	return headline, nil
}

func (news *NewsStorage) Update(headline *models.News) (*models.News, error) {
	headline.Updated_at = time.Now().Unix()

	dialect := goqu.Dialect("mysql")

	updatequery, _, err := dialect.Update(
		"news",
	).Set(
		goqu.Record{
			"title":         headline.Title,
			"headline_text": headline.Text,
			"photo":         headline.Photo,
			"author_id":     headline.Author,
			"updated_at":    headline.Updated_at,
		},
	).Prepared(
		true,
	).Where(
		goqu.Ex{
			"id": headline.ID,
		},
	).ToSQL()
	if err != nil {
		return nil, err
	}

	if _, err := news.database.Exec(
		updatequery,
		headline.Title,
		headline.Text,
		headline.Photo,
		headline.Author,
		headline.Updated_at,
	); err != nil {
		return nil, err
	}

	return headline, nil
}

func (news *NewsStorage) CreateParticipants(headline *models.News) (*models.News, error) {
	dialect := goqu.Dialect("mysql")

	insertbase := dialect.Insert(
		"participants",
	).Prepared(
		true,
	)

	for _, participant := range headline.Participants {
		query, _, err := insertbase.Rows(
			goqu.Record{
				"userid":      participant.UserID,
				"headline_id": headline.ID,
			},
		).ToSQL()
		if err != nil {
			return nil, err
		}

		if _, err := news.database.Exec(
			query,
			participant.UserID,
			headline.ID,
		); err != nil {
			return nil, err
		}
	}
	return headline, nil
}

func (news *NewsStorage) CreateTags(headline *models.News) (*models.News, error) {
	dialect := goqu.Dialect("mysql")

	insertbase1 := dialect.Insert(
		"tags",
	).Prepared(
		true,
	)

	insertbase2 := dialect.Insert(
		"tags_news",
	).Prepared(
		true,
	)

	for _, tag := range headline.Tag {
		insertquery1, _, err := insertbase1.Rows(
			goqu.Record{
				"tag": tag.Tag,
			},
		).ToSQL()
		if err != nil {
			return nil, err
		}

		if _, err := news.database.Exec(
			insertquery1,
			tag.Tag,
		); err != nil {
			return nil, err
		}

		if err := news.database.QueryRow(
			`SELECT LAST_INSERT_ID()`,
		).Scan(
			&tag.TagID,
		); err != nil {
			return nil, err
		}

		insertquery2, _, err := insertbase2.Rows(
			goqu.Record{
				"headline_id": headline.ID,
				"tag_id":      tag.TagID,
			},
		).ToSQL()
		if err != nil {
			return nil, err
		}
		if _, err := news.database.Exec(
			insertquery2,
			headline.ID,
			tag.TagID,
		); err != nil {
			return nil, err
		}
	}
	return headline, nil
}

func (news *NewsStorage) DeleteTags(tag *models.Tag) error {
	dialect := goqu.Dialect("mysql")

	deletequery, _, err := dialect.Delete(
		"tags",
	).Where(
		goqu.Ex{
			"tag_id": tag.TagID,
		},
	).ToSQL()
	if err != nil {
		return err
	}

	if _, err := news.database.Exec(
		deletequery,
		tag.TagID,
	); err != nil {
		return err
	}

	return nil
}

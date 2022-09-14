package db

import (
	"database/sql"
	"fmt"

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

func checker(mass []int64, val int64) bool {
	for _, i := range mass {
		if i == val {
			return true
		}
	}
	return false
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

func tagRemove(slice []*models.Tag, s int64) []*models.Tag {
	return append(slice[:s], slice[s+1:]...)
}

func userRemove(slice []*models.User, s int64) []*models.User {
	return append(slice[:s], slice[s+1:]...)
}

func sliceRemove(slice []int64, s int64) []int64 {
	return append(slice[:s], slice[s+1:]...)
}

package db

import (
	"aptizer.com/internal/app/models"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"

	"database/sql"
	"errors"
	"time"
)

type UsersStorage struct {
	DB *sql.DB
}

// CreateHash - Necessary for hashing the user's password.
func (m *UsersStorage) CreateHash(user *models.User) *models.User {
	hash, err := m.HashPassword(user.Hash)
	if err != nil {
		return nil
	}
	user.Hash = hash
	return user
}

// CreateNewUser - Necessary to enter the data of a new user in the database.
func (m *UsersStorage) CreateNewUser(user *models.User) (*models.User, error) {
	user = m.CreateHash(user)

	dialect := goqu.Dialect("mysql")

	query, _, err := dialect.Insert(
		"users",
	).Prepared(
		true,
	).Rows(
		goqu.Record{
			"username":        user.Name,
			"usersurname":     user.Surname,
			"userpatrynomic":  user.Patrynomic,
			"useremail":       user.Mail,
			"userphone":       user.Phone,
			"userhash":        user.Hash,
			"userdescription": user.Description,
			"userphoto":       user.Photo,
			"userrole":        user.Role.RoleID,
		},
	).ToSQL()
	if err != nil {
		return nil, err
	}

	if _, err := m.DB.Exec(
		query,
		user.Name,
		user.Surname,
		user.Patrynomic,
		user.Phone,
		user.Mail,
		user.Hash,
		user.Role.RoleID,
	); err != nil {
		return nil, err
	}
	if err := m.DB.QueryRow(
		`SELECT LAST_INSERT_ID()`,
	).Scan(
		&user.UserID,
	); err != nil {
		return nil, err
	}
	user.Hash = ""

	return user, nil
}

// SetRefreshToken - Check and save the refresh-token to the database.
func (m *UsersStorage) SetRefreshToken(rt *models.RefreshToken) error {
	dialect := goqu.Dialect("mysql")

	selectquery, _, err := dialect.From(
		"refresh_token",
	).Select(
		"id",
	).Where(
		goqu.Ex{
			"useragent": rt.UserAgent,
			"userid":    rt.UserID,
		},
	).ToSQL()
	if err != nil {
		return err
	}

	if err := m.DB.QueryRow(selectquery).Scan(
		&rt.ID,
	); err != nil {
		if err == sql.ErrNoRows {

			insertquery, _, err := dialect.Insert(
				"refresh_token",
			).Prepared(
				true,
			).Rows(
				goqu.Record{
					"userid":        rt.UserID,
					"refresh_token": rt.RefreshToken,
					"useragent":     rt.UserAgent,
					"expires_in":    rt.ExpiresIn,
				},
			).ToSQL()
			if err != nil {
				return err
			}

			if _, err := m.DB.Exec(
				insertquery,
				rt.UserID,
				rt.RefreshToken,
				rt.UserAgent,
				rt.ExpiresIn,
			); err != nil {
				return err
			}
			if err := m.DB.QueryRow(`SELECT LAST_INSERT_ID()`).Scan(&rt.ID); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	updatequery, _, err := dialect.Update(
		"refresh_token",
	).Prepared(
		true,
	).Set(
		goqu.Record{
			"refresh_token": rt.RefreshToken,
			"expires_in":    rt.ExpiresIn,
			"updated_at":    goqu.L("NOW()"),
		},
	).Where(
		goqu.Ex{
			"id": rt.ID,
		},
	).ToSQL()
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(
		updatequery,
		rt.RefreshToken,
		rt.ExpiresIn,
		rt.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// CheckRefreshToken - Check the expiration date and validity of the user's refresh-token with the one in the database.
func (m *UsersStorage) CheckRefreshToken(rt *models.RefreshToken) (*models.RefreshToken, error) {
	dialect := goqu.Dialect("mysql")

	selectquery, _, err := dialect.From(
		"refresh_token",
	).Select(
		"id",
		"expires_in",
		"refresh_token",
	).Where(
		goqu.Ex{
			"userid":    rt.UserID,
			"useragent": rt.UserAgent,
		},
	).ToSQL()
	if err != nil {
		return nil, err
	}

	if err := m.DB.QueryRow(selectquery).Scan(
		&rt.ID,
		&rt.ExpiresIn,
		&rt.RefreshToken,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, m.SetRefreshToken(rt)
		} else {
			return nil, err
		}
	}

	now := time.Now().Unix()
	if now > rt.ExpiresIn {
		return nil, models.ExpiredToken
	}

	return rt, nil
}

// GetUsersList - Returns a list of all registered users.
func (m *UsersStorage) GetUsersList() ([]*models.User, error) {
	dialect := goqu.Dialect("mysql")

	selectquery, _, err := dialect.From(
		goqu.T(
			"users",
		).As(
			"u",
		),
	).Select(
		"u.userid",
		"u.username",
		"u.usersurname",
		"u.userpatrynomic",
		"u.userphone",
		"u.useremail",
		"u.userdescription",
		"u.userphoto",
		"u.userrole",
		"r.user_role",
	).Where(
		goqu.Ex{
			"u.userid": goqu.Op{
				"isNot": nil,
			},
		},
	).RightJoin(
		goqu.T(
			"roles",
		).As(
			"r",
		),
		goqu.On(
			goqu.Ex{
				"u.userrole": goqu.I(
					"r.role_id",
				),
			},
		),
	).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := m.DB.Query(selectquery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		s := &models.User{}
		err = rows.Scan(
			&s.UserID,
			&s.Name,
			&s.Surname,
			&s.Patrynomic,
			&s.Phone,
			&s.Mail,
			&s.Description,
			&s.Photo,
			&s.Role.RoleID,
			&s.Role.Role,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID - Returns user data by his id (userid).
func (m *UsersStorage) GetUserByID(id int64) (*models.User, error) {
	dialect := goqu.Dialect("mysql")

	selectquery, _, err := dialect.From(
		goqu.T(
			"users",
		).As(
			"u",
		),
	).Select(
		"u.userid",
		"u.username",
		"u.usersurname",
		"u.userpatrynomic",
		"u.userphone",
		"u.useremail",
		"u.userdescription",
		"u.userphoto",
		"u.userrole",
		"r.user_role",
	).Where(
		goqu.Ex{
			"u.userid": id,
		},
	).RightJoin(
		goqu.T(
			"roles",
		).As(
			"r",
		),
		goqu.On(
			goqu.Ex{
				"u.userrole": goqu.I(
					"r.role_id",
				),
			},
		),
	).ToSQL()
	if err != nil {
		return nil, err
	}

	row := m.DB.QueryRow(selectquery)

	s := &models.User{}

	if err := row.Scan(
		&s.UserID,
		&s.Name,
		&s.Surname,
		&s.Patrynomic,
		&s.Phone,
		&s.Mail,
		&s.Description,
		&s.Photo,
		&s.Role.RoleID,
		&s.Role.Role,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return s, nil
}

// GetUserByPhone - Returns user data by his phone (userphone).
func (m *UsersStorage) GetUserByPhone(phone string) (*models.User, error) {
	dialect := goqu.Dialect("mysql")

	selectquery, _, err := dialect.From(
		goqu.T(
			"users",
		).As(
			"u",
		),
	).Prepared(
		true,
	).Select(
		"u.userid",
		"u.username",
		"u.usersurname",
		"u.userpatrynomic",
		"u.userphone",
		"u.useremail",
		"u.userdescription",
		"u.userphoto",
		"u.userrole",
		"r.user_role",
	).Where(
		goqu.Ex{
			"u.userphone": phone,
		},
	).RightJoin(
		goqu.T(
			"roles",
		).As(
			"r",
		),
		goqu.On(
			goqu.Ex{
				"u.userrole": goqu.I(
					"r.role_id",
				),
			},
		),
	).ToSQL()
	if err != nil {
		return nil, err
	}

	row := m.DB.QueryRow(selectquery, phone)

	s := &models.User{}

	if err := row.Scan(
		&s.UserID,
		&s.Name,
		&s.Surname,
		&s.Patrynomic,
		&s.Phone,
		&s.Mail,
		&s.Description,
		&s.Photo,
		&s.Role.RoleID,
		&s.Role.Role,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return s, nil
}

// ChangeUser - Makes changes to the user's data.
func (m *UsersStorage) ChangeUser(old *models.User) (*models.User, error) {

	hash, err := m.HashPassword(old.Hash)
	if err != nil {
		return nil, err
	}
	old.Hash = hash

	dialect := goqu.Dialect("mysql")

	updatequery, _, err := dialect.Update(
		"users",
	).Set(
		goqu.Record{
			"username":        old.Name,
			"usersurname":     old.Surname,
			"userpatrynomic":  old.Patrynomic,
			"useremail":       old.Mail,
			"userphone":       old.Phone,
			"userhash":        old.Hash,
			"userdescription": old.Description,
			"userphoto":       old.Photo,
			"userrole":        old.Role.RoleID,
		},
	).Prepared(
		true,
	).Where(
		goqu.Ex{
			"userid": old.UserID,
		},
	).ToSQL()
	if err != nil {
		return nil, err
	}

	_, err = m.DB.Exec(
		updatequery,
		old.Name,
		old.Surname,
		old.Patrynomic,
		old.Mail,
		old.Phone,
		old.Hash,
		old.Description,
		old.Photo,
		old.Role.RoleID,
	)
	if err != nil {
		return nil, err
	}
	old.Hash = ""

	return old, nil
}

// DeleteUserByID - Deletes all data (account) of the user.
func (m *UsersStorage) DeleteUserByID(id int64) (int64, error) {
	dialect := goqu.Dialect("mysql")

	deletequery_1, _, err := dialect.Delete(
		"refresh_token",
	).Where(
		goqu.Ex{
			"userid": id,
		},
	).ToSQL()
	if err != nil {
		return 0, err
	}

	deletequery_2, _, err := dialect.Delete(
		"users",
	).Where(
		goqu.Ex{
			"userid": id,
		},
	).ToSQL()
	if err != nil {
		return 0, err
	}

	if _, err = m.DB.Exec(deletequery_1); err != nil {
		return 0, err
	}

	if _, err = m.DB.Exec(deletequery_2); err != nil {
		return 0, err
	}
	return id, nil
}

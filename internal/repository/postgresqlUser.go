package repository

import (
	"context"

	model "github.com/AndiVS/myapp3.0/internal/model"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

// InsertUser function for inserting item from a table
func (repos *Postgres) InsertUser(c context.Context, user *model.User) error {
	row := repos.pool.QueryRow(c,
		"INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3) RETURNING username", user.Username, user.Password, false)

	err := row.Scan(&user.Username)
	if err != nil {
		log.Errorf("Unable to INSERT: %v", err)
		return err
	}

	return err
}

// SelectUser function for selecting item from a table
func (repos *Postgres) SelectUser(c context.Context, username string) (*model.User, error) {
	var user model.User
	row := repos.pool.QueryRow(c,
		"SELECT username, password, is_admin  FROM users WHERE username = $1", username)

	err := row.Scan(&user.Username, &user.Password, &user.IsAdmin)
	if err != nil {
		log.Errorf("Unable to SELECT: %v", err)
		return &user, err
	}

	log.Printf("sec")

	return &user, err
}

// SelectAllUser function for selecting items from a table
func (repos *Postgres) SelectAllUser(c context.Context) ([]*model.User, error) {
	var user []*model.User

	row, err := repos.pool.Query(c,
		"SELECT username, password, is_admin  FROM users")

	for row.Next() {
		var us model.User
		err = row.Scan(&us.Username, &us.Password, &us.IsAdmin)
		if err == pgx.ErrNoRows {
			return user, err
		}
		user = append(user, &us)
	}

	return user, err
}

// UpdateUser function for updating item from a table
func (repos *Postgres) UpdateUser(c context.Context, username string, isAdmin bool) error {
	_, err := repos.pool.Exec(c,
		"UPDATE users SET is_admin = $2 WHERE username = $1", username, isAdmin)
	if err != nil {
		log.Errorf("Failed updating data in db: %s\n", err)
		return err
	}

	return nil
}

// DeleteUser function for deleting item from a table
func (repos *Postgres) DeleteUser(c context.Context, username string) error {
	ct, err := repos.pool.Exec(c, "DELETE FROM users WHERE username = $1", username)

	if err != nil {
		return err
	} else if ct.RowsAffected() == 0 {
		log.Errorf("Not found : %s\n", err)
		return ErrNotFound
	}

	return nil
}

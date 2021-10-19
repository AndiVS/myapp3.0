package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	model "myapp3.0/internal/model"
)

// Users used for structuring, function for working with users
type Users interface {
	InsertU(c context.Context, rec *model.User) error
	SelectU(c context.Context, username, password string) (*model.User, error)
	SelectAllU(c context.Context) ([]*model.User, error)
	UpdateU(c context.Context, username string, isAdmin bool) error
	DeleteU(c context.Context, username string) error
}

// InsertU function for inserting item from a table
func (repos *Postgres) InsertU(c context.Context, rec *model.User) error {
	row := repos.pool.QueryRow(c,
		"INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3) RETURNING username", rec.Username, rec.Password, false)

	err := row.Scan(&rec.Username)
	if err != nil {
		log.Errorf("Unable to INSERT: %v", err)
		return err
	}

	return err
}

// SelectU function for selecting item from a table
func (repos *Postgres) SelectU(c context.Context, username, password string) (*model.User, error) {
	var rc model.User
	row := repos.pool.QueryRow(c,
		"SELECT username, password, is_admin  FROM users WHERE username = $1 AND password = $2", username, password)

	err := row.Scan(&rc.Username, &rc.Password, &rc.IsAdmin)
	if err != nil {
		log.Errorf("Unable to SELECT: %v", err)
		return &rc, err
	}

	log.Printf("sec")

	return &rc, err
}

// SelectAllU function for selecting items from a table
func (repos *Postgres) SelectAllU(c context.Context) ([]*model.User, error) {
	var rec []*model.User

	row, err := repos.pool.Query(c,
		"SELECT username, password, is_admin  FROM users")

	for row.Next() {
		var rc model.User
		err = row.Scan(&rc.Username, &rc.Password, &rc.IsAdmin)
		if err == pgx.ErrNoRows {
			return rec, err
		}
		rec = append(rec, &rc)
	}

	return rec, err
}

// UpdateU function for updating item from a table
func (repos *Postgres) UpdateU(c context.Context, username string, isAdmin bool) error {
	_, err := repos.pool.Exec(c,
		"UPDATE users SET is_admin = $2 WHERE username = $1", username, isAdmin)
	if err != nil {
		log.Errorf("Failed updating data in db: %s\n", err)
		return err
	}

	return nil
}

// DeleteU function for deleting item from a table
func (repos *Postgres) DeleteU(c context.Context, username string) error {
	_, err := repos.pool.Exec(c, "DELETE FROM users WHERE username = $1", username)

	if err != nil {
		return err
	}

	return nil
}

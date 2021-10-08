package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	model "myapp3.0/internal/model"
)

// SelectAllU function for selecting items from a table
func (repos *Repository) SelectAllU(c context.Context) ([]*model.User, error) {
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

// SelectU function for selecting item from a table
func (repos *Repository) SelectU(c context.Context, username string) (*model.User, error) {
	var rc model.User
	row := repos.pool.QueryRow(c,
		"SELECT username, password, is_admin  FROM users WHERE id = $1", username)

	err := row.Scan(&rc.Username, &rc.Password, &rc.IsAdmin)
	if err != nil {
		log.Errorf("Unable to SELECT: %v", err)
		return &rc, err
	}

	log.Printf("sec")

	return &rc, err
}

// InsertU function for inserting item from a table
func (repos *Repository) InsertU(c context.Context, rec *model.User) error {
	row := repos.pool.QueryRow(c,
		"INSERT INTO users (username, password) VALUES ($1, $2)", rec.Username, rec.Password)

	err := row.Scan(&rec.Username)
	if err != nil {
		log.Errorf("Unable to INSERT: %v", err)
		return err
	}

	return err
}

// UpdateU function for updating item from a table
func (repos *Repository) UpdateU(c context.Context, username string, isAdmin bool) error {
	_, err := repos.pool.Exec(c,
		"UPDATE users SET is_admin = $2 WHERE username = $1", username, isAdmin)

	if err != nil {
		log.Errorf("Failed updating data in db: %s\n", err)
		return err
	}

	return nil
}

// DeleteU function for deleting item from a table
func (repos *Repository) DeleteU(c context.Context, username string) error {
	_, err := repos.pool.Exec(c, "DELETE FROM users WHERE username = $1", username)

	if err != nil {
		return err
	}

	return nil
}
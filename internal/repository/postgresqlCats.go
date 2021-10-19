// Package repository contains code for handling different types of databases
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	model "myapp3.0/internal/model"
)

// Cats used for structuring, function for working with records
type Cats interface {
	InsertC(c context.Context, rec *model.Record) error
	SelectC(c context.Context, id uuid.UUID) (*model.Record, error)
	SelectAllC(c context.Context) ([]*model.Record, error)
	UpdateC(c context.Context, rec *model.Record) error
	DeleteC(c context.Context, id uuid.UUID) error
}

// InsertC function for inserting item from a table
func (repos *Postgres) InsertC(c context.Context, rec *model.Record) error {
	rec.ID = uuid.New()
	row := repos.pool.QueryRow(c,
		"INSERT INTO cats (_id, name, type) VALUES ($1, $2, $3) RETURNING _id", rec.ID, rec.Name, rec.Type)

	err := row.Scan(&rec.ID)
	if err != nil {
		log.Errorf("Unable to INSERT: %v", err)
		return err
	}

	return err
}

// SelectC function for selecting item from a table
func (repos *Postgres) SelectC(c context.Context, id uuid.UUID) (*model.Record, error) {
	var rec model.Record
	row := repos.pool.QueryRow(c,
		"SELECT _id, name, type FROM cats WHERE _id = $1", id)

	err := row.Scan(&rec.ID, &rec.Name, &rec.Type)
	if err != nil {
		log.Errorf("Unable to SELECT: %v", err)
		return &rec, err
	}

	log.Printf("sec")

	return &rec, err
}

// SelectAllC function for selecting items from a table
func (repos *Postgres) SelectAllC(c context.Context) ([]*model.Record, error) {
	var rec []*model.Record

	row, err := repos.pool.Query(c,
		"SELECT _id, name, type  FROM cats")

	for row.Next() {
		var rc model.Record
		err = row.Scan(&rc.ID, &rc.Name, &rc.Type)
		if err == pgx.ErrNoRows {
			return rec, err
		}
		rec = append(rec, &rc)
	}

	return rec, err
}

// UpdateC function for updating item from a table
func (repos *Postgres) UpdateC(c context.Context, rec *model.Record) error {
	_, err := repos.pool.Exec(c,
		"UPDATE cats SET name = $2, type = $3 WHERE _id = $1", rec.ID, rec.Name, rec.Type)

	if err != nil {
		log.Errorf("Failed updating data in db: %s\n", err)
		return err
	}

	return nil
}

// DeleteC function for deleting item from a table
func (repos *Postgres) DeleteC(c context.Context, id uuid.UUID) error {
	_, err := repos.pool.Exec(c, "DELETE FROM cats WHERE _id = $1", id)

	if err != nil {
		return err
	}

	return nil
}

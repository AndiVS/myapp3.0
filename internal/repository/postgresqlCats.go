// Package repository contains code for handling different types of databases
package repository

import (
	"context"
	"errors"

	model "github.com/AndiVS/myapp3.0/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

// InsertCat function for inserting item from a table
func (repos *Postgres) InsertCat(c context.Context, cat *model.Cat) (uuid.UUID, error) {
	row := repos.pool.QueryRow(c,
		"INSERT INTO cats (_id, name, type) VALUES ($1, $2, $3) RETURNING _id", cat.ID, cat.Name, cat.Type)

	err := row.Scan(&cat.ID)
	if err != nil {
		log.Errorf("Unable to INSERT: %v", err)
		return cat.ID, err
	}

	return cat.ID, err
}

// SelectCat function for selecting item from a table
func (repos *Postgres) SelectCat(c context.Context, id uuid.UUID) (*model.Cat, error) {
	var cat model.Cat
	row := repos.pool.QueryRow(c,
		"SELECT _id, name, type FROM cats WHERE _id = $1", id)

	err := row.Scan(&cat.ID, &cat.Name, &cat.Type)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Errorf("Not found : %s\n", err)
		return &cat, ErrNotFound
	} else if err != nil {
		return &cat, err
	}

	log.Printf("sec")

	return &cat, err
}

// SelectAllCat function for selecting items from a table
func (repos *Postgres) SelectAllCat(c context.Context) ([]*model.Cat, error) {
	var cats []*model.Cat

	row, err := repos.pool.Query(c,
		"SELECT _id, name, type  FROM cats")

	for row.Next() {
		var rc model.Cat
		err = row.Scan(&rc.ID, &rc.Name, &rc.Type)
		if err == pgx.ErrNoRows {
			return cats, err
		}
		cats = append(cats, &rc)
	}

	return cats, err
}

// UpdateCat function for updating item from a table
func (repos *Postgres) UpdateCat(c context.Context, cat *model.Cat) error {
	_, err := repos.pool.Exec(c,
		"UPDATE cats SET name = $2, type = $3 WHERE _id = $1", cat.ID, cat.Name, cat.Type)

	if err != nil {
		log.Errorf("Failed updating data in db: %s\n", err)
		return err
	}

	return nil
}

// DeleteCat function for deleting item from a table
func (repos *Postgres) DeleteCat(c context.Context, id uuid.UUID) error {
	_, err := repos.pool.Exec(c, "DELETE FROM cats WHERE _id = $1", id)

	if err != nil {
		return err
	}

	return nil
}

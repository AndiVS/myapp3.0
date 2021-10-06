package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"myapp3.0/internal/model"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (repos *Repository) SelectAll(c context.Context) ([]*model.Record, error) {

	var rec []*model.Record

	row, err := repos.pool.Query(c,
		"SELECT id, name, type  FROM catsbase")

	for row.Next() {
		var rc model.Record
		err = row.Scan(&rc.Id, &rc.Name, &rc.Type)
		if err == pgx.ErrNoRows {
			return rec, err
		}
		rec = append(rec, &rc)
	}

	return rec, err
}

func (repos *Repository) Select(id *int, c context.Context) (*model.Record, error) {

	var rec model.Record
	row := repos.pool.QueryRow(c,
		"SELECT id, name, type FROM catsbase WHERE id = $1", id)

	err := row.Scan(&rec.Id, &rec.Name, &rec.Type)
	if err != nil {
		log.Errorf("Unable to SELECT: %v", err)
		return &rec, err

	}

	log.Printf("sec")

	return &rec, err
}

func (repos *Repository) Insert(rec *model.Record, c context.Context) error {

	row := repos.pool.QueryRow(c,
		"INSERT INTO catsbase (name, type) VALUES ($1, $2) RETURNING id", rec.Name, rec.Type)

	err := row.Scan(&rec.Id)
	if err != nil {
		log.Errorf("Unable to INSERT: %v", err)
		return err
	}

	return err
}

func (repos *Repository) Update(rec *model.Record, c context.Context) error {

	_, err := repos.pool.Exec(c,
		"UPDATE catsbase SET name = $2, type = $3 WHERE id = $1", rec.Id, rec.Name, rec.Type)

	if err != nil {
		log.Errorf("Failed updating data in db: %s\n", err)
		return err
	}

	return nil
}

func (repos *Repository) Delete(id *int, c context.Context) error {

	_, err := repos.pool.Exec(c, "DELETE FROM catsbase WHERE id = $1", id)

	if err != nil {
		return err
	}

	return nil
}

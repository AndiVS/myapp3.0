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

func (repos Repository) SelectAll(c context.Context) ([]*model.Record, error) {

	var rec []*model.Record

	conn, err := repos.pool.Acquire(c)
	if err != nil {
		log.Errorf("Unable to acquire a database connection for SelectAll: %v\n", err)
		return rec, err
	}
	defer conn.Release()

	row, err := conn.Query(c,
		"SELECT * FROM catsbase")

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

func (repos Repository) Select(rec *model.Record, c context.Context) (model.Record, error) {

	conn, err := repos.pool.Acquire(c)
	if err != nil {
		log.Printf("Failed reading the request body for addCats: %s\n", err)
		return *rec, err
	}
	defer conn.Release()

	row := conn.QueryRow(c,
		"SELECT id, name, type FROM catsbase WHERE id = $1", rec.Id)

	err = row.Scan(&rec.Id, &rec.Name, &rec.Type)
	if err == pgx.ErrNoRows {
		log.Errorf("No such row: %v", err)
		return *rec, err
	}
	if err != nil {
		log.Errorf("Unable to SELECT: %v", err)
		return *rec, err

	}

	log.Printf("sec")

	return *rec, err
}

func (repos Repository) Insert(rec *model.Record, c context.Context) (uint64, error) {

	var id uint64
	conn, err := repos.pool.Acquire(c)
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		return id, err
	}
	defer conn.Release()

	row := conn.QueryRow(c,
		"INSERT INTO catsbase (name, type) VALUES ($1, $2) RETURNING id", rec.Name, rec.Type)

	err = row.Scan(&id)
	if err != nil {
		log.Errorf("Unable to INSERT: %v", err)
		return id, err
	}

	return id, err
}

func (repos Repository) Update(rec *model.Record, c context.Context) error {

	conn, err := repos.pool.Acquire(c)
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(c,
		"UPDATE catsbase SET name = $2, type = $3 WHERE id = $1", rec.Id, rec.Name, rec.Type)

	if err != nil {
		log.Errorf("Failed updating data in db: %s\n", err)
		return err
	}

	if ct.RowsAffected() == 0 {
		log.Errorf("Failed reading the request body for addCats: %s\n", err)
		return err
	}

	return nil
}

func (repos Repository) Delete(rec *model.Record, c context.Context) error {
	conn, err := repos.pool.Acquire(c)
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(c, "DELETE FROM catsbase WHERE id = $1", rec.Id)

	if err != nil {
		log.Errorf("Unable to DELETE: %v", err)
		return err
	}

	if ct.RowsAffected() == 0 {
		return err
	}

	return nil
}

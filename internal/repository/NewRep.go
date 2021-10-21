package repository

import (
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"

	"reflect"
)

var (
	// ErrNotFound means entity is not found in repository
	ErrNotFound = errors.New("not found")
)

// Postgres struct for pool
type Postgres struct {
	pool *pgxpool.Pool
}

// Mongo struct for pool
type Mongo struct {
	collectionC *mongo.Collection
	collectionU *mongo.Collection
}

// NewRepository set new repository for mongo and postgress
func NewRepository(db interface{}) Cats {
	var pool *pgxpool.Pool
	var mongoDB *mongo.Database

	switch reflect.TypeOf(db) {
	case reflect.TypeOf(pool):
		return &Postgres{pool: db.(*pgxpool.Pool)}
	case reflect.TypeOf(mongoDB):
		return &Mongo{
			collectionC: db.(*mongo.Database).Collection("cats"),
			collectionU: db.(*mongo.Database).Collection("users"),
		}
	}
	return nil
}

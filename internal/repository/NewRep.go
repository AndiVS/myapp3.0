package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

// RepositoryPostgres struct for pool
type RepositoryPostgres struct {
	pool *pgxpool.Pool
}

// RepositoryMongo struct for pool
type RepositoryMongo struct {
	collectionC *mongo.Collection
	collectionU *mongo.Collection
}

func NewRepository(db interface{}) Cats {
	var pool *pgxpool.Pool
	var mongoDB *mongo.Database

	switch reflect.TypeOf(db) {
	case reflect.TypeOf(pool):
		return &RepositoryPostgres{pool: db.(*pgxpool.Pool)}
	case reflect.TypeOf(mongoDB):
		return &RepositoryMongo{
			collectionC: db.(*mongo.Database).Collection("cats"),
			collectionU: db.(*mongo.Database).Collection("user"),
		}
	}
	return nil
}

package repository

import (
	"context"
	"errors"

	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"

	"reflect"
)

var (
	// ErrNotFound means entity is not found in repository
	ErrNotFound = errors.New("not found")
)

// Postgres struct for Pool
type Postgres struct {
	Pool *pgxpool.Pool
}

// Mongo struct for Pool
type Mongo struct {
	collectionCats  *mongo.Collection
	collectionUsers *mongo.Collection
}

// NewRepository set new repository for mongo and postgres
func NewRepository(db interface{}) Cats {
	var pool *pgxpool.Pool
	var mongoDB *mongo.Database

	switch reflect.TypeOf(db) {
	case reflect.TypeOf(pool):
		return &Postgres{Pool: db.(*pgxpool.Pool)}
	case reflect.TypeOf(mongoDB):
		return &Mongo{
			collectionCats:  db.(*mongo.Database).Collection("cats"),
			collectionUsers: db.(*mongo.Database).Collection("users"),
		}
	}
	return nil
}

// Cats used for structuring, function for working with records
type Cats interface {
	InsertCat(c context.Context, cat *model.Cat) (uuid.UUID, error)
	SelectCat(c context.Context, id uuid.UUID) (*model.Cat, error)
	SelectAllCat(c context.Context) ([]*model.Cat, error)
	UpdateCat(c context.Context, cat *model.Cat) error
	DeleteCat(c context.Context, id uuid.UUID) error
}

// Users used for structuring, function for working with users
type Users interface {
	InsertUser(c context.Context, user *model.User) error
	SelectUser(c context.Context, username string) (*model.User, error)
	SelectAllUser(c context.Context) ([]*model.User, error)
	UpdateUser(c context.Context, username string, isAdmin bool) error
	DeleteUser(c context.Context, username string) error
}

package repository

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	model "myapp3.0/internal/model"
)

// InsertU function for inserting item from a table
func (rep *RepositoryMongo) InsertU(c context.Context, rec *model.User) error {
	id, erro := uuid.New()
	if erro != nil {
		log.Errorf("asd", erro)
	}

	cat := model.Us{
		ID:       id,
		Username: rec.Username,
		Password: rec.Password,
		IsAdmin:  rec.IsAdmin,
	}

	_, err := rep.collectionU.InsertOne(c, cat)
	if err != nil {
		return err
	}
	return err
}

// SelectU function for selecting item from a table
func (rep *RepositoryMongo) SelectU(c context.Context, username, password string) (*model.User, error) {
	var rec model.User
	err := rep.collectionU.FindOne(c, bson.M{"username": username, "password": password}).Decode(&rec)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Errorf("Not found : %s\n", err)
		return &rec, err
	} else if err != nil {
		return &rec, err
	}
	return &rec, nil
}

// SelectAllU function for selecting items from a table
func (rep *RepositoryMongo) SelectAllU(c context.Context) ([]*model.User, error) {
	cursor, err := rep.collectionU.Find(c, bson.M{})
	if err != nil {
		return nil, err
	}

	var result []*model.User

	for cursor.Next(c) {
		rec := new(model.User)
		if err := cursor.Decode(rec); err != nil {
			return nil, err
		}
		result = append(result, rec)
	}

	if err := cursor.Close(c); err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateU function for updating item from a table
func (rep *RepositoryMongo) UpdateU(c context.Context, username string, isAdmin bool) error {
	if r, err := rep.collectionU.UpdateOne(c, bson.M{"username": username}, bson.M{"$set": bson.M{"is_admin": isAdmin}}); err != nil {
		return err
	} else if r.MatchedCount == 0 {
		log.Errorf("Not found : %s\n", err)
		return err
	}
	return nil
}

// DeleteU function for deleting item from a table
func (rep *RepositoryMongo) DeleteU(c context.Context, username string) error {
	if r, err := rep.collectionU.DeleteOne(c, bson.M{"username": username}); err != nil {
		return err
	} else if r.DeletedCount == 0 {
		log.Errorf("Not found : %s\n", err)
		return err
	}
	return nil
}

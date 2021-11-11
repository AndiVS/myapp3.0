package repository

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	model "myapp3.0/internal/model"
)

// InsertUser function for inserting item from a table
func (rep *Mongo) InsertUser(c context.Context, user *model.User) error {
	_, err := rep.collectionUsers.InsertOne(c, user)
	if err != nil {
		return err
	}
	return err
}

// SelectUser function for selecting item from a table
func (rep *Mongo) SelectUser(c context.Context, username string) (*model.User, error) {
	var user model.User
	err := rep.collectionUsers.FindOne(c, bson.M{"username": username}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Errorf("Not found : %s\n", err)
		return &user, err
	} else if err != nil {
		return &user, err
	}
	return &user, nil
}

// SelectAllUser function for selecting items from a table
func (rep *Mongo) SelectAllUser(c context.Context) ([]*model.User, error) {
	cursor, err := rep.collectionUsers.Find(c, bson.M{})
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

// UpdateUser function for updating item from a table
func (rep *Mongo) UpdateUser(c context.Context, username string, isAdmin bool) error {
	if r, err := rep.collectionUsers.UpdateOne(c, bson.M{"username": username}, bson.M{"$set": bson.M{"is_admin": isAdmin}}); err != nil {
		return err
	} else if r.MatchedCount == 0 {
		log.Errorf("Not found : %s\n", err)
		return err
	}
	return nil
}

// DeleteUser function for deleting item from a table
func (rep *Mongo) DeleteUser(c context.Context, username string) error {
	if r, err := rep.collectionUsers.DeleteOne(c, bson.M{"username": username}); err != nil {
		return err
	} else if r.DeletedCount == 0 {
		log.Errorf("Not found : %s\n", err)
		return ErrNotFound
	}
	return nil
}

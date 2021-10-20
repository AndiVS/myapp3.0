package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	model "myapp3.0/internal/model"
)

// InsertC function for inserting item from a table
func (rep *Mongo) InsertC(c context.Context, rec *model.Record) (uuid.UUID, error) {

	_, err := rep.collectionC.InsertOne(c, rec)
	if err != nil {
		return rec.ID, err
	}
	return rec.ID, err
}

// SelectC function for selecting item from a table
func (rep *Mongo) SelectC(c context.Context, id uuid.UUID) (*model.Record, error) {
	var rec model.Record
	err := rep.collectionC.FindOne(c, bson.M{"_id": id}).Decode(&rec)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Errorf("Not found : %s\n", err)
		return &rec, ErrNotFound
	} else if err != nil {
		return &rec, err
	}
	return &rec, nil
}

// SelectAllC function for selecting items from a table
func (rep *Mongo) SelectAllC(c context.Context) ([]*model.Record, error) {
	var result []*model.Record

	cursor, err := rep.collectionC.Find(c, bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(c) {
		rec := new(model.Record)
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

// UpdateC function for updating item from a table
func (rep *Mongo) UpdateC(c context.Context, rec *model.Record) error {
	if r, err := rep.collectionC.UpdateOne(c, bson.M{"_id": rec.ID}, bson.M{"$set": bson.M{"name": rec.Name, "type": rec.Type}}); err != nil {
		return err
	} else if r.MatchedCount == 0 {
		return err
	} else if r.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteC function for deleting item from a table
func (rep *Mongo) DeleteC(c context.Context, id uuid.UUID) error {
	if r, err := rep.collectionC.DeleteOne(c, bson.M{"_id": id}); err != nil {
		return err
	} else if r.DeletedCount == 0 {
		log.Errorf("Not found : %s\n", err)
		return err
	}
	return nil
}

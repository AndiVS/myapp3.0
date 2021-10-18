package repository

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	model "myapp3.0/internal/model"
)

// InsertC function for inserting item from a table
func (rep *RepositoryMongo) InsertC(c context.Context, rec *model.Record) error {

	_, err := rep.collectionC.InsertOne(c, rec)
	if err != nil {
		return err
	}
	return err
}

// SelectC function for selecting item from a table
func (rep *RepositoryMongo) SelectC(c context.Context, id string) (*model.Record, error) {
	var rec model.Record
	err := rep.collectionC.FindOne(c, bson.M{"id": id}).Decode(&rec)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Errorf("Not found : %s\n", err)
		return &rec, err
	} else if err != nil {
		return &rec, err
	}
	return &rec, nil
}

// SelectAllC function for selecting items from a table
func (rep *RepositoryMongo) SelectAllC(c context.Context) ([]*model.Record, error) {
	cursor, err := rep.collectionC.Find(c, bson.M{})
	if err != nil {
		return nil, err
	}

	var result []*model.Record

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
func (rep *RepositoryMongo) UpdateC(c context.Context, rec *model.Record) error {
	if r, err := rep.collectionC.UpdateOne(c, bson.M{"id": rec.ID}, bson.M{"$set": bson.M{"name": rec.Name, "type": rec.Type}}); err != nil {
		return err
	} else if r.MatchedCount == 0 {
		log.Errorf("Not found : %s\n", err)
		return err
	}
	return nil
}

// DeleteC function for deleting item from a table
func (rep *RepositoryMongo) DeleteC(c context.Context, id string) error {
	if r, err := rep.collectionC.DeleteOne(c, bson.M{"id": id}); err != nil {
		return err
	} else if r.DeletedCount == 0 {
		log.Errorf("Not found : %s\n", err)
		return err
	}
	return nil
}

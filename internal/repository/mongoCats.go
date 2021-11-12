package repository

import (
	"context"
	"errors"

	model "github.com/AndiVS/myapp3.0/internal/model"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertCat function for inserting item from a table
func (rep *Mongo) InsertCat(c context.Context, cat *model.Cat) (uuid.UUID, error) {
	_, err := rep.collectionCats.InsertOne(c, cat)
	if err != nil {
		return cat.ID, err
	}
	return cat.ID, err
}

// SelectCat function for selecting item from a table
func (rep *Mongo) SelectCat(c context.Context, id uuid.UUID) (*model.Cat, error) {
	var cat model.Cat
	err := rep.collectionCats.FindOne(c, bson.M{"_id": id}).Decode(&cat)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Errorf("Not found : %s\n", err)
		return &cat, ErrNotFound
	} else if err != nil {
		return &cat, err
	}
	return &cat, nil
}

// SelectAllCat function for selecting items from a table
func (rep *Mongo) SelectAllCat(c context.Context) ([]*model.Cat, error) {
	var cats []*model.Cat

	cursor, err := rep.collectionCats.Find(c, bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(c) {
		rec := new(model.Cat)
		if err := cursor.Decode(rec); err != nil {
			return nil, err
		}
		cats = append(cats, rec)
	}

	if err := cursor.Close(c); err != nil {
		return nil, err
	}
	return cats, nil
}

// UpdateCat function for updating item from a table
func (rep *Mongo) UpdateCat(c context.Context, cat *model.Cat) error {
	if r, err := rep.collectionCats.UpdateOne(c, bson.M{"_id": cat.ID}, bson.M{"$set": bson.M{"name": cat.Name, "type": cat.Type}}); err != nil {
		return err
	} else if r.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteCat function for deleting item from a table
func (rep *Mongo) DeleteCat(c context.Context, id uuid.UUID) error {
	if r, err := rep.collectionCats.DeleteOne(c, bson.M{"_id": id}); err != nil {
		return err
	} else if r.DeletedCount == 0 {
		log.Errorf("Not found : %s\n", err)
		return err
	}
	return nil
}

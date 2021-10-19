// Package model contain model of struct
package model

import "github.com/google/uuid"

// Record struct that contain record info
type Record struct {
	ID   uuid.UUID `param:"_id" query:"_id" header:"_id" form:"_id" json:"_id" xml:"_id" bson:"_id"`
	Name string    `param:"name" query:"name" header:"name" form:"name" json:"name" xml:"name" bson:"name"`
	Type string    `param:"type" query:"type" header:"type" form:"type" json:"type" xml:"type" bson:"type"`
}

// Package model contain model of struct
package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Cat struct that contain record info
type Cat struct {
	ID   uuid.UUID `param:"_id" query:"_id" header:"_id" form:"_id" json:"_id" xml:"_id" bson:"_id" msg:"_id"`
	Name string    `param:"name" query:"name" header:"name" form:"name" json:"name" xml:"name" bson:"name" msg:"name"`
	Type string    `param:"type" query:"type" header:"type" form:"type" json:"type" xml:"type" bson:"type" msg:"type"`
}

// MarshalBinary Marshal cat for redis stream
func (cat *Cat) MarshalBinary() ([]byte, error) {
	return json.Marshal(cat)
}

// UnmarshalBinary Marshal cat for redis stream
func (cat *Cat) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, cat)
}

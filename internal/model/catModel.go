// Package model contain model of struct
package model

// Record struct that contain record info
type Record struct {
	ID   int    `param:"id" query:"id" header:"id" form:"id" json:"id" xml:"id"`
	Name string `param:"name" query:"name" header:"name" form:"name" json:"name" xml:"name"`
	Type string `param:"type" query:"type" header:"type" form:"type" json:"type" xml:"type"`
}

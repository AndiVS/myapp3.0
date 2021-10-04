package model

type Record struct {
	Id   int    `json:"id" xml:"id" form:"id" query:"id"`
	Name string `json:"name" xml:"name" form:"name" query:"name"`
	Type string `json:"type" xml:"type" form:"type" query:"type"`
}
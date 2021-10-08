package model

import "github.com/golang-jwt/jwt"

type Claims struct {
	jwt.StandardClaims
	UserName string `param:"username" query:"username" header:"username" form:"username" xml:"username" json:"username,omitempty"`
}

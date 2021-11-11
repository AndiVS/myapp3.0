package model

import "github.com/golang-jwt/jwt"

// Claims for JWT
type Claims struct {
	jwt.StandardClaims
	Username string `param:"username" query:"username" header:"username" form:"username" xml:"username" json:"username,omitempty"`
	IsAdmin  bool   `param:"is_admin" query:"is_admin" header:"is_admin" form:"is_admin" xml:"is_admin" json:"is_admin,omitempty"`
}

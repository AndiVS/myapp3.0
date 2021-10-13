package model

// User for JWT
type User struct {
	Username string `param:"username" query:"username" header:"username" form:"username" xml:"username" json:"username,omitempty"`
	Password string `param:"password" query:"password" header:"password" form:"password" xml:"password" json:"password,omitempty"`
	IsAdmin  bool   `param:"is_admin" query:"is_admin" header:"is_admin" form:"is_admin" xml:"is_admin" json:"is_admin,omitempty"`
}

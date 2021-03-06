// Package middlewares for JWT
package middlewares

import (
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"net/http"
)

// Check for checking roll isAdmin
func Check(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("user") == nil {
			return next(c)
		}
		u := c.Get("user").(*jwt.Token)

		claims := u.Claims.(*model.Claims)

		if !claims.IsAdmin {
			return echo.NewHTTPError(http.StatusNotAcceptable, "have no access")
		}
		return next(c)
	}
}

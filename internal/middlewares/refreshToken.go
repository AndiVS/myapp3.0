package middlewares

import (
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/AndiVS/myapp3.0/internal/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"net/http"
	"time"
)

func TokenRefresherMiddleware(access, refresh *service.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Get("user") == nil {
				return next(c)
			}
			u := c.Get("user").(*jwt.Token)

			claims := u.Claims.(*model.Claims)

			if time.Until(time.Unix(claims.ExpiresAt, 0)) < 10*time.Minute {
				// Gets the refresh token from the cookie.
				rc, err := c.Cookie("refreshToken")
				if err == nil && rc != nil {
					// Parses token and checks if it valid.
					tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
						return refresh.SecretKey, nil
					})
					if err != nil {
						if err == jwt.ErrSignatureInvalid {
							c.Response().Writer.WriteHeader(http.StatusUnauthorized)
						}
					}

					if tkn != nil && tkn.Valid {
						// If everything is good, update tokens.
						_, _, _ = service.GenerateTokens(&model.User{
							Username: claims.Username,
							IsAdmin:  claims.IsAdmin,
						}, access, refresh)
					}
				}
			}

			return next(c)
		}
	}
}

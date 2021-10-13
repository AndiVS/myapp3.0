// Package service encapsulates
package service

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"

	"net/http"
	"time"
)

// Authorizer for generating token
type Authorizer struct {
	Rep *repository.Repository

	hashSalt           string
	authenticationKey  []byte
	refreshKey         []byte
	auntExpireDuration time.Duration
	refExpireDuration  time.Duration
}

// NewAuthorizer func for creating new authorizer
func NewAuthorizer(repositor *repository.Repository, hashSalt string, authenticationKey, refreshKey []byte, auntExpireDuration, refExpireDuration time.Duration) *Authorizer {
	return &Authorizer{
		Rep:                repositor,
		hashSalt:           hashSalt,
		authenticationKey:  authenticationKey,
		refreshKey:         refreshKey,
		auntExpireDuration: auntExpireDuration,
		refExpireDuration:  refExpireDuration,
	}
}

// GetU provides cat
func (author *Authorizer) GetU(c context.Context, username, password string) (*model.User, error) {
	return author.Rep.SelectU(c, username, password)
}

// GetAllU provides all cats
func (author *Authorizer) GetAllU(c context.Context) ([]*model.User, error) {
	return author.Rep.SelectAllU(c)
}

// AddU record about cat
func (author *Authorizer) AddU(c context.Context, rec *model.User) error {
	pwd := sha256.New()
	pwd.Write([]byte(rec.Password))
	pwd.Write([]byte(author.hashSalt))
	rec.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	return author.Rep.InsertU(c, rec)
}

// UpdateU updating record about cat
func (author *Authorizer) UpdateU(c context.Context, username string, isAdmin bool) error {
	return author.Rep.UpdateU(c, username, isAdmin)
}

// DeleteU record about cat
func (author *Authorizer) DeleteU(c context.Context, username string) error {
	return author.Rep.DeleteU(c, username)
}

// SignIn generate token
func (author *Authorizer) SignIn(c echo.Context, user *model.User) error {
	pwd := sha256.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(author.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := author.Rep.SelectU(c.Request().Context(), user.Username, user.Password)
	if err != nil {
		return err
	}

	return author.GenerateTokensAndSetCookies(user, c)
}

// GenerateTokensAndSetCookies func for token generation
func (author *Authorizer) GenerateTokensAndSetCookies(user *model.User, c echo.Context) error {
	accessToken, err := generateToken(user, author.auntExpireDuration, author.authenticationKey)
	if err != nil {
		return err
	}

	refreshToken, err := generateToken(user, author.refExpireDuration, author.refreshKey)
	if err != nil {
		return err
	}

	setUserCookie(user, time.Now().Add(author.auntExpireDuration), c)

	setTokenCookie("refreshToken", refreshToken, time.Now().Add(author.refExpireDuration), c)

	return c.JSON(http.StatusOK, echo.Map{
		"token": accessToken,
	})
}

//
func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)
}

// GenerateToken func for token generation
func generateToken(user *model.User, expire time.Duration, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func setUserCookie(user *model.User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = user.Username
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}

// TokenRefresherMiddleware func for refreshing toke
func (author *Authorizer) TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
					return author.refreshKey, nil
				})
				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					// If everything is good, update tokens.
					_ = author.GenerateTokensAndSetCookies(&model.User{
						Username: claims.Username,
						IsAdmin:  claims.IsAdmin,
					}, c)
				}
			}
		}

		return next(c)
	}
}

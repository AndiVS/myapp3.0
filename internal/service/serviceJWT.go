package service

import (
	"crypto/sha256"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"myapp3.0/internal/model"

	"net/http"
	"time"
)

// JWTManager is a JSON web token manager
type JWTManager struct {
	SecretKey     []byte
	TokenDuration time.Duration
}

// NewJWTManager returns a new JWT manager
func NewJWTManager(secretKey []byte, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{SecretKey: secretKey, TokenDuration: tokenDuration}
}

// GenerateTokens func for token generation
func GenerateTokens(user *model.User, access *JWTManager, refresh *JWTManager) (string, string, error) {
	accessToken, err := GenerateToken(user, access)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateToken(user, refresh)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func PasswordChek(pass1, pass2 string) bool {
	return pass1 == pass2
}

// PasswordGenerator generate password from hash and pass string
func PasswordGenerator(password, hashSalt string) string {
	pwd := sha256.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(hashSalt))

	return fmt.Sprintf("%x", pwd.Sum(nil))
}

func GenerateToken(user *model.User, manager *JWTManager) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.TokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	})

	tokenString, err := token.SignedString(manager.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)
}

func SetUserCookie(user *model.User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = user.Username
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}

// Verify verifies the access token string and return a user claim if the token is valid
func (manager *JWTManager) Verify(Token string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(
		Token,
		&model.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return manager.SecretKey, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

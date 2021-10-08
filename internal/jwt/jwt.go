// Package jwt for authorization
package jwt

import (
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"myapp3.0/internal/model"

	"net/http"
	"time"
)

// CheckJWT check JWT
func CheckJWT(tknStr string, jwtKey []byte) (*model.Claims, error) {
	claims := &model.Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}
	if !tkn.Valid {
		return nil, err
	}

	return claims, nil
}

// CheckJWTFromCookie load JWT check from Cookie and check it
func CheckJWTFromCookie(cookie *http.Cookie, jwtKey []byte) error {
	_, err := CheckJWT(cookie.Value, jwtKey)

	return err
}

// CreateJWT create new JWT
func CreateJWT(claims *model.Claims, expirationTime *time.Time, jwtKey []byte) (string, error) {
	if expirationTime != nil {
		claims.StandardClaims.ExpiresAt = expirationTime.Unix()
	}

	log.Errorf("JWT: claims %v", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	log.Errorf("JWT: tokenString %v", tokenString)

	return tokenString, nil
}

// CreateJWTCookie create new JWT as Cookie
func CreateJWTCookie(claims *model.Claims, jwtExpiresAt int, jwtKey []byte) (*http.Cookie, error) {
	var expirationTime *time.Time

	// jwtExpiresAt > 0 установим expiry time
	if jwtExpiresAt > 0 {
		t := time.Now().Add(time.Duration(jwtExpiresAt * int(time.Second)))
		expirationTime = &t
	}

	// создадим новый токен
	tokenString, err := CreateJWT(claims, expirationTime, jwtKey)
	if err != nil {
		return nil, err
	}

	// подготовим Cookie
	cookie := http.Cookie{
		Name:  "token",
		Value: tokenString,
	}

	if jwtExpiresAt > 0 {
		// set an expiry time is the same as the token itself
		cookie.Expires = *expirationTime
	} else {
		cookie.MaxAge = 0 // без ограничения времени жизни
	}

	return &cookie, nil
}

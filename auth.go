package main

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtData struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var jwtSecret = []byte("078fcc13c}84$f6c923c89394{bbd0#e")

func CreateJWT(email string, expTimeInMinutes float64) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expTimeInMinutes) * time.Minute)
	jwtData := &JwtData{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtData)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err;
	}
	return tokenString, nil;
}

func VerifyJWT(jwtToken string) (*JwtData, error) {
	jwtData := &JwtData{}
	tkn, err := jwt.ParseWithClaims(jwtToken, jwtData, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return jwtData, errors.New("invalid jwt signature") 
		}
		return jwtData, errors.New("invalid jwt") 
	}
	if !tkn.Valid {
		return jwtData, errors.New("invalid token") 
	}
	return jwtData, nil
}
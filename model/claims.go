package model

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Uuid string `json:"uuid"`
	jwt.StandardClaims
}

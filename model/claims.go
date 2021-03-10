package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/vectorman1/analysis/analysis-api/model/db"
)

type Claims struct {
	Uuid        string         `json:"uuid"`
	PrivateRole db.PrivateRole `json:"privateRole"`
	jwt.StandardClaims
}

package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/vectorman1/analysis/analysis-api/model/db/entities"
)

type Claims struct {
	Uuid        string               `json:"uuid"`
	PrivateRole entities.PrivateRole `json:"privateRole"`
	jwt.StandardClaims
}

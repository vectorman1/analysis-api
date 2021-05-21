package common

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/vectorman1/analysis/analysis-api/domain/user/model"
)

type Claims struct {
	Uuid        string            `json:"uuid"`
	PrivateRole model.PrivateRole `json:"privateRole"`
	jwt.StandardClaims
}

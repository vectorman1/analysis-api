package db

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/pgtype"
)

type User struct {
	ID        uint
	Uuid      pgtype.UUID
	Username  string
	Password  string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	DeletedAt pgtype.Timestamptz
}

type Claims struct {
	Uuid string `json:"uuid"`
	jwt.StandardClaims
}

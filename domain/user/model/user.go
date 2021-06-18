package model

import (
	"github.com/jackc/pgtype"
	"github.com/vectorman1/analysis/analysis-api/generated/user_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PrivateRole uint

const (
	Default PrivateRole = iota
	Admin
)

type User struct {
	ID          uint
	Uuid        pgtype.UUID
	PrivateRole PrivateRole
	Username    string
	Password    string
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	DeletedAt   pgtype.Timestamptz
}

func (e *User) ToProto() *user_service.User {
	var u string
	e.Uuid.AssignTo(&u)

	return &user_service.User{
		Id:          uint64(e.ID),
		Uuid:        u,
		Username:    e.Username,
		Password:    e.Password,
		PrivateRole: uint32(e.PrivateRole),
		CreatedAt:   timestamppb.New(e.CreatedAt.Time),
		UpdatedAt:   timestamppb.New(e.UpdatedAt.Time),
	}
}

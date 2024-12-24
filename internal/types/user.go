package types

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthPayloadErrMessage struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type RegisterUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserOutput struct {
	ID        uuid.UUID        `json:"id"`
	Email     string           `json:"email"`
	Admin     bool             `json:"admin"`
	Password  string           `json:"password"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
	UpdatedAt pgtype.Timestamp `json:"updatedAt"`
}

type RegisterUserErrMessage struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserOutput struct {
	Token string `json:"token"`
}

type LoginUserErrMessage struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

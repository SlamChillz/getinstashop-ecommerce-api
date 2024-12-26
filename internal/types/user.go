package types

import (
	"github.com/google/uuid"
	"time"
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
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Admin     bool      `json:"admin"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RegisterUserErrMessage struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// CreateUserOk For Swagger Docs
type CreateUserOk struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    RegisterUserOutput `json:"data"`
}

// CreateUserError For Swagger Docs
type CreateUserError struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Error   RegisterUserErrMessage `json:"error"`
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

// LoginUserError For Swagger Docs
type LoginUserError struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Error   RegisterUserErrMessage `json:"error"`
}

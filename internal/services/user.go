package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/utils"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/validators"
	"github.com/slamchillz/getinstashop-ecommerce-api/pkg/token"
	"net/http"
	"strings"
)

// UserService provides business logic for user operations.
type UserService struct {
	store    db.Store
	jwtToken *token.JWT
}

// NewUserService creates a new UserService instance.
func NewUserService(store db.Store, jwtToken *token.JWT) *UserService {
	return &UserService{
		store:    store,
		jwtToken: jwtToken,
	}
}

// CreateUser creates a new user in the database.
func (s *UserService) CreateUser(ctx context.Context, user types.RegisterUserInput) (types.RegisterUserOutput, types.RegisterUserErrMessage, int, error) {
	// Validate user input (e.g., email uniqueness)
	var newUserOutput types.RegisterUserOutput
	errMessage, err := validators.ValidateAuthPayload(types.AuthPayload(user))
	if err != nil {
		return newUserOutput, errMessage, http.StatusBadRequest, err
	}
	// Create the user in the database
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return newUserOutput, errMessage, http.StatusInternalServerError, err
	}
	newUser, err := s.store.CreateUser(ctx, db.CreateUserParams{
		ID:       uuid.New(),
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				errMessage.Email = "email already taken"
				return newUserOutput, errMessage, http.StatusBadRequest, err
			}
		}
		return newUserOutput, errMessage, http.StatusInternalServerError, err
	}
	newUserOutput = types.RegisterUserOutput{
		ID:        newUser.ID,
		Email:     newUser.Email,
		Admin:     newUser.Admin,
		CreatedAt: newUser.CreatedAt.Time,
		UpdatedAt: newUser.UpdatedAt.Time,
	}
	return newUserOutput, errMessage, http.StatusCreated, nil
}

// LoginUser validates user credentials and returns a valid token for subsequent requests
func (s *UserService) LoginUser(ctx context.Context, user types.LoginUserInput) (types.LoginUserOutput, types.RegisterUserErrMessage, int, error) {
	// Validate user input (e.g., email uniqueness)
	var output types.LoginUserOutput
	errMessage, err := validators.ValidateAuthPayload(types.AuthPayload(user))
	if err != nil {
		return output, errMessage, http.StatusBadRequest, err
	}
	dbUser, err := s.store.GetUserById(ctx, user.Email)
	if err != nil {
		if strings.Replace(sql.ErrNoRows.Error(), "sql: ", "", 1) == err.Error() {
			errMessage.Email = "email does not exist"
			return output, errMessage, http.StatusNotFound, err
		}
		return output, errMessage, http.StatusInternalServerError, err
	}
	err = utils.CheckPassword(dbUser.Password, user.Password)
	if err != nil {
		errMessage.Password = "invalid password"
		return output, errMessage, http.StatusBadRequest, err
	}
	accessToken, err := s.jwtToken.CreateToken(dbUser.ID, dbUser.Admin)
	if err != nil {
		return output, errMessage, http.StatusInternalServerError, err
	}
	output = types.LoginUserOutput{
		Token: accessToken,
	}
	return output, errMessage, http.StatusOK, nil
}

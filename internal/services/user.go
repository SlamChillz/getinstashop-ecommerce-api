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
func (s *UserService) CreateUser(ctx context.Context, user types.RegisterUserInput) (*types.RegisterUserOutput, *types.RegisterUserErrMessage, int, error) {
	// Validate user input (e.g., email uniqueness)
	errMessage, err := validators.ValidateAuthPayload(types.AuthPayload(user))
	if err != nil {
		return nil, (*types.RegisterUserErrMessage)(&errMessage), http.StatusBadRequest, err
	}
	// Create the user in the database
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
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
				return nil, &types.RegisterUserErrMessage{
					Email: "email has been taken",
				}, http.StatusBadRequest, err
			}
		}
		return nil, nil, http.StatusInternalServerError, err
	}
	var newUserOutput = types.RegisterUserOutput{
		ID:        newUser.ID,
		Email:     newUser.Email,
		Admin:     newUser.Admin,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}
	return &newUserOutput, nil, http.StatusCreated, nil
}

// LoginUser validates user credentials and returns a valid token for subsequent requests
func (s *UserService) LoginUser(ctx context.Context, user types.LoginUserInput) (*types.LoginUserOutput, *types.LoginUserErrMessage, int, error) {
	// Validate user input (e.g., email uniqueness)
	errMessage, err := validators.ValidateAuthPayload(types.AuthPayload(user))
	if err != nil {
		return nil, (*types.LoginUserErrMessage)(&errMessage), http.StatusBadRequest, err
	}
	dbUser, err := s.store.GetUserById(ctx, user.Email)
	if err != nil {
		if strings.Replace(sql.ErrNoRows.Error(), "sql: ", "", 1) == err.Error() {
			return nil, &types.LoginUserErrMessage{
				Email: "email does not exist",
			}, http.StatusNotFound, err
		}
		return nil, nil, http.StatusInternalServerError, err
	}
	err = utils.CheckPassword(dbUser.Password, user.Password)
	if err != nil {
		return nil, &types.LoginUserErrMessage{
			Password: "invalid password",
		}, http.StatusBadRequest, err
	}
	accessToken, err := s.jwtToken.CreateToken(dbUser.ID, dbUser.Admin)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}
	var output = types.LoginUserOutput{
		Token: accessToken,
	}
	return &output, nil, http.StatusOK, nil
}

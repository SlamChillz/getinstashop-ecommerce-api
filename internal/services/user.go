package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/utils"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/validators"
	"net/http"
)

// UserService provides business logic for user operations.
type UserService struct {
	store db.Store
}

// NewUserService creates a new UserService instance.
func NewUserService(store db.Store) *UserService {
	return &UserService{
		store: store,
	}
}

// CreateUser creates a new user in the database.
func (s *UserService) CreateUser(ctx context.Context, user types.RegisterUserInput) (*types.RegisterUserOutput, *types.RegisterUserErrMessage, int, error) {
	// Validate user input (e.g., email uniqueness)
	errMessage, err := validators.RegisterUser(user)
	if err != nil {
		return nil, &errMessage, http.StatusBadRequest, err
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
